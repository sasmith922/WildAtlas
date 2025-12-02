package scraper

import (
	"testing"
)

func TestNewScraper(t *testing.T) {
	s := NewScraper()
	if s == nil {
		t.Fatal("expected scraper to be created")
	}
	if s.cache == nil {
		t.Error("expected cache to be initialized")
	}
	if s.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}
}

func TestGetCountryName(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"US", "United States"},
		{"us", "United States"},
		{"BR", "Brazil"},
		{"CN", "China"},
		{"AU", "Australia"},
		{"XX", "XX"}, // Unknown code should return as-is
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := GetCountryName(tt.code)
			if result != tt.expected {
				t.Errorf("GetCountryName(%s) = %s, expected %s", tt.code, result, tt.expected)
			}
		})
	}
}

func TestGetSpeciesByCountry(t *testing.T) {
	s := NewScraper()

	// Test with known country code
	data, err := s.GetSpeciesByCountry("US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.CountryCode != "US" {
		t.Errorf("expected country code US, got %s", data.CountryCode)
	}

	if data.Country != "United States" {
		t.Errorf("expected country name United States, got %s", data.Country)
	}

	if len(data.Species) == 0 {
		t.Error("expected at least one species")
	}

	// Verify species data structure
	for _, sp := range data.Species {
		if sp.Name == "" {
			t.Error("species name should not be empty")
		}
	}
}

func TestGetSpeciesByCountry_Caching(t *testing.T) {
	s := NewScraper()

	// First call
	data1, err := s.GetSpeciesByCountry("US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second call should return cached data
	data2, err := s.GetSpeciesByCountry("US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both should have the same last updated timestamp (cached)
	if data1.LastUpdated != data2.LastUpdated {
		t.Error("expected cached data to have same timestamp")
	}
}

func TestGetSpeciesByCountry_CaseInsensitive(t *testing.T) {
	s := NewScraper()

	data1, err := s.GetSpeciesByCountry("us")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data2, err := s.GetSpeciesByCountry("US")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data1.CountryCode != data2.CountryCode {
		t.Error("expected same country code regardless of case")
	}
}

func TestGetSampleData(t *testing.T) {
	s := NewScraper()

	// Test known country with sample data
	species := s.getSampleData("US")
	if len(species) == 0 {
		t.Error("expected sample data for US")
	}

	// Verify California Condor is in the list
	found := false
	for _, sp := range species {
		if sp.Name == "California Condor" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected California Condor in US sample data")
	}

	// Test unknown country
	species = s.getSampleData("ZZ")
	if len(species) != 1 {
		t.Error("expected default data for unknown country")
	}
}
