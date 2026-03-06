package supplier

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
)

type Hertz struct {
	r io.Reader
}

func NewHertz() Hertz {
	return Hertz{}
}

func NewHertzWithReader(r io.Reader) Hertz {
	return Hertz{r: r}
}

func (h Hertz) Supplier() SupplierName {
	return SupplierHertz
}

func (h Hertz) GetLocationsPage(cursor string) (LocationPage, error) {
	if h.r == nil {
		return LocationPage{}, errors.New("GetLocationPage required hertz reader to be initialized")
	}

	csvr := csv.NewReader(h.r)
	records, err := csvr.ReadAll()
	if err != nil {
		return LocationPage{}, err
	}

	expectedHeader := []string{"location_id", "country_code", "country", "location_name"}
	if len(records) == 0 || !slices.Equal(records[0], expectedHeader) {
		return LocationPage{}, errors.New("invalid CSV format for Hertz locations")
	}

	locations := make([]Location, 0, len(records)-1)
	for i, record := range records[1:] {
		locationID := record[0]
		countryCode := record[1]
		country := record[2]
		locationName := record[3]

		if locationID == "" || countryCode == "" || country == "" || locationName == "" {
			return LocationPage{}, fmt.Errorf("missing required fields in CSV at line %d", i+2)
		}

		locations = append(locations, Location{
			ID:          locationID,
			CountryCode: countryCode,
			Country:     country,
			Name:        locationName,
		})
	}

	return LocationPage{Locations: locations}, nil
}
