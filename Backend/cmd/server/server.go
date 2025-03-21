package server

import (
	"flag"
	"log"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there, I love " + r.URL.Path[1:] + "!"))
}

// corsMiddleware sets the CORS headers for each request.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins and specific methods/headers.
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

// StartServer initializes and starts the HTTP server.
func StartServer() {
	port := flag.String("port", "8000", "Port to run the HTTP server on")
	flag.Parse()

	mux := http.NewServeMux()

	// Register route handlers.
	mux.HandleFunc("/", defaultHandler)

	// Apply CORS middleware.
	handlerWithCORS := corsMiddleware(mux)
	addr := ":" + *port
	log.Printf("Starting server on %s", addr)

	// Start listening and serving requests.
	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
