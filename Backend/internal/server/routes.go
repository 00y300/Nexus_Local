package server

import (
	"fmt"
	"net/http"
)

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

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// routes sets up the URL routes and associates them with their handlers.
func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello/", s.defaultHandler)

	// OAuth routes
	mux.HandleFunc("/", s.AuthApp.Root)
	mux.HandleFunc("/login", s.AuthApp.Login)
	mux.HandleFunc("/redirect", s.AuthApp.OAuthCallback)

	return corsMiddleware(mux)
}
