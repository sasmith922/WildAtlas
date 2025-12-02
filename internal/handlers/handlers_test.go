package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sasmith922/WildAtlas/internal/scraper"
)

func TestHealthCheck(t *testing.T) {
	h := NewHandler()

	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HealthCheck)

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

func TestGetSpeciesByCountry_ValidCode(t *testing.T) {
	h := NewHandler()

	req, err := http.NewRequest("GET", "/api/species/US", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response scraper.CountryData
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

func TestGetSpeciesByCountry_InvalidCode(t *testing.T) {
	h := NewHandler()

	req, err := http.NewRequest("GET", "/api/species/ABC", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetSpeciesByCountry_MissingCode(t *testing.T) {
	h := NewHandler()

	req, err := http.NewRequest("GET", "/api/species/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetSpeciesByCountry_MethodNotAllowed(t *testing.T) {
	h := NewHandler()

	req, err := http.NewRequest("POST", "/api/species/US", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestGetSpeciesByCountry_CORSHeaders(t *testing.T) {
	h := NewHandler()

	req, err := http.NewRequest("GET", "/api/species/US", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetSpeciesByCountry)

	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS header Access-Control-Allow-Origin")
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Error("missing Content-Type header")
	}
}
