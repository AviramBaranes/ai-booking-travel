package supplier

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var secrets struct {
	AgentCode    string
	FlexPassword string
	FlexBaseURL  string
}

type Flex struct {
	httpClient     *http.Client
	countriesCache []flexCountry
}

const defaultTimeout = 5 * time.Minute

func NewFlex() Flex {
	return Flex{
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
}

func (f Flex) Supplier() SupplierName {
	return SupplierFlex
}

func (f *Flex) GetLocationsPage(cursor string) (LocationPage, error) {
	countries, err := f.getCountries()
	if err != nil {
		return LocationPage{}, err
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
		return LocationPage{}, fmt.Errorf("cursor not found in countries list")
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

type flexCountry struct {
	name string
	code string
}

func (f *Flex) getCountries() ([]flexCountry, error) {
	if f.countriesCache != nil && len(f.countriesCache) > 0 {
		return f.countriesCache, nil
	}

	form := url.Values{}
	form.Set("Language", "")
	form.Set("AdditionalParameters", "")

	body, err := f.postForm("GetCountries", form)
	if err != nil {
		return nil, err
	}

	var resp struct {
		CountrySet struct {
			Countries []struct {
				Name string `xml:"CountryName"`
				Code string `xml:"CountryCode"`
			} `xml:"Country"`
		} `xml:"CountrySet"`
	}

	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, err
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

	var resp struct {
		LocationSet []struct {
			Full []struct {
				ID   string `xml:"LocationID"`
				Name string `xml:"LocationName"`
				IATA string `xml:"LocationIATA"`
			} `xml:"FullLocationDetails"`
		} `xml:"LocationSet"`
	}

	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, err
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

func (f Flex) postForm(ep string, form url.Values) ([]byte, error) {
	form.Set("AgentCode", secrets.AgentCode)
	form.Set("Password", secrets.FlexPassword)

	url := secrets.FlexBaseURL + "/" + ep
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json,application/xml,application/xml;q=0.9,*/*;q=0.8")

	res, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return b, nil
}
