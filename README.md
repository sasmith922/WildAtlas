# WildAtlas

Scope Fall 2025 Curriculum Project

## Overview

WildAtlas is an interactive web application that displays endangered species information for countries around the world. Click on any country on the map to view information about its endangered wildlife.

## Features

- 🗺️ Interactive world map using Leaflet.js
- 🦁 Endangered species data for multiple countries
- 🔍 Web scraping capability for Wikipedia endangered species lists
- 💾 Server-side caching for improved performance
- 📱 Responsive design for all screen sizes

## Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript with Leaflet.js
- **Web Scraping**: goquery library

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/sasmith922/WildAtlas.git
cd WildAtlas
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
go build -o wildatlas ./cmd/server
```

4. Run the server:
```bash
./wildatlas
```

5. Open your browser and navigate to `http://localhost:8080`

### Environment Variables

- `PORT`: Server port (default: 8080)

## API Endpoints

### Get Species by Country
```
GET /api/species/{countryCode}
```

Returns endangered species data for the specified country code (ISO 3166-1 alpha-2).

**Example:**
```bash
curl http://localhost:8080/api/species/US
```

**Response:**
```json
{
  "country": "United States",
  "country_code": "US",
  "species": [
    {
      "name": "California Condor",
      "scientific_name": "Gymnogyps californianus",
      "status": "Critically Endangered",
      "population": "~500",
      "habitat": "Mountains and forests",
      "threats": "Habitat loss, lead poisoning"
    }
  ],
  "last_updated": "2025-12-02T00:00:00Z"
}
```

### Health Check
```
GET /api/health
```

Returns server health status.

## Project Structure

```
WildAtlas/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── handlers.go      # HTTP request handlers
│   │   └── handlers_test.go # Handler tests
│   └── scraper/
│       ├── scraper.go       # Web scraping logic
│       └── scraper_test.go  # Scraper tests
├── static/
│   └── index.html           # Frontend application
├── go.mod
├── go.sum
└── README.md
```

## Running Tests

```bash
go test -v ./...
```

## License

This project is for educational purposes as part of the Scope Fall 2025 Curriculum.
