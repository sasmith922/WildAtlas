package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/sasmith922/WildAtlas/internal/scraper"
)

// Handler contains the HTTP handlers for the API
type Handler struct {
	scraper *scraper.Scraper
}

// NewHandler creates a new Handler instance
func NewHandler() *Handler {
	return &Handler{
		scraper: scraper.NewScraper(),
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// GetSpeciesByCountry handles GET requests for species by country
func (h *Handler) GetSpeciesByCountry(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "method_not_allowed",
			Message: "Only GET method is allowed",
		})
		return
	}

	// Extract country code from URL path
	// Expected path: /api/species/{countryCode}
	path := strings.TrimPrefix(r.URL.Path, "/api/species/")
	countryCode := strings.TrimSpace(path)

	if countryCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "missing_country_code",
			Message: "Country code is required. Use /api/species/{countryCode}",
		})
		return
	}

	// Validate country code (should be 2 letters)
	if len(countryCode) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "invalid_country_code",
			Message: "Country code must be a 2-letter ISO 3166-1 alpha-2 code",
		})
		return
	}

	// Get species data
	data, err := h.scraper.GetSpeciesByCountry(countryCode)
	if err != nil {
		log.Printf("Error fetching species data for %s: %v", countryCode, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "fetch_error",
			Message: "Failed to fetch species data",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
