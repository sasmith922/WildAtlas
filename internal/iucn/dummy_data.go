package iucn

// getDummyData returns dummy data for specific countries if API is failing.
func getDummyData(code string) *CountryData {
	switch code {
	case "CA": // Canada
		return &CountryData{
			Country:     "Canada",
			CountryCode: "CA",
			Species: []Species{
				{
					Name:           "Polar Bear",
					ScientificName: "Ursus maritimus",
					Status:         "Vulnerable",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Mammalia",
					Order:          "Carnivora",
					Family:         "Ursidae",
					Url:            "https://www.iucnredlist.org/species/22823/14871490",
				},
				{
					Name:           "Vancouver Island Marmot",
					ScientificName: "Marmota vancouverensis",
					Status:         "Critically Endangered",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Mammalia",
					Order:          "Rodentia",
					Family:         "Sciuridae",
					Url:            "https://www.iucnredlist.org/species/12828/111561606",
				},
				{
					Name:           "Whooping Crane",
					ScientificName: "Grus americana",
					Status:         "Endangered",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Aves",
					Order:          "Gruiformes",
					Family:         "Gruidae",
					Url:            "https://www.iucnredlist.org/species/22692156/111562000",
				},
			},
		}
	case "BR": // Brazil
		return &CountryData{
			Country:     "Brazil",
			CountryCode: "BR",
			Species: []Species{
				{
					Name:           "Golden Lion Tamarin",
					ScientificName: "Leontopithecus rosalia",
					Status:         "Endangered",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Mammalia",
					Order:          "Primates",
					Family:         "Callitrichidae",
					Url:            "https://www.iucnredlist.org/species/11506/192319267",
				},
				{
					Name:           "Jaguar",
					ScientificName: "Panthera onca",
					Status:         "Near Threatened",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Mammalia",
					Order:          "Carnivora",
					Family:         "Felidae",
					Url:            "https://www.iucnredlist.org/species/15953/123791436",
				},
				{
					Name:           "Hyacinth Macaw",
					ScientificName: "Anodorhynchus hyacinthinus",
					Status:         "Vulnerable",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Aves",
					Order:          "Psittaciformes",
					Family:         "Psittacidae",
					Url:            "https://www.iucnredlist.org/species/22685516/93077457",
				},
			},
		}
	case "AU": // Australia
		return &CountryData{
			Country:     "Australia",
			CountryCode: "AU",
			Species: []Species{
				{
					Name:           "Koala",
					ScientificName: "Phascolarctos cinereus",
					Status:         "Vulnerable",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Mammalia",
					Order:          "Diprotodontia",
					Family:         "Phascolarctidae",
					Url:            "https://www.iucnredlist.org/species/16892/166496779",
				},
				{
					Name:           "Tasmanian Devil",
					ScientificName: "Sarcophilus harrisii",
					Status:         "Endangered",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Mammalia",
					Order:          "Dasyuromorphia",
					Family:         "Dasyuridae",
					Url:            "https://www.iucnredlist.org/species/40540/10331066",
				},
				{
					Name:           "Regent Honeyeater",
					ScientificName: "Anthochaera phrygia",
					Status:         "Critically Endangered",
					Kingdom:        "Animalia",
					Phylum:         "Chordata",
					Class:          "Aves",
					Order:          "Passeriformes",
					Family:         "Meliphagidae",
					Url:            "https://www.iucnredlist.org/species/22704415/219632355",
				},
			},
		}
	default:
		return nil
	}
}
