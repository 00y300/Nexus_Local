package server

import (
	"log"
	"net/http"

	"nexus.local/internal/auth"
)

// Server holds the dependencies for the HTTP server.
type Server struct {
	AuthApp *auth.App
}

// NewServer constructs a new Server instance.
func NewServer(authApp *auth.App) *Server {
	return &Server{AuthApp: authApp}
}

// Start runs the HTTP server on the given address.
func (s *Server) Start(addr string) error {
	handler := s.routes()
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, handler)
}
