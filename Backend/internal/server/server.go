package server

import (
	"database/sql"
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

	// Public OAuth endpoints
	mux.HandleFunc("/", s.AuthApp.Root)
	mux.HandleFunc("/login", s.AuthApp.Login)
	mux.HandleFunc("/redirect", s.AuthApp.OAuthCallback)

	// Public: list items
	mux.HandleFunc("/items", s.getItemsHandler)

	// Admin-only: add & update stock
	mux.Handle(
		"/items/add",
		s.AuthApp.AuthMiddleware(http.HandlerFunc(s.addItemHandler)),
	)
	mux.Handle(
		"/items/update",
		s.AuthApp.AuthMiddleware(http.HandlerFunc(s.updateStockHandler)),
	)

	// Public or user-scoped orders endpoint
	mux.HandleFunc("/orders", s.ordersHandler)

	return mux
}

// Start runs the HTTP server with CORS enabled.
func (s *Server) Start(addr string) error {
	handler := corsMiddleware(s.routes())
	return http.ListenAndServe(addr, handler)
}

// corsMiddleware sets permissive CORS headers.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
