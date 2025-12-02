package iucn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const baseURL = "https://api.iucnredlist.org/api/v4"

// Client handles interactions with the IUCN Red List API.
type Client struct {
	Token             string
	HTTPClient        *http.Client
	CountryCodeToName map[string]string
}

// NewClient creates a new IUCN API client.
func NewClient(token string) *Client {
	c := &Client{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		CountryCodeToName: make(map[string]string),
	}
	// Best effort to populate country names at startup
	go c.FetchAllCountries()
	return c
}

// Structs for v4 API

type CountryListResponse struct {
	Countries []struct {
		Description struct {
			En string `json:"en"`
		} `json:"description"`
		Code string `json:"code"`
	} `json:"countries"`
}

type CountryResponse struct {
	Country struct {
		Description struct {
			En string `json:"en"`
		} `json:"description"`
		Code string `json:"code"`
	} `json:"country"`
	Assessments []Assessment `json:"assessments"`
}

type Assessment struct {
	TaxonScientificName string `json:"taxon_scientific_name"`
	RedListCategoryCode string `json:"red_list_category_code"`
}

type TaxonResponse struct {
	Taxon TaxonDetails `json:"taxon"`
}

type TaxonDetails struct {
	ScientificName string       `json:"scientific_name"`
	KingdomName    string       `json:"kingdom_name"`
	PhylumName     string       `json:"phylum_name"`
	ClassName      string       `json:"class_name"`
	OrderName      string       `json:"order_name"`
	FamilyName     string       `json:"family_name"`
	CommonNames    []CommonName `json:"common_names"`
}

type CommonName struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Main     bool   `json:"main"`
}

// Species represents a species returned by the IUCN API.
type Species struct {
	Name           string `json:"name"`            // Common Name (English)
	ScientificName string `json:"scientific_name"` // Scientific Name
	Status         string `json:"status"`          // Mapped from Code
	Kingdom        string `json:"kingdom"`
	Phylum         string `json:"phylum"`
	Class          string `json:"class"`
	Order          string `json:"order"`
	Family         string `json:"family"`
}

// CountryData represents the aggregated data for a country.
type CountryData struct {
	Country     string    `json:"country"`
	CountryCode string    `json:"country_code"`
	Species     []Species `json:"species"`
}

// FetchAllCountries fetches the list of all countries to build a name mapping.
func (c *Client) FetchAllCountries() error {
	url := fmt.Sprintf("%s/countries?token=%s", baseURL, c.Token)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch countries: %d", resp.StatusCode)
	}

	var listResp CountryListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return err
	}

	for _, country := range listResp.Countries {
		c.CountryCodeToName[country.Code] = country.Description.En
	}
	return nil
}

// GetSpeciesByCountry fetches endangered species for a specific country.
func (c *Client) GetSpeciesByCountry(countryCode string) (*CountryData, error) {
	countryCode = strings.ToUpper(countryCode)

	url := fmt.Sprintf("%s/countries/%s?token=%s", baseURL, countryCode, c.Token)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var countryResp CountryResponse
	if err := json.NewDecoder(resp.Body).Decode(&countryResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var speciesList []Species
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Limit concurrency to avoid overwhelming the API
	sem := make(chan struct{}, 10)

	for _, assessment := range countryResp.Assessments {
		// Filter for endangered categories
		if isEndangered(assessment.RedListCategoryCode) {

			// Limit to 20 species for performance
			mu.Lock()
			if len(speciesList) >= 20 {
				mu.Unlock()
				break
			}
			// Pre-allocate slot
			speciesList = append(speciesList, Species{})
			idx := len(speciesList) - 1
			mu.Unlock()

			wg.Add(1)
			go func(idx int, a Assessment) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire semaphore
				defer func() { <-sem }() // Release semaphore

				details, err := c.GetSpeciesDetails(a.TaxonScientificName)

				mu.Lock()
				defer mu.Unlock()

				s := Species{
					ScientificName: a.TaxonScientificName,
					Status:         mapCategory(a.RedListCategoryCode),
				}

				if err == nil && details != nil {
					s.Kingdom = details.KingdomName
					s.Phylum = details.PhylumName
					s.Class = details.ClassName
					s.Order = details.OrderName
					s.Family = details.FamilyName

					// Find English common name
					for _, cn := range details.CommonNames {
						if cn.Language == "eng" {
							s.Name = cn.Name
							if cn.Main {
								break
							}
						}
					}
				}

				// Fallback if no common name found
				if s.Name == "" {
					s.Name = s.ScientificName
				}

				speciesList[idx] = s
			}(idx, assessment)
		}
	}

	wg.Wait()

	countryName := countryResp.Country.Description.En
	if countryName == "" {
		countryName = c.GetCountryName(countryCode)
	}

	return &CountryData{
		Country:     countryName,
		CountryCode: countryCode,
		Species:     speciesList,
	}, nil
}

// GetSpeciesDetails fetches detailed taxonomy information for a species.
func (c *Client) GetSpeciesDetails(scientificName string) (*TaxonDetails, error) {
	parts := strings.Split(scientificName, " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid scientific name: %s", scientificName)
	}
	genus := parts[0]
	species := parts[1]

	url := fmt.Sprintf("%s/taxa/scientific_name?genus_name=%s&species_name=%s&token=%s", baseURL, genus, species, c.Token)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	var taxonResp TaxonResponse
	if err := json.NewDecoder(resp.Body).Decode(&taxonResp); err != nil {
		return nil, err
	}

	return &taxonResp.Taxon, nil
}

func isEndangered(code string) bool {
	switch code {
	case "CR", "EN", "VU":
		return true
	default:
		return false
	}
}

func mapCategory(code string) string {
	switch code {
	case "CR":
		return "Critically Endangered"
	case "EN":
		return "Endangered"
	case "VU":
		return "Vulnerable"
	case "NT":
		return "Near Threatened"
	default:
		return code
	}
}

func (c *Client) GetCountryName(code string) string {
	if name, ok := c.CountryCodeToName[strings.ToUpper(code)]; ok {
		return name
	}
	return code
}
