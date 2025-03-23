package server

import (
	"fmt"
	"log"
	"net/http"

	"nexus.local/internal/auth"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Hi there, I love %s!", r.URL.Path[len("/hello/"):])
	w.Write([]byte(msg))
}

// corsMiddleware sets the CORS headers for each request.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS requests.
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// StartServer initializes and starts the HTTP server on port 8080.
func StartServer(app *auth.App) {
	mux := http.NewServeMux()

	// Register your default route at /hello
	mux.HandleFunc("/hello/", defaultHandler)

	// Register OAuth routes:
	mux.HandleFunc("/", app.Root)
	mux.HandleFunc("/login", app.Login)
	mux.HandleFunc("/redirect", app.OAuthCallback)

	// Wrap everything with the CORS middleware.
	handlerWithCORS := corsMiddleware(mux)

	addr := ":8080"
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
