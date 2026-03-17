package broker

import (
	"encoding/csv"
	"errors"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrHertzReaderNotInitialized = errors.New("Hertz reader is not initialized")
	ErrHertzInvalidCSVFormat     = errors.New("invalid CSV format for Hertz locations")
	ErrHertzMissingFields        = errors.New("missing required fields in CSV for Hertz locations")
)

// GetLocationsPage retrieves a page of locations from the Hertz broker based on the provided cursor, which is not used in this implementation since all locations are returned in a single page.
func (h Hertz) GetLocationsPage(cursor string) (LocationPage, error) {
	if h.r == nil {
		return LocationPage{}, ErrHertzReaderNotInitialized
	}

	csvr := csv.NewReader(h.r)
	records, err := csvr.ReadAll()
	if err != nil {
		return LocationPage{}, fmt.Errorf("error reading CSV data for Hertz locations: %w", err)
	}

	expectedHeader := []string{"location_id", "country_code", "country", "location_name"}
	if len(records) == 0 || !slices.Equal(records[0], expectedHeader) {
		return LocationPage{}, ErrHertzInvalidCSVFormat
	}

	locations := make([]Location, 0, len(records)-1)
	for i, record := range records[1:] {
		locationID := record[0]
		countryCode := record[1]
		country := record[2]
		locationName := record[3]

		if locationID == "" || countryCode == "" || country == "" || locationName == "" {
			return LocationPage{}, fmt.Errorf("%w at line %d", ErrHertzMissingFields, i+2)
		}

		iata := ""
		if len(locationID) == 3 {
			iata = locationID
		}

		locations = append(locations, Location{
			ID:          locationID,
			CountryCode: countryCode,
			Country:     country,
			Iata:        iata,
			Name:        normalizeHertzLocationName(locationName),
		})
	}

	return LocationPage{Locations: locations}, nil
}

// normalizeHertzLocationName normalizes the location name from the Hertz CSV by trimming whitespace, converting to lowercase, and capitalizing each word.
func normalizeHertzLocationName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	words := strings.Fields(name)
	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, " ")
}
