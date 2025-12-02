// Initialize the map
const map = L.map('map').setView([20, 0], 2);

// Add tile layer (using a dark theme to match our design)
L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
    subdomains: 'abcd',
    maxZoom: 19
}).addTo(map);

// Country boundaries GeoJSON URL
const countriesGeoJsonUrl = 'https://raw.githubusercontent.com/datasets/geo-countries/master/data/countries.geojson';

let geojsonLayer;
let selectedCountry = null;

// Style for countries
function style(feature) {
    return {
        fillColor: '#0f3460',
        weight: 1,
        opacity: 1,
        color: '#16213e',
        fillOpacity: 0.7
    };
}

// Highlight style on hover
function highlightStyle(feature) {
    return {
        fillColor: '#e94560',
        weight: 2,
        opacity: 1,
        color: '#fff',
        fillOpacity: 0.7
    };
}

// Selected style
function selectedStyle(feature) {
    return {
        fillColor: '#e94560',
        weight: 3,
        opacity: 1,
        color: '#fff',
        fillOpacity: 0.8
    };
}

// Reset highlight
function resetHighlight(e) {
    if (e.target !== selectedCountry) {
        geojsonLayer.resetStyle(e.target);
    }
}

// Highlight on hover
function highlightFeature(e) {
    const layer = e.target;
    if (layer !== selectedCountry) {
        layer.setStyle(highlightStyle());
    }
    layer.bringToFront();
}

// Handle country click
function onCountryClick(e) {
    const layer = e.target;
    const properties = layer.feature.properties;
    const countryCode = properties.ISO_A2;
    const countryName = properties.ADMIN;

    // Reset previous selection
    if (selectedCountry) {
        geojsonLayer.resetStyle(selectedCountry);
    }

    // Set new selection
    selectedCountry = layer;
    layer.setStyle(selectedStyle());

    // Zoom to country
    map.fitBounds(layer.getBounds(), { padding: [50, 50] });

    // Fetch species data
    fetchSpeciesData(countryCode, countryName);
}

// Add interactions to each country
function onEachFeature(feature, layer) {
    layer.on({
        mouseover: highlightFeature,
        mouseout: resetHighlight,
        click: onCountryClick
    });
}

// Fetch and display species data
async function fetchSpeciesData(countryCode, countryName) {
    const sidebar = document.getElementById('sidebar-content');

    // Show loading state
    sidebar.innerHTML = `
        <div class="sidebar-header">
            <h2>Endangered Species</h2>
        </div>
        <div class="loading">
            <div class="loading-spinner"></div>
            <p>Loading species data...</p>
        </div>
    `;

    try {
        const response = await fetch(`/api/species/${countryCode}`);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        displaySpeciesData(data);
    } catch (error) {
        console.error('Error fetching species data:', error);
        sidebar.innerHTML = `
            <div class="sidebar-header">
                <h2>Endangered Species</h2>
            </div>
            <div class="error-message">
                <p>Failed to load species data for ${countryName}</p>
                <p>Please try again later.</p>
            </div>
        `;
    }
}

// Display species data in sidebar
function displaySpeciesData(data) {
    const sidebar = document.getElementById('sidebar-content');

    let speciesHtml = data.species.map((species, index) => {
        const statusClass = getStatusClass(species.status);
        const uniqueId = `species-${index}`;

        return `
            <div class="species-card" onclick="toggleDetails('${uniqueId}')">
                <div class="species-header">
                    <div>
                        <h4 class="species-name">${escapeHtml(species.name)}</h4>
                        <p class="scientific-name">${escapeHtml(species.scientific_name)}</p>
                    </div>
                    <span class="status-badge ${statusClass}">${escapeHtml(species.status)}</span>
                </div>
                <div id="${uniqueId}" class="species-details hidden">
                    <div class="detail-grid">
                        <div class="detail-item">
                            <span class="label">Kingdom:</span>
                            <span class="value">${escapeHtml(species.kingdom || 'N/A')}</span>
                        </div>
                        <div class="detail-item">
                            <span class="label">Phylum:</span>
                            <span class="value">${escapeHtml(species.phylum || 'N/A')}</span>
                        </div>
                        <div class="detail-item">
                            <span class="label">Class:</span>
                            <span class="value">${escapeHtml(species.class || 'N/A')}</span>
                        </div>
                        <div class="detail-item">
                            <span class="label">Order:</span>
                            <span class="value">${escapeHtml(species.order || 'N/A')}</span>
                        </div>
                        <div class="detail-item">
                            <span class="label">Family:</span>
                            <span class="value">${escapeHtml(species.family || 'N/A')}</span>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }).join('');

    sidebar.innerHTML = `
        <div class="sidebar-header">
            <h2>Endangered Species</h2>
        </div>
        <div class="country-info">
            <h3 class="country-name">${escapeHtml(data.country)}</h3>
            <p class="country-code">Country Code: ${escapeHtml(data.country_code)}</p>
        </div>
        <div class="species-list">
            ${speciesHtml}
        </div>
    `;
}

function toggleDetails(id) {
    const element = document.getElementById(id);
    if (element.classList.contains('hidden')) {
        element.classList.remove('hidden');
    } else {
        element.classList.add('hidden');
    }
}

// Get CSS class for status badge
function getStatusClass(status) {
    const lowerStatus = status.toLowerCase();
    if (lowerStatus.includes('critically')) return 'status-critically-endangered';
    if (lowerStatus.includes('endangered')) return 'status-endangered';
    if (lowerStatus.includes('vulnerable')) return 'status-vulnerable';
    if (lowerStatus.includes('near threatened')) return 'status-near-threatened';
    return 'status-default';
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Load countries GeoJSON
fetch(countriesGeoJsonUrl)
    .then(response => response.json())
    .then(data => {
        geojsonLayer = L.geoJSON(data, {
            style: style,
            onEachFeature: onEachFeature
        }).addTo(map);
    })
    .catch(error => {
        console.error('Error loading countries GeoJSON:', error);
    });
