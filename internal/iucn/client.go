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

// =====================================================
// v4 STRUCTS
// =====================================================

// /countries response
type CountryListResponse struct {
	Countries []struct {
		Code        string `json:"code"`
		Description struct {
			En string `json:"en"`
		} `json:"description"`
	} `json:"countries"`
}

// /countries/{code} response
type CountryResponse struct {
	Country struct {
		Code        string `json:"code"`
		Description struct {
			En string `json:"en"`
		} `json:"description"`
	} `json:"country"`

	Assessments []Assessment `json:"assessments"`
}

type Assessment struct {
	TaxonScientificName string `json:"taxon_scientific_name"`
	RedListCategoryCode string `json:"red_list_category_code"`
	Url                 string `json:"url"`
}

// /taxa/scientific_name response
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

// frontend format
type Species struct {
	Name           string `json:"name"`
	ScientificName string `json:"scientific_name"`
	Status         string `json:"status"`
	Kingdom        string `json:"kingdom"`
	Phylum         string `json:"phylum"`
	Class          string `json:"class"`
	Order          string `json:"order"`
	Family         string `json:"family"`
	Url            string `json:"url"`
}

type CountryData struct {
	Country     string    `json:"country"`
	CountryCode string    `json:"country_code"`
	Species     []Species `json:"species"`
}

// =====================================================
// CLIENT INITIALIZATION
// =====================================================

type Client struct {
	Token             string
	HTTPClient        *http.Client
	CountryCodeToName map[string]string
}

func NewClient(token string) *Client {
	c := &Client{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		CountryCodeToName: make(map[string]string),
	}

	// async background loading
	go c.FetchAllCountries()
	return c
}

// =====================================================
// FETCH ALL COUNTRIES (v4)
// =====================================================

func (c *Client) FetchAllCountries() error {
	url := fmt.Sprintf("%s/countries?token=%s", baseURL, c.Token)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed country list status: %d", resp.StatusCode)
	}

	var data CountryListResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	for _, ct := range data.Countries {
		c.CountryCodeToName[ct.Code] = ct.Description.En
	}

	return nil
}

// =====================================================
// GET SPECIES FOR COUNTRY (v4)
// =====================================================

func (c *Client) GetSpeciesByCountry(countryCode string) (*CountryData, error) {
	countryCode = strings.ToUpper(countryCode)

	// =====================================================
	// DUMMY DATA FALLBACK (Requested by User)
	// =====================================================
	if data := getDummyData(countryCode); data != nil {
		return data, nil
	}

	// API Logic (Commented out or bypassed for these countries, but kept for others if needed)
	// For now, we will proceed with API calls for other countries, or you can comment this out entirely.

	url := fmt.Sprintf("%s/countries/%s?token=%s", baseURL, countryCode, c.Token)
	fmt.Printf("DEBUG: Fetching URL: %s\n", url)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("DEBUG: API Status Code: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var data CountryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}

	// ============================
	// BUILD SPECIES LIST
	// ============================
	var speciesList []Species
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Semaphore to limit the number of concurrent API requests to 10.
	// This prevents us from overwhelming the IUCN API and getting rate-limited.
	sem := make(chan struct{}, 10)

	for _, a := range data.Assessments {
		// We only care about species that are actually endangered.
		// Categories: CR (Critically Endangered), EN (Endangered), VU (Vulnerable)
		if !isEndangered(a.RedListCategoryCode) {
			continue
		}

		// Thread-safe check to limit the total number of species we display to 20.
		// This ensures the page loads reasonably fast.
		mu.Lock()
		if len(speciesList) >= 20 {
			mu.Unlock()
			break
		}
		// Pre-allocate a slot in the slice
		speciesList = append(speciesList, Species{})
		idx := len(speciesList) - 1
		mu.Unlock()

		wg.Add(1)
		// Launch a goroutine to fetch details for this species in the background
		go func(idx int, a Assessment) {
			defer wg.Done()

			// Acquire token from semaphore (blocks if 10 requests are already running)
			sem <- struct{}{}
			defer func() { <-sem }() // Release token when done

			// Fetch detailed taxonomy (Kingdom, Phylum, etc.)
			details, _ := c.GetSpeciesDetails(a.TaxonScientificName)

			// Update the species list in a thread-safe way
			mu.Lock()
			defer mu.Unlock()

			s := Species{
				ScientificName: a.TaxonScientificName,
				Status:         mapCategory(a.RedListCategoryCode),
				Url:            a.Url,
			}

			if details != nil {
				s.Kingdom = details.KingdomName
				s.Phylum = details.PhylumName
				s.Class = details.ClassName
				s.Order = details.OrderName
				s.Family = details.FamilyName

				// Iterate through common names to find the English one
				for _, cn := range details.CommonNames {
					if cn.Language == "eng" {
						s.Name = cn.Name
						break
					}
				}
			}
			// Fallback: use scientific name if no common name is found
			if s.Name == "" {
				s.Name = s.ScientificName
			}

			speciesList[idx] = s
		}(idx, a)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	name := data.Country.Description.En
	if name == "" {
		name = c.GetCountryName(countryCode)
	}

	return &CountryData{
		Country:     name,
		CountryCode: countryCode,
		Species:     speciesList,
	}, nil
}

// =====================================================
// GET SPECIES DETAILS (v4)
// =====================================================

func (c *Client) GetSpeciesDetails(scientificName string) (*TaxonDetails, error) {
	parts := strings.Split(scientificName, " ")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid scientific name: %s", scientificName)
	}

	genus, species := parts[0], parts[1]

	url := fmt.Sprintf(
		"%s/taxa/scientific_name?genus_name=%s&species_name=%s&token=%s",
		baseURL, genus, species, c.Token,
	)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("taxonomy status %d", resp.StatusCode)
	}

	var tr TaxonResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}

	return &tr.Taxon, nil
}

// =====================================================

func isEndangered(code string) bool {
	switch code {
	case "CR", "EN", "VU":
		return true
	}
	return false
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
	}
	return code
}

func (c *Client) GetCountryName(code string) string {
	if n, ok := c.CountryCodeToName[strings.ToUpper(code)]; ok {
		return n
	}
	return code
}
