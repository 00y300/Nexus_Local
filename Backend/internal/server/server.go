package server

import (
	"database/sql"
	"log"
	"net/http"

	"nexus.local/internal/auth"
)

type Server struct {
	AuthApp *auth.App
	DB      *sql.DB
}

func NewServer(authApp *auth.App, db *sql.DB) *Server {
	return &Server{AuthApp: authApp, DB: db}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.AuthApp.Root)
	mux.HandleFunc("/login", s.AuthApp.Login)
	mux.HandleFunc("/redirect", s.AuthApp.OAuthCallback)
	mux.HandleFunc("/items", s.getItemsHandler)
	mux.HandleFunc("/orders", s.placeOrderHandler)
	return mux
}

func (s *Server) Start(addr string) error {
	handler := corsMiddleware(s.routes())
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

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
