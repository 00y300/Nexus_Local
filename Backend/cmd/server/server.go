package server

import (
	"fmt"
	"log"
	"net/http"

	"nexus.local/internal/auth"
)

// Server holds the dependencies for the HTTP server.
type Server struct {
	AuthApp *auth.App
}

// NewServer constructs a new Server instance using dependency injection.
func NewServer(authApp *auth.App) *Server {
	return &Server{
		AuthApp: authApp,
	}
}

// defaultHandler is a simple handler that shows a greeting message.
func (s *Server) defaultHandler(w http.ResponseWriter, r *http.Request) {
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

// routes sets up the URL routes and associates them with their handlers.
func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello/", s.defaultHandler)

	// Dependency injection for OAuth routes.
	mux.HandleFunc("/", s.AuthApp.Root)
	mux.HandleFunc("/login", s.AuthApp.Login)
	mux.HandleFunc("/redirect", s.AuthApp.OAuthCallback)

	return mux
}

// Start starts the HTTP server on the given address.
func (s *Server) Start(addr string) error {
	mux := s.routes()

	// Wrap all routes with the CORS middleware.
	handlerWithCORS := corsMiddleware(mux)

	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, handlerWithCORS)
}
