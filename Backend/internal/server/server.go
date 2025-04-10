// internal/server/server.go
package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"nexus.local/internal/auth"
)

// GraphUser models the subset of fields we care about from MS Graph /me
type GraphUser struct {
	ODataContext      string   `json:"@odata.context"`
	BusinessPhones    []string `json:"businessPhones"`
	DisplayName       string   `json:"displayName"`
	GivenName         string   `json:"givenName"`
	JobTitle          *string  `json:"jobTitle"`
	Mail              *string  `json:"mail"`
	MobilePhone       *string  `json:"mobilePhone"`
	OfficeLocation    *string  `json:"officeLocation"`
	PreferredLanguage *string  `json:"preferredLanguage"`
	Surname           string   `json:"surname"`
	UserPrincipalName string   `json:"userPrincipalName"`
	ID                string   `json:"id"`
}

// Server holds your OAuth app and the database pool.
type Server struct {
	AuthApp *auth.App
	DB      *sql.DB
}

// NewServer constructs a Server with its dependencies.
func NewServer(authApp *auth.App, db *sql.DB) *Server {
	return &Server{AuthApp: authApp, DB: db}
}

// routes wires up all handlers.
func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// OAuth endpoints
	mux.HandleFunc("/", s.AuthApp.Root)
	mux.HandleFunc("/login", s.AuthApp.Login)
	mux.HandleFunc("/redirect", s.AuthApp.OAuthCallback)

	// CRUD endpoints
	mux.HandleFunc("/items", s.getItemsHandler)
	mux.Handle("/items/add",
		s.AuthApp.AuthMiddleware(http.HandlerFunc(s.addItemHandler)),
	)
	mux.Handle("/items/update",
		s.AuthApp.AuthMiddleware(http.HandlerFunc(s.updateStockHandler)),
	)
	mux.HandleFunc("/orders", s.ordersHandler)

	// Graph profile + DB upsert
	mux.HandleFunc("/me", s.profileHandler)

	// Logout endpoint — clears the auth cookies
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		clear := func(name string) {
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				// SameSite: http.SameSiteNoneMode,
				// Secure:   false, // set to true in production
				SameSite: http.SameSiteNoneMode,
				Secure:   true, // ← must be true if SameSite=None
				MaxAge:   -1,
			})
		}
		clear("id_token")
		clear("access_token")
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}

// Start runs the HTTP server with CORS enabled.
func (s *Server) Start(addr string) error {
	handler := corsMiddleware(s.routes())
	return http.ListenAndServe(addr, handler)
}

// corsMiddleware sets CORS headers and allows credentials.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only allow your front‑end origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		// Allow cookies to be sent/received
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// profileHandler calls Graph /me, upserts the user into MySQL, then returns the JSON.
func (s *Server) profileHandler(w http.ResponseWriter, r *http.Request) {
	// 1) Grab the access_token from the cookie
	ck, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}
	at := ck.Value

	// 2) Call Graph /me
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+at)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "graph request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 3) Read the entire response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read graph response", http.StatusInternalServerError)
		return
	}

	// If Graph didn’t return 200, just pass it through
	if resp.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// 4) Decode into our GraphUser struct
	var user GraphUser
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "failed to parse graph response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5) Upsert into `users` table
	phonesJSON, err := json.Marshal(user.BusinessPhones)
	if err != nil {
		http.Error(w, "failed to marshal phones: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Exec(`
        INSERT INTO users (
            id,
            display_name,
            given_name,
            surname,
            job_title,
            mail,
            mobile_phone,
            office_location,
            preferred_language,
            user_principal_name,
            business_phones
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            display_name        = VALUES(display_name),
            given_name          = VALUES(given_name),
            surname             = VALUES(surname),
            job_title           = VALUES(job_title),
            mail                = VALUES(mail),
            mobile_phone        = VALUES(mobile_phone),
            office_location     = VALUES(office_location),
            preferred_language  = VALUES(preferred_language),
            user_principal_name = VALUES(user_principal_name),
            business_phones     = VALUES(business_phones)
    `,
		user.ID,
		user.DisplayName,
		user.GivenName,
		user.Surname,
		user.JobTitle,
		user.Mail,
		user.MobilePhone,
		user.OfficeLocation,
		user.PreferredLanguage,
		user.UserPrincipalName,
		phonesJSON,
	)
	if err != nil {
		http.Error(w, "failed to upsert user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6) Return the Graph JSON to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
