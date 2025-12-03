package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sasmith922/WildAtlas/internal/iucn"
)

// GetSpeciesByCountry handles GET requests to retrieve endangered species data.
// It acts as a bridge between the frontend and the IUCN API.
func GetSpeciesByCountry(w http.ResponseWriter, r *http.Request) {
	// Allow Cross-Origin Resource Sharing (CORS) so the frontend can call this API
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Extract country code from the URL path (e.g., /api/species/CA -> CA)
	path := strings.TrimPrefix(r.URL.Path, "/api/species/")
	countryCode := strings.TrimSpace(path)
	log.Printf("Handler received request for country: %s", countryCode) // Force log to stderr

	// Basic validation of the country code
	if countryCode == "" || len(countryCode) != 2 {
		http.Error(w, "Invalid country code", http.StatusBadRequest)
		return
	}

	// Initialize IUCN client with the API token from environment variables
	token := os.Getenv("IUCN_API_TOKEN")
	if token == "" {
		log.Println("WARNING: IUCN_API_TOKEN is not set")
	}
	client := iucn.NewClient(token)

	// Fetch data using the client
	data, err := client.GetSpeciesByCountry(countryCode)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// Serialize the data to JSON and send it back to the client
	json.NewEncoder(w).Encode(data)
}

// HealthCheck handles health check requests.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
