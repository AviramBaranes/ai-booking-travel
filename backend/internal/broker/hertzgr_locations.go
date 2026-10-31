package broker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"encore.dev/rlog"
)

// hertzGRLocationsURL is Autohellas' public station feed. It lists every station of every country
// and brand, and Autohellas keeps it current, so the import reads it live rather than from a file.
const hertzGRLocationsURL = "https://prdota01.hertz.gr/OTAServices/Autohellaslocations.xml"

// HertzGRBrands are the brands Autohellas sells, as the OTA vendor codes its feed groups stations
// by: Hertz, Thrifty, Dollar and Firefly.
var HertzGRBrands = []string{"ZE", "ZT", "ZR", "FF"}

// hertzGRCountries are the countries whose stations we import. Our contract covers Greece only,
// while the feed lists all of Autohellas' countries. Montenegro is "MN" there, not ISO "ME".
var hertzGRCountries = map[string]struct{}{
	"GR": {},
}

// hertzGRAirportIATA gives the IATA code of each airport station, so that it merges onto the
// canonical airport location other brokers already created. The feed has no IATA codes, and its
// names do not give them reliably (Aghialos A/P is Volos, VOL). A brand working from its own desk
// at an airport has its own station, which maps to the same code: Mykonos is GR001 for Hertz and
// GR221 for Thrifty and Dollar.
var hertzGRAirportIATA = map[string]string{
	"GR001": "JMK", // Mykonos
	"GR016": "SMI", // Samos
	"GR017": "VOL", // Aghialos (Volos)
	"GR030": "ATH", // Athens
	"GR034": "HER", // Heraklion
	"GR039": "RHO", // Rhodes
	"GR041": "CFU", // Corfu
	"GR042": "SKG", // Thessaloniki
	"GR043": "CHQ", // Chania
	"GR044": "KGS", // Kos
	"GR064": "KLX", // Kalamata
	"GR067": "MJT", // Mytilini
	"GR068": "JTR", // Santorini
	"GR074": "JKH", // Chios
	"GR076": "PVK", // Preveza
	"GR079": "IOA", // Ioannina
	"GR084": "EFL", // Kefalonia
	"GR086": "AXD", // Alexandroupoli
	"GR089": "ZTH", // Zakynthos
	"GR097": "KVA", // Kavala
	"GR117": "GPA", // Araxos
	"GR120": "JSI", // Skiathos
	"GR121": "AOK", // Karpathos
	"GR125": "PAS", // Paros
	"GR126": "JSH", // Sitia
	"GR137": "MLO", // Milos
	"GR253": "JNX", // Naxos
	"GR221": "JMK", // Mykonos
	"GR222": "RHO", // Rhodes
	"GR223": "CFU", // Corfu
	"GR224": "SKG", // Thessaloniki
	"GR225": "CHQ", // Chania
	"GR226": "EFL", // Kefalonia
	"GR230": "KGS", // Kos
	"GR423": "ATH", // Athens
	"GR424": "JTR", // Santorini
	"GR425": "PAS", // Paros
	"GR426": "PVK", // Preveza
}

// hertzGRNameAbbreviations spells out the abbreviations station names use ("Mykonos A/P").
var hertzGRNameAbbreviations = strings.NewReplacer("A/P", "Airport", "D/T", "Downtown")

var ErrHertzGRInvalidFeed = errors.New("invalid Hertz GR locations feed")

// hertzGRLocationsFeed is the station feed: countries, each split into brands, each listing the
// stations that sell it. A station shared by several brands appears under each of them.
type hertzGRLocationsFeed struct {
	Countries []struct {
		Name   string `xml:"Name,attr"`
		Code   string `xml:"Code,attr"`
		Brands []struct {
			Code     string           `xml:"Code,attr"`
			Stations []hertzGRStation `xml:"Location"`
		} `xml:"Brand"`
	} `xml:"Country"`
}

// hertzGRStation holds the station fields we import. The feed also carries coordinates (with
// Longtitude and Latitude holding each other's values), phone and seasonal opening hours, which
// broker.Location has no place for.
type hertzGRStation struct {
	Code        string `xml:"Code,attr"`
	Name        string `xml:"Name,attr"`
	Airport     string `xml:"Airport,attr"`
	City        string `xml:"Address>CityName"`
	AddressLine string `xml:"Address>AddressLine"`
}

// GetLocationsPage retrieves the stations of the countries we import from the Autohellas feed. The
// cursor is unused since all stations are returned in a single page.
func (h *HertzGR) GetLocationsPage(cursor string) (LocationPage, error) {
	r, err := h.openLocationsFeed()
	if err != nil {
		return LocationPage{}, err
	}
	defer r.Close()

	var feed hertzGRLocationsFeed
	if err := xml.NewDecoder(r).Decode(&feed); err != nil {
		return LocationPage{}, fmt.Errorf("hertz gr locations: decode feed: %w", err)
	}

	locations := hertzGRLocations(feed)
	if len(locations) == 0 {
		return LocationPage{}, ErrHertzGRInvalidFeed
	}

	return LocationPage{Locations: locations}, nil
}

func (h *HertzGR) openLocationsFeed() (io.ReadCloser, error) {
	if h.r != nil {
		return io.NopCloser(h.r), nil
	}

	res, err := h.httpClient.Get(hertzGRLocationsURL)
	if err != nil {
		return nil, fmt.Errorf("hertz gr locations: fetch feed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("hertz gr locations: fetch feed: status %d", res.StatusCode)
	}

	return res.Body, nil
}

// hertzGRLocations flattens the feed into one location per station code, carrying every brand the
// station appears under, in feed order.
func hertzGRLocations(feed hertzGRLocationsFeed) []Location {
	locations := make([]Location, 0)
	indexOf := make(map[string]int)

	for _, country := range feed.Countries {
		countryCode := strings.ToUpper(strings.TrimSpace(country.Code))
		if _, ok := hertzGRCountries[countryCode]; !ok {
			continue
		}

		for _, brand := range country.Brands {
			brandCode := strings.TrimSpace(brand.Code)
			if brandCode == "" {
				rlog.Warn("skipping Hertz GR brand with no code", "country", countryCode)
				continue
			}

			for _, s := range brand.Stations {
				code := strings.TrimSpace(s.Code)
				if code == "" {
					continue
				}

				if i, ok := indexOf[code]; ok {
					if !slices.Contains(locations[i].Brands, brandCode) {
						locations[i].Brands = append(locations[i].Brands, brandCode)
					}
					continue
				}

				name := hertzGRStationName(s)
				if name == "" {
					rlog.Warn("skipping Hertz GR station with no name or address", "code", code)
					continue
				}

				indexOf[code] = len(locations)
				locations = append(locations, Location{
					ID:          code,
					Name:        name,
					Country:     strings.TrimSpace(country.Name),
					CountryCode: countryCode,
					City:        strings.TrimSpace(s.City),
					Iata:        hertzGRStationIATA(code, s),
					Brands:      []string{brandCode},
				})
			}
		}
	}

	for i := range locations {
		slices.Sort(locations[i].Brands)
	}

	return locations
}

// hertzGRStationName spells out the feed's abbreviations. A station with no name (the Heraklion
// port desk) is named by its address line instead.
func hertzGRStationName(s hertzGRStation) string {
	name := s.Name
	if strings.TrimSpace(name) == "" {
		name = s.AddressLine
	}

	return strings.Join(strings.Fields(hertzGRNameAbbreviations.Replace(name)), " ")
}

// hertzGRStationIATA looks the station up in hertzGRAirportIATA. An airport station missing from it
// still imports, as a location of its own, and the warning says which code to add.
func hertzGRStationIATA(code string, s hertzGRStation) string {
	if iata, ok := hertzGRAirportIATA[code]; ok {
		return iata
	}

	if strings.EqualFold(strings.TrimSpace(s.Airport), "true") {
		rlog.Warn("Hertz GR airport station has no IATA in hertzGRAirportIATA, importing it without one",
			"code", code, "name", s.Name)
	}

	return ""
}
