package supplier

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
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

	if len(records) == 0 || len(records[0]) != 3 || records[0][0] != "location_id" || records[0][1] != "country_code" || records[0][2] != "location_name" {
		return LocationPage{}, errors.New("invalid CSV format for Hertz locations")
	}

	locations := make([]Location, 0, len(records)-1)
	for i, record := range records[1:] {
		locationID := record[0]
		countryCode := record[1]
		locationName := record[2]

		if locationID == "" || countryCode == "" || locationName == "" {
			return LocationPage{}, fmt.Errorf("missing required fields in CSV at line %d", i+2)
		}

		locations = append(locations, Location{
			ID:          locationID,
			CountryCode: countryCode,
			Name:        locationName,
		})
	}

	return LocationPage{Locations: locations}, nil
}
