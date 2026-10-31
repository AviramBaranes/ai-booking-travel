package broker

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func hertzGRTestPage(t *testing.T) LocationPage {
	t.Helper()

	f, err := os.Open("testdata/hertzgr_locations.xml")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	page, err := NewHertzGRWithReader(f).GetLocationsPage("")
	if err != nil {
		t.Fatalf("GetLocationsPage: %v", err)
	}

	return page
}

// The fixture is a verbatim excerpt of the live feed: six Greek stations, listed under every brand
// that serves them, and one Cypriot station.
func TestHertzGRGetLocationsPage(t *testing.T) {
	page := hertzGRTestPage(t)

	byCode := make(map[string]Location, len(page.Locations))
	for _, l := range page.Locations {
		byCode[l.ID] = l
	}

	t.Run("returns each Greek station once, in one page", func(t *testing.T) {
		if page.NextPage != "" {
			t.Errorf("NextPage = %q, want empty", page.NextPage)
		}
		if got, want := len(page.Locations), 6; got != want {
			t.Errorf("len(Locations) = %d, want %d", got, want)
		}
	})

	t.Run("leaves out countries we have no contract for", func(t *testing.T) {
		if _, ok := byCode["CY040"]; ok {
			t.Error("CY040 was imported, want only Greek stations")
		}
	})

	t.Run("collects every brand a station is listed under", func(t *testing.T) {
		cases := map[string][]string{
			"GR001": {"ZE"},
			"GR221": {"ZR", "ZT"},
			"GR030": {"FF", "ZE", "ZT"},
			"GR021": {"ZE", "ZT"},
			"GR136": {"ZE"},
		}
		for code, want := range cases {
			if got := byCode[code].Brands; !reflect.DeepEqual(got, want) {
				t.Errorf("%s brands = %v, want %v", code, got, want)
			}
		}
	})

	t.Run("maps an airport station onto its canonical location", func(t *testing.T) {
		want := Location{
			ID:          "GR030",
			Name:        "Athens International Airport",
			Country:     "Greece",
			City:        "Athens",
			CountryCode: "GR",
			Iata:        "ATH",
			Brands:      []string{"FF", "ZE", "ZT"},
		}
		if got := byCode["GR030"]; !reflect.DeepEqual(got, want) {
			t.Errorf("GR030 = %+v, want %+v", got, want)
		}
	})

	t.Run("gives a brand's own airport desk the airport's IATA", func(t *testing.T) {
		for _, code := range []string{"GR001", "GR221"} {
			if got := byCode[code].Iata; got != "JMK" {
				t.Errorf("%s Iata = %q, want JMK", code, got)
			}
		}
	})

	t.Run("leaves stations away from the airport without an IATA", func(t *testing.T) {
		for _, code := range []string{"GR002", "GR021", "GR136"} {
			if got := byCode[code].Iata; got != "" {
				t.Errorf("%s Iata = %q, want empty", code, got)
			}
		}
	})

	t.Run("spells out the name abbreviations", func(t *testing.T) {
		cases := map[string]string{
			"GR001": "Mykonos Airport",
			"GR002": "Kefalonia Downtown",
			"GR021": "Hilton",
		}
		for code, want := range cases {
			if got := byCode[code].Name; got != want {
				t.Errorf("%s Name = %q, want %q", code, got, want)
			}
		}
	})

	t.Run("names a station without a name by its address", func(t *testing.T) {
		got := byCode["GR136"]
		if got.Name != "Heraklion Port terminal" || got.City != "Crete, Heraklion" {
			t.Errorf("GR136 = name %q, city %q, want the address line and the feed's city", got.Name, got.City)
		}
	})
}

func TestHertzGRGetLocationsPageEdgeCases(t *testing.T) {
	parse := func(t *testing.T, feed string) (LocationPage, error) {
		t.Helper()
		return NewHertzGRWithReader(strings.NewReader(feed)).GetLocationsPage("")
	}

	t.Run("imports an airport station missing from the IATA map without one", func(t *testing.T) {
		page, err := parse(t, `<AutoHellasLocations><Country Name="Greece" Code="GR"><Brand Code="ZE">
			<Location Code="GR999" Name="New Island A/P" Airport="True"><Address><CityName>New Island</CityName></Address></Location>
		</Brand></Country></AutoHellasLocations>`)
		if err != nil {
			t.Fatalf("GetLocationsPage: %v", err)
		}
		if len(page.Locations) != 1 || page.Locations[0].Iata != "" || page.Locations[0].Name != "New Island Airport" {
			t.Errorf("Locations = %+v, want New Island Airport without an IATA", page.Locations)
		}
	})

	t.Run("skips a station with neither a name nor an address, and a brand with no code", func(t *testing.T) {
		page, err := parse(t, `<AutoHellasLocations><Country Name="Greece" Code="GR">
			<Brand Code="ZE">
				<Location Code="GR998" Name="" Airport="False"><Address><CityName>Athens</CityName></Address></Location>
				<Location Code="GR997" Name="Kept" Airport="False"><Address><CityName>Athens</CityName></Address></Location>
			</Brand>
			<Brand Code=""><Location Code="GR997" Name="Kept" Airport="False"></Location></Brand>
		</Country></AutoHellasLocations>`)
		if err != nil {
			t.Fatalf("GetLocationsPage: %v", err)
		}
		if len(page.Locations) != 1 || page.Locations[0].ID != "GR997" || !reflect.DeepEqual(page.Locations[0].Brands, []string{"ZE"}) {
			t.Errorf("Locations = %+v, want only GR997 with brand ZE", page.Locations)
		}
	})

	t.Run("rejects a feed with no station in our countries", func(t *testing.T) {
		_, err := parse(t, `<AutoHellasLocations><Country Name="Cyprus" Code="CY"><Brand Code="ZE">
			<Location Code="CY040" Name="Larnaca A/P" Airport="True"></Location>
		</Brand></Country></AutoHellasLocations>`)
		if !errors.Is(err, ErrHertzGRInvalidFeed) {
			t.Errorf("err = %v, want ErrHertzGRInvalidFeed", err)
		}
	})

	t.Run("reports a feed that is not XML", func(t *testing.T) {
		if _, err := parse(t, "<html>maintenance"); err == nil {
			t.Error("err = nil, want a decode error")
		}
	})
}
