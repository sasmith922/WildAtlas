package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Species represents an endangered species
type Species struct {
	Name           string `json:"name"`
	ScientificName string `json:"scientific_name"`
	Status         string `json:"status"`
	Population     string `json:"population"`
	Habitat        string `json:"habitat"`
	Threats        string `json:"threats"`
}

// CountryData represents endangered species data for a country
type CountryData struct {
	Country     string    `json:"country"`
	CountryCode string    `json:"country_code"`
	Species     []Species `json:"species"`
	LastUpdated string    `json:"last_updated"`
}

// Scraper handles web scraping for endangered species data
type Scraper struct {
	cache      map[string]*CountryData
	cacheMutex sync.RWMutex
	httpClient *http.Client
}

// NewScraper creates a new Scraper instance
func NewScraper() *Scraper {
	return &Scraper{
		cache: make(map[string]*CountryData),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// countryCodeToName maps ISO 3166-1 alpha-2 codes to country names
var countryCodeToName = map[string]string{
	"US": "United States",
	"CA": "Canada",
	"MX": "Mexico",
	"BR": "Brazil",
	"AR": "Argentina",
	"GB": "United Kingdom",
	"FR": "France",
	"DE": "Germany",
	"IT": "Italy",
	"ES": "Spain",
	"PT": "Portugal",
	"CN": "China",
	"JP": "Japan",
	"IN": "India",
	"AU": "Australia",
	"NZ": "New Zealand",
	"ZA": "South Africa",
	"KE": "Kenya",
	"TZ": "Tanzania",
	"EG": "Egypt",
	"NG": "Nigeria",
	"RU": "Russia",
	"ID": "Indonesia",
	"MY": "Malaysia",
	"TH": "Thailand",
	"VN": "Vietnam",
	"PH": "Philippines",
	"KR": "South Korea",
	"PK": "Pakistan",
	"BD": "Bangladesh",
	"CO": "Colombia",
	"PE": "Peru",
	"VE": "Venezuela",
	"CL": "Chile",
	"EC": "Ecuador",
	"BO": "Bolivia",
	"PY": "Paraguay",
	"UY": "Uruguay",
	"CR": "Costa Rica",
	"PA": "Panama",
	"CU": "Cuba",
	"GT": "Guatemala",
	"HN": "Honduras",
	"NI": "Nicaragua",
	"SV": "El Salvador",
	"BZ": "Belize",
	"MG": "Madagascar",
	"MW": "Malawi",
	"ZM": "Zambia",
	"ZW": "Zimbabwe",
	"BW": "Botswana",
	"NA": "Namibia",
	"AO": "Angola",
	"MZ": "Mozambique",
	"CD": "Democratic Republic of the Congo",
	"CG": "Republic of the Congo",
	"GA": "Gabon",
	"CM": "Cameroon",
	"GH": "Ghana",
	"CI": "Ivory Coast",
	"SN": "Senegal",
	"ML": "Mali",
	"NE": "Niger",
	"TD": "Chad",
	"SD": "Sudan",
	"ET": "Ethiopia",
	"SO": "Somalia",
	"UG": "Uganda",
	"RW": "Rwanda",
	"BI": "Burundi",
	"NP": "Nepal",
	"BT": "Bhutan",
	"LK": "Sri Lanka",
	"MM": "Myanmar",
	"LA": "Laos",
	"KH": "Cambodia",
	"SG": "Singapore",
	"BN": "Brunei",
	"PG": "Papua New Guinea",
	"FJ": "Fiji",
	"NO": "Norway",
	"SE": "Sweden",
	"FI": "Finland",
	"DK": "Denmark",
	"IS": "Iceland",
	"IE": "Ireland",
	"NL": "Netherlands",
	"BE": "Belgium",
	"LU": "Luxembourg",
	"CH": "Switzerland",
	"AT": "Austria",
	"PL": "Poland",
	"CZ": "Czech Republic",
	"SK": "Slovakia",
	"HU": "Hungary",
	"RO": "Romania",
	"BG": "Bulgaria",
	"GR": "Greece",
	"TR": "Turkey",
	"UA": "Ukraine",
	"BY": "Belarus",
	"RS": "Serbia",
	"HR": "Croatia",
	"SI": "Slovenia",
	"BA": "Bosnia and Herzegovina",
	"MK": "North Macedonia",
	"AL": "Albania",
	"ME": "Montenegro",
	"XK": "Kosovo",
	"MD": "Moldova",
	"EE": "Estonia",
	"LV": "Latvia",
	"LT": "Lithuania",
}

// GetCountryName returns the country name for a given ISO code
func GetCountryName(code string) string {
	if name, ok := countryCodeToName[strings.ToUpper(code)]; ok {
		return name
	}
	return code
}

// GetSpeciesByCountry fetches endangered species data for a given country
func (s *Scraper) GetSpeciesByCountry(countryCode string) (*CountryData, error) {
	countryCode = strings.ToUpper(countryCode)

	// Check cache first
	s.cacheMutex.RLock()
	if data, ok := s.cache[countryCode]; ok {
		s.cacheMutex.RUnlock()
		return data, nil
	}
	s.cacheMutex.RUnlock()

	// Scrape data from worldwildlife.org or similar source
	data, err := s.scrapeSpeciesData(countryCode)
	if err != nil {
		return nil, err
	}

	// Update cache
	s.cacheMutex.Lock()
	s.cache[countryCode] = data
	s.cacheMutex.Unlock()

	return data, nil
}

// scrapeSpeciesData scrapes endangered species data from the web
func (s *Scraper) scrapeSpeciesData(countryCode string) (*CountryData, error) {
	countryName := GetCountryName(countryCode)

	// Try to scrape from Wikipedia's endangered species by country
	species, err := s.scrapeWikipedia(countryName)
	if err != nil {
		// If scraping fails, return sample data as fallback
		species = s.getSampleData(countryCode)
	}

	return &CountryData{
		Country:     countryName,
		CountryCode: countryCode,
		Species:     species,
		LastUpdated: time.Now().Format(time.RFC3339),
	}, nil
}

// scrapeWikipedia attempts to scrape endangered species from Wikipedia
func (s *Scraper) scrapeWikipedia(countryName string) ([]Species, error) {
	// Build Wikipedia URL for endangered species
	searchName := strings.ReplaceAll(countryName, " ", "_")
	url := fmt.Sprintf("https://en.wikipedia.org/wiki/List_of_endangered_species_in_%s", searchName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WildAtlas/1.0 (Educational Project)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var species []Species

	// Parse tables for species information
	doc.Find("table.wikitable tbody tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return // Skip header row
		}

		cells := row.Find("td")
		if cells.Length() >= 2 {
			name := strings.TrimSpace(cells.Eq(0).Text())
			scientificName := strings.TrimSpace(cells.Eq(1).Text())

			if name != "" && len(species) < 20 {
				sp := Species{
					Name:           name,
					ScientificName: scientificName,
					Status:         "Endangered",
				}

				if cells.Length() >= 3 {
					sp.Status = strings.TrimSpace(cells.Eq(2).Text())
				}

				species = append(species, sp)
			}
		}
	})

	if len(species) == 0 {
		return nil, fmt.Errorf("no species found")
	}

	return species, nil
}

// getSampleData returns sample endangered species data for demonstration
func (s *Scraper) getSampleData(countryCode string) []Species {
	// Sample data for various regions
	sampleDataByRegion := map[string][]Species{
		"US": {
			{Name: "California Condor", ScientificName: "Gymnogyps californianus", Status: "Critically Endangered", Population: "~500", Habitat: "Mountains and forests", Threats: "Habitat loss, lead poisoning"},
			{Name: "Florida Panther", ScientificName: "Puma concolor coryi", Status: "Endangered", Population: "~200", Habitat: "Swamps and forests", Threats: "Habitat fragmentation, vehicle collisions"},
			{Name: "Hawaiian Monk Seal", ScientificName: "Neomonachus schauinslandi", Status: "Endangered", Population: "~1,400", Habitat: "Hawaiian Islands", Threats: "Climate change, marine debris"},
			{Name: "Red Wolf", ScientificName: "Canis rufus", Status: "Critically Endangered", Population: "~20 in wild", Habitat: "Forests and wetlands", Threats: "Hybridization, habitat loss"},
			{Name: "Ocelot", ScientificName: "Leopardus pardalis", Status: "Endangered", Population: "~100 in US", Habitat: "Dense thorny scrubland", Threats: "Habitat loss, vehicle strikes"},
		},
		"BR": {
			{Name: "Golden Lion Tamarin", ScientificName: "Leontopithecus rosalia", Status: "Endangered", Population: "~3,200", Habitat: "Atlantic Forest", Threats: "Deforestation, illegal pet trade"},
			{Name: "Hyacinth Macaw", ScientificName: "Anodorhynchus hyacinthinus", Status: "Vulnerable", Population: "~6,500", Habitat: "Pantanal wetlands", Threats: "Illegal trade, habitat loss"},
			{Name: "Jaguar", ScientificName: "Panthera onca", Status: "Near Threatened", Population: "~170,000", Habitat: "Rainforests and wetlands", Threats: "Deforestation, poaching"},
			{Name: "Amazon River Dolphin", ScientificName: "Inia geoffrensis", Status: "Endangered", Population: "Unknown", Habitat: "Amazon River system", Threats: "Dam construction, pollution"},
			{Name: "Black Lion Tamarin", ScientificName: "Leontopithecus chrysopygus", Status: "Endangered", Population: "~1,000", Habitat: "Atlantic Forest", Threats: "Habitat fragmentation"},
		},
		"CN": {
			{Name: "Giant Panda", ScientificName: "Ailuropoda melanoleuca", Status: "Vulnerable", Population: "~1,800", Habitat: "Mountain bamboo forests", Threats: "Habitat loss, low birth rate"},
			{Name: "South China Tiger", ScientificName: "Panthera tigris amoyensis", Status: "Critically Endangered", Population: "~0 in wild", Habitat: "Temperate forests", Threats: "Poaching, habitat loss"},
			{Name: "Chinese Alligator", ScientificName: "Alligator sinensis", Status: "Critically Endangered", Population: "~150 in wild", Habitat: "Yangtze River wetlands", Threats: "Habitat destruction, pollution"},
			{Name: "Yangtze Finless Porpoise", ScientificName: "Neophocaena asiaeorientalis", Status: "Critically Endangered", Population: "~1,000", Habitat: "Yangtze River", Threats: "Pollution, boat traffic"},
			{Name: "Crested Ibis", ScientificName: "Nipponia nippon", Status: "Endangered", Population: "~2,600", Habitat: "Wetlands and rice paddies", Threats: "Habitat loss, pesticides"},
		},
		"IN": {
			{Name: "Bengal Tiger", ScientificName: "Panthera tigris tigris", Status: "Endangered", Population: "~3,000", Habitat: "Forests and grasslands", Threats: "Poaching, habitat loss"},
			{Name: "Asian Elephant", ScientificName: "Elephas maximus", Status: "Endangered", Population: "~27,000 in India", Habitat: "Forests and grasslands", Threats: "Habitat fragmentation, human conflict"},
			{Name: "Indian Rhinoceros", ScientificName: "Rhinoceros unicornis", Status: "Vulnerable", Population: "~3,700", Habitat: "Grasslands and riverine areas", Threats: "Poaching, habitat loss"},
			{Name: "Ganges River Dolphin", ScientificName: "Platanista gangetica", Status: "Endangered", Population: "~3,500", Habitat: "Ganges River system", Threats: "Pollution, dam construction"},
			{Name: "Snow Leopard", ScientificName: "Panthera uncia", Status: "Vulnerable", Population: "~500 in India", Habitat: "High mountain regions", Threats: "Poaching, climate change"},
		},
		"AU": {
			{Name: "Koala", ScientificName: "Phascolarctos cinereus", Status: "Vulnerable", Population: "~100,000", Habitat: "Eucalyptus forests", Threats: "Habitat loss, disease, bushfires"},
			{Name: "Numbat", ScientificName: "Myrmecobius fasciatus", Status: "Endangered", Population: "~1,000", Habitat: "Eucalyptus woodlands", Threats: "Predation by foxes and cats"},
			{Name: "Leadbeater's Possum", ScientificName: "Gymnobelideus leadbeateri", Status: "Critically Endangered", Population: "~1,500", Habitat: "Mountain ash forests", Threats: "Logging, bushfires"},
			{Name: "Northern Hairy-nosed Wombat", ScientificName: "Lasiorhinus krefftii", Status: "Critically Endangered", Population: "~300", Habitat: "Semi-arid grasslands", Threats: "Competition with cattle, drought"},
			{Name: "Tasmanian Devil", ScientificName: "Sarcophilus harrisii", Status: "Endangered", Population: "~25,000", Habitat: "Tasmanian forests", Threats: "Devil facial tumour disease"},
		},
		"KE": {
			{Name: "Black Rhinoceros", ScientificName: "Diceros bicornis", Status: "Critically Endangered", Population: "~750 in Kenya", Habitat: "Savannas and forests", Threats: "Poaching for horn"},
			{Name: "African Wild Dog", ScientificName: "Lycaon pictus", Status: "Endangered", Population: "~600 in Kenya", Habitat: "Savannas and grasslands", Threats: "Habitat fragmentation, human conflict"},
			{Name: "Grevy's Zebra", ScientificName: "Equus grevyi", Status: "Endangered", Population: "~2,800", Habitat: "Semi-arid grasslands", Threats: "Habitat loss, competition with livestock"},
			{Name: "Hirola", ScientificName: "Beatragus hunteri", Status: "Critically Endangered", Population: "~500", Habitat: "Semi-arid grasslands", Threats: "Drought, habitat loss, disease"},
			{Name: "Mountain Bongo", ScientificName: "Tragelaphus eurycerus isaaci", Status: "Critically Endangered", Population: "~100 in wild", Habitat: "Mountain forests", Threats: "Poaching, habitat loss"},
		},
		"MG": {
			{Name: "Aye-aye", ScientificName: "Daubentonia madagascariensis", Status: "Endangered", Population: "Unknown", Habitat: "Rainforests", Threats: "Deforestation, persecution"},
			{Name: "Indri", ScientificName: "Indri indri", Status: "Critically Endangered", Population: "~10,000", Habitat: "Rainforests", Threats: "Habitat loss, hunting"},
			{Name: "Silky Sifaka", ScientificName: "Propithecus candidus", Status: "Critically Endangered", Population: "~250", Habitat: "Mountain rainforests", Threats: "Habitat loss, hunting"},
			{Name: "Radiated Tortoise", ScientificName: "Astrochelys radiata", Status: "Critically Endangered", Population: "Unknown", Habitat: "Spiny forests", Threats: "Illegal pet trade, habitat loss"},
			{Name: "Ploughshare Tortoise", ScientificName: "Astrochelys yniphora", Status: "Critically Endangered", Population: "~500", Habitat: "Bamboo scrub", Threats: "Illegal pet trade"},
		},
		"ID": {
			{Name: "Sumatran Tiger", ScientificName: "Panthera tigris sumatrae", Status: "Critically Endangered", Population: "~400", Habitat: "Tropical rainforests", Threats: "Poaching, deforestation"},
			{Name: "Sumatran Orangutan", ScientificName: "Pongo abelii", Status: "Critically Endangered", Population: "~14,000", Habitat: "Tropical rainforests", Threats: "Habitat loss, illegal trade"},
			{Name: "Javan Rhinoceros", ScientificName: "Rhinoceros sondaicus", Status: "Critically Endangered", Population: "~70", Habitat: "Tropical rainforests", Threats: "Poaching, habitat loss"},
			{Name: "Sumatran Rhinoceros", ScientificName: "Dicerorhinus sumatrensis", Status: "Critically Endangered", Population: "~80", Habitat: "Tropical rainforests", Threats: "Poaching, habitat loss"},
			{Name: "Komodo Dragon", ScientificName: "Varanus komodoensis", Status: "Endangered", Population: "~3,000", Habitat: "Islands of Indonesia", Threats: "Habitat loss, climate change"},
		},
	}

	// Return specific data if available, otherwise return generic endangered species
	if data, ok := sampleDataByRegion[countryCode]; ok {
		return data
	}

	// Default data for countries without specific information
	return []Species{
		{Name: "Data unavailable", ScientificName: "N/A", Status: "Please check IUCN Red List", Population: "Unknown", Habitat: "Various", Threats: "Multiple factors"},
	}
}
