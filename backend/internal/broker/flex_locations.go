package broker

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
)

// GetLocationsPage retrieves a page of locations from the Flex broker based on the provided cursor, which represents the last country code retrieved. It returns a LocationPage containing the list of locations and a cursor for the next page, or an error if the retrieval fails.
func (f *Flex) GetLocationsPage(cursor string) (LocationPage, error) {
	countries, err := f.getCountries()
	if err != nil {
		return LocationPage{}, fmt.Errorf("flex get countries: %w", err)
	}

	if len(countries) == 0 {
		return LocationPage{}, nil
	}

	currentCursorIndex := 0
	if cursor != "" {
		currentCursorIndex = -1
		for i, c := range countries {
			if c.code == cursor {
				currentCursorIndex = i
				break
			}
		}
	}

	if currentCursorIndex == -1 {
		return LocationPage{}, fmt.Errorf("cursor not found in countries list: %s", cursor)
	}

	if currentCursorIndex >= len(countries) {
		return LocationPage{}, nil
	}

	var nextPage string
	if currentCursorIndex+1 < len(countries) {
		nextPage = countries[currentCursorIndex+1].code
	}

	c := countries[currentCursorIndex]
	locs, err := f.getLocationsFullDetails(c.name, c.code)
	if err != nil {
		// we want to allow continuing to the next country even if one country fails, so we return an empty page with the error
		return LocationPage{
			Locations: nil,
			NextPage:  nextPage,
		}, err
	}

	return LocationPage{
		Locations: locs,
		NextPage:  nextPage,
	}, nil
}

// flexCountry represents a country in the Flex broker, including its name and code.
type flexCountry struct {
	name string
	code string
}

// getCountriesResponse is the xml response for GetCountries.
type getCountriesResponse struct {
	CountrySet struct {
		Countries []struct {
			Name string `xml:"CountryName"`
			Code string `xml:"CountryCode"`
		} `xml:"Country"`
	} `xml:"CountrySet"`
}

// getCountries retrieves the list of countries from the Flex broker, using a cache to avoid unnecessary API calls. It returns a slice of flexCountry structs containing the country name and code, or an error if the retrieval fails.
func (f *Flex) getCountries() ([]flexCountry, error) {
	if len(f.countriesCache) > 0 {
		return f.countriesCache, nil
	}

	form := url.Values{}
	form.Set("Language", "")
	form.Set("AdditionalParameters", "")

	body, err := f.postForm("GetCountries", form)
	if err != nil {
		return nil, err
	}

	var resp getCountriesResponse

	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("flex GetCountries unmarshal response: %w", err)
	}

	out := make([]flexCountry, 0, len(resp.CountrySet.Countries))
	for _, c := range resp.CountrySet.Countries {
		out = append(out, flexCountry{
			name: c.Name,
			code: c.Code,
		})
	}

	f.countriesCache = out
	return out, nil
}

// getFullLocationDetailsResponse is the xml response for GetFullLocationDetails.
type getFullLocationDetailsResponse struct {
	LocationSet []struct {
		Full []struct {
			ID   string `xml:"LocationID"`
			Name string `xml:"LocationName"`
			IATA string `xml:"LocationIATA"`
		} `xml:"FullLocationDetails"`
	} `xml:"LocationSet"`
}

// getLocationsFullDetails retrieves the full details of locations for a given country from the Flex broker.
func (f *Flex) getLocationsFullDetails(country string, countryCode string) ([]Location, error) {
	form := url.Values{}
	form.Set("Language", "")
	form.Set("AdditionalParameters", "")
	form.Set("LocationID", "")
	form.Set("Supplier", "")
	form.Set("IATA", "")
	form.Set("State", "")
	form.Set("Area", "")
	form.Set("Country", countryCode)

	isRequiredLongTimeout := countryCode == "US" || countryCode == "CA"
	if isRequiredLongTimeout {
		f.httpClient.Timeout = 10 * time.Minute
	}

	body, err := f.postForm("GetFullLocationDetails", form)

	if isRequiredLongTimeout {
		f.httpClient.Timeout = defaultTimeout
	}

	if err != nil {
		return nil, fmt.Errorf("GetFullLocationDetails country=%s: %w", countryCode, err)
	}

	var resp getFullLocationDetailsResponse

	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("flex GetFullLocationDetails unmarshal response country=%s: %w", countryCode, err)
	}

	capHint := 0
	if len(resp.LocationSet) > 0 {
		capHint = len(resp.LocationSet[0].Full)
	}

	out := make([]Location, 0, capHint)
	for _, loc := range resp.LocationSet {
		for _, full := range loc.Full {
			out = append(out, Location{
				ID:          full.ID,
				Name:        full.Name,
				Iata:        full.IATA,
				Country:     country,
				CountryCode: countryCode,
			})
		}
	}

	return out, nil
}
