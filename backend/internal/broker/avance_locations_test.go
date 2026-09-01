package broker

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func avanceTestPage(t *testing.T) LocationPage {
	t.Helper()

	f, err := os.Open("testdata/avance_stations.xlsx")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	page, err := NewAvanceWithReader(f).GetLocationsPage("")
	if err != nil {
		t.Fatalf("GetLocationsPage: %v", err)
	}

	return page
}

func TestAvanceGetLocationsPage(t *testing.T) {
	page := avanceTestPage(t)

	byCode := make(map[string]Location, len(page.Locations))
	for _, l := range page.Locations {
		byCode[l.ID] = l
	}

	t.Run("returns every station in one page", func(t *testing.T) {
		if page.NextPage != "" {
			t.Errorf("NextPage = %q, want empty", page.NextPage)
		}
		// The sheet holds 90 stations; 4 delivery-service rows carry no country and are skipped.
		if got, want := len(page.Locations), 86; got != want {
			t.Errorf("len(Locations) = %d, want %d", got, want)
		}
	})

	t.Run("skips stations missing a country", func(t *testing.T) {
		for _, code := range []string{"ITFP", "ITAPT", "JMKPR", "SKGPR"} {
			if _, ok := byCode[code]; ok {
				t.Errorf("station %s was imported, want skipped for missing country", code)
			}
		}
	})

	t.Run("maps station fields onto the canonical location", func(t *testing.T) {
		got, ok := byCode["ATHAP"]
		if !ok {
			t.Fatal("ATHAP missing")
		}
		want := Location{
			ID:          "ATHAP",
			Name:        "Athens International Airport",
			City:        "Spata1",
			CountryCode: "GR",
			Country:     "Greece",
			Iata:        "ATH",
		}
		if got != want {
			t.Errorf("ATHAP = %+v, want %+v", got, want)
		}
	})

	t.Run("gives an IATA to airport desks only", func(t *testing.T) {
		withIata := make([]string, 0)
		for _, l := range page.Locations {
			if l.Iata != "" {
				withIata = append(withIata, l.ID)
			}
		}
		if got, want := len(withIata), 15; got != want {
			t.Errorf("stations with IATA = %d (%v), want %d", got, withIata, want)
		}

		for code, want := range map[string]string{
			"ATHAP": "ATH", // Airport
			"RHOAP": "RHO", // Airport 1
			"JMKAP": "JMK",
			"ZTHAP": "ZTH",
		} {
			if got := byCode[code].Iata; got != want {
				t.Errorf("%s Iata = %q, want %q", code, got, want)
			}
		}

		// Shuttle, meet-and-greet and office pickups at an airport are a different arrangement.
		for _, code := range []string{"CFUAP", "PASAP", "AOKAP", "KLXAP", "JTRAP", "SKGAP", "ATHDT"} {
			if got := byCode[code].Iata; got != "" {
				t.Errorf("%s Iata = %q, want empty (%q)", code, got, byCode[code].Name)
			}
		}
	})
}

func TestAvanceGetLocationsPageErrors(t *testing.T) {
	t.Run("reader not initialised", func(t *testing.T) {
		if _, err := NewAvance().GetLocationsPage(""); !errors.Is(err, ErrAvanceReaderNotInitialized) {
			t.Errorf("err = %v, want ErrAvanceReaderNotInitialized", err)
		}
	})

	t.Run("not a workbook", func(t *testing.T) {
		_, err := NewAvanceWithReader(strings.NewReader("not an xlsx")).GetLocationsPage("")
		if err == nil {
			t.Fatal("want an error")
		}
	})

	t.Run("empty reader", func(t *testing.T) {
		_, err := NewAvanceWithReader(bytes.NewReader(nil)).GetLocationsPage("")
		if err == nil {
			t.Fatal("want an error")
		}
	})
}

func TestAvanceStationIata(t *testing.T) {
	tests := []struct {
		code, stationType, want string
	}{
		{"ATHAP", "Airport", "ATH"},
		{"RHOAP", "Airport 1", "RHO"},
		{"RHOAP", " Airport 1 ", "RHO"},
		{"athap", "Airport", "ATH"},
		{"ATHDT", "Office", ""},
		{"CFUAP", "Airport Shuttle", ""},
		{"AOKAP", "Meet and Greet", ""},
		{"RHOPR", "Port", ""},
		{"ALD", "Hotel", ""},
		{"AB", "Airport", ""},
	}

	for _, tt := range tests {
		if got := avanceStationIata(tt.code, tt.stationType); got != tt.want {
			t.Errorf("avanceStationIata(%q, %q) = %q, want %q", tt.code, tt.stationType, got, tt.want)
		}
	}
}

func TestParseAvanceCredentials(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := parseAvanceCredentials("10268, TBAVA ,ABTAV")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := avanceCredentials{accountNo: "10268", linkCode: "TBAVA", agentCode: "ABTAV"}
		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	for _, s := range []string{"", "10268", "10268,TBAVA", "10268,TBAVA,ABTAV,extra", "10268,,ABTAV"} {
		if _, err := parseAvanceCredentials(s); err == nil {
			t.Errorf("parseAvanceCredentials(%q) = nil error, want an error", s)
		}
	}
}
