package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sasmith922/WildAtlas/internal/handlers"
)

// main is the entry point of the application.
// It sets up the HTTP server, defines routes, and starts listening for incoming requests.
func main() {
	// Get port from environment variable (useful for deployment platforms like Heroku)
	// or use default port 8080 for local development.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up a new HTTP request multiplexer (router).
	mux := http.NewServeMux()

	// --- API Routes ---
	// Endpoint to get endangered species data for a specific country.
	// The country code is passed as part of the URL path.
	mux.HandleFunc("/api/species/", handlers.GetSpeciesByCountry)

	// Endpoint for health checks, useful for monitoring uptime.
	mux.HandleFunc("/api/health", handlers.HealthCheck)

	// --- Static File Serving ---
	// Serve the frontend files (HTML, CSS, JS) from the "static" directory.
	// This allows the Go server to host the entire application.
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	// Start the HTTP server.
	log.Printf("WildAtlas server starting on port %s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	// ListenAndServe starts an HTTP server with a given address and handler.
	// It returns an error if the server fails to start.
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
