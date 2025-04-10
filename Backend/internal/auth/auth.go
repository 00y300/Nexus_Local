package auth

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// PageData holds the data passed to the HTML template.
type PageData struct {
	LoginURL    string
	RedirectURL string
}

// App holds the OAuth2 config, OIDC verifier, and HTML template.
type App struct {
	OAuthCfg *oauth2.Config
	Verifier *oidc.IDTokenVerifier
	Tmpl     *template.Template
}

// NewApp constructs a new App with OAuth2 config, verifier, and template.
func NewApp(oauthCfg *oauth2.Config, verifier *oidc.IDTokenVerifier, tmpl *template.Template) *App {
	return &App{
		OAuthCfg: oauthCfg,
		Verifier: verifier,
		Tmpl:     tmpl,
	}
}

var (
	// ErrNoAuthHeader is returned when no Authorization header is present.
	ErrNoAuthHeader = errors.New("no Authorization header")
	// ErrInvalidToken is returned when the token is missing or invalid.
	ErrInvalidToken = errors.New("invalid token")
	// ErrForbidden is returned when the user lacks the "admin" role.
	ErrForbidden = errors.New("forbidden")
)

// logScopes extracts the "scope" field from the token response,
// splits it on spaces, and logs the resulting slice.
func logScopes(token *oauth2.Token) {
	raw, ok := token.Extra("scope").(string)
	if !ok || raw == "" {
		log.Println("⚠️  no scopes returned in token response")
		return
	}
	scopes := strings.Fields(raw)
	log.Printf("✅ granted scopes: %v\n", scopes)
}

// AuthMiddleware ensures the request has a valid Bearer token
// and that the "roles" claim includes "admin".
func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := r.Header.Get("Authorization")
		if hdr == "" {
			http.Error(w, ErrNoAuthHeader.Error(), http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(hdr, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		rawIDToken := parts[1]

		// Verify signature & standard claims
		idToken, err := a.Verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		// Extract roles claim
		var claims struct {
			Roles []string `json:"roles"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Enforce admin role
		if !slices.Contains(claims.Roles, "admin") {
			http.Error(w, ErrForbidden.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Root renders the login page.
func (a *App) Root(w http.ResponseWriter, r *http.Request) {
	state := "some-random-state"
	loginURL := a.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	data := PageData{LoginURL: loginURL, RedirectURL: a.OAuthCfg.RedirectURL}
	if err := a.Tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Login redirects the user to Azure AD for authentication.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	state := "some-random-state"
	url := a.OAuthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Printf("Redirecting to: %s", url)
	http.Redirect(w, r, url, http.StatusFound)
}

// OAuthCallback handles the code exchange, logs scopes, sets ID and access token cookies,
// and redirects back to your Next.js admin page.
func (a *App) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	// 1) Get the authorization code from the query params
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	// 2) Exchange the code for an OAuth2 token
	token, err := a.OAuthCfg.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3) Log the granted scopes for debugging
	logScopes(token)

	// 4) (Optional) Log raw tokens for debugging
	log.Printf("Access Token: %s", token.AccessToken)
	if rt := token.RefreshToken; rt != "" {
		log.Printf("Refresh Token: %s", rt)
	}
	if idtRaw, ok := token.Extra("id_token").(string); ok {
		log.Printf("ID Token: %s", idtRaw)
	}

	// 5) Extract the ID token and Access token
	idt, _ := token.Extra("id_token").(string)
	at := token.AccessToken

	// 6a) Set the ID token as an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "id_token",
		Value:    idt,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // enable in production
		// SameSite: http.SameSiteLaxMode,
	})

	// 6b) Set the Access token as an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    at,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // enable in production
		// SameSite: http.SameSiteLaxMode,
	})

	// 7) Redirect back to your Next.js admin page
	http.Redirect(w, r, "http://localhost:3000/admin/add-item", http.StatusFound)
}
