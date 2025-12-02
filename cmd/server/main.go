package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sasmith922/WildAtlas/internal/handlers"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create handlers
	h := handlers.NewHandler()

	// Set up routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/species/", h.GetSpeciesByCountry)
	mux.HandleFunc("/api/health", h.HealthCheck)

	// Serve static files (frontend)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	// Start server
	log.Printf("WildAtlas server starting on port %s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
