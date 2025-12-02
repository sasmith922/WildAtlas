package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sasmith922/WildAtlas/internal/iucn"
)

// TestHealthCheck verifies that the health check endpoint returns a 200 OK status
// and the expected JSON response.
func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response["status"] != "healthy" {
		t.Errorf("handler returned unexpected body: got %v want healthy", response["status"])
	}
}

// TestGetSpeciesByCountry_ValidCode verifies that a valid country code (US)
// returns a 200 OK status and correct species data.
func TestGetSpeciesByCountry_ValidCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/species/US", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response iucn.CountryData
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.CountryCode != "US" {
		t.Errorf("expected country code US, got %v", response.CountryCode)
	}

	if response.Country != "United States" {
		t.Errorf("expected country name United States, got %v", response.Country)
	}

	if len(response.Species) == 0 {
		t.Error("expected at least one species")
	}
}

// TestGetSpeciesByCountry_InvalidCode verifies that an invalid country code (not 2 letters)
// returns a 400 Bad Request status.
func TestGetSpeciesByCountry_InvalidCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/species/ABC", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// TestGetSpeciesByCountry_MissingCode verifies that missing the country code
// returns a 400 Bad Request status.
func TestGetSpeciesByCountry_MissingCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/species/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
