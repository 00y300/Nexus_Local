package server

import (
	"database/sql"
	"log"
	"net/http"

	"nexus.local/internal/auth"
)

// Server holds your OAuth app and the database pool.
type Server struct {
	AuthApp *auth.App
	DB      *sql.DB
}

// NewServer constructs a Server with its dependencies.
func NewServer(authApp *auth.App, db *sql.DB) *Server {
	return &Server{
		AuthApp: authApp,
		DB:      db,
	}
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
	mux.HandleFunc("/items/add", s.addItemHandler)
	mux.HandleFunc("/items/update", s.updateStockHandler)
	mux.HandleFunc("/orders", s.ordersHandler)

	return mux
}

// Start runs the HTTP server with CORS enabled.
func (s *Server) Start(addr string) error {
	handler := corsMiddleware(s.routes())
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

// corsMiddleware sets permissive CORS headers.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
