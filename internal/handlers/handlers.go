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
func GetSpeciesByCountry(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Extract country code
	path := strings.TrimPrefix(r.URL.Path, "/api/species/")
	countryCode := strings.TrimSpace(path)

	if countryCode == "" || len(countryCode) != 2 {
		http.Error(w, "Invalid country code", http.StatusBadRequest)
		return
	}

	// Initialize IUCN client
	token := os.Getenv("IUCN_API_TOKEN")
	client := iucn.NewClient(token)

	// Fetch data
	data, err := client.GetSpeciesByCountry(countryCode)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// HealthCheck handles health check requests.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
