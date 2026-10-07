package broker

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
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
		// The prepared sheet holds 88 stations; 4 delivery-service rows carry no country and are
		// skipped.
		if got, want := len(page.Locations), 84; got != want {
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
		got, ok := byCode["022"]
		if !ok {
			t.Fatal("022 missing")
		}
		want := Location{
			ID:          "022",
			Name:        "Athens International Airport",
			City:        "Spata1",
			CountryCode: "GR",
			Country:     "Greece",
			Iata:        "ATH",
		}
		if got != want {
			t.Errorf("022 = %+v, want %+v", got, want)
		}
	})

	t.Run("takes the IATA from its column, whatever the code or station type", func(t *testing.T) {
		withIata := 0
		for _, l := range page.Locations {
			if l.Iata != "" {
				withIata++
			}
		}
		if got, want := withIata, 26; got != want {
			t.Errorf("stations with IATA = %d, want %d", got, want)
		}

		for code, want := range map[string]string{
			"022":   "ATH", // numeric live code
			"005":   "JTR", // typed Office
			"CFU1":  "CFU", // Airport Shuttle
			"AOKAP": "AOK", // Meet and Greet
			"RHOAP": "RHO",
			"010":   "", // Mykonos Main Station, 400m from the airport
			"001":   "", // Athens Downtown
		} {
			if got := byCode[code].Iata; got != want {
				t.Errorf("%s Iata = %q, want %q", code, got, want)
			}
		}
	})
}

func TestAvanceGetLocationsPageErrors(t *testing.T) {
	t.Run("reader not initialised", func(t *testing.T) {
		if _, err := NewAvance(AvanceConfig{}).GetLocationsPage(""); !errors.Is(err, ErrAvanceReaderNotInitialized) {
			t.Errorf("err = %v, want ErrAvanceReaderNotInitialized", err)
		}
	})

	t.Run("not a workbook", func(t *testing.T) {
		_, err := NewAvanceWithReader(strings.NewReader("not an xlsx")).GetLocationsPage("")
		if err == nil {
			t.Fatal("want an error")
		}
	})

	t.Run("sheet without the IATA column", func(t *testing.T) {
		// Avance's sheet as sent, before the Apps Script adds the column.
		f := excelize.NewFile()
		defer f.Close()
		if err := f.SetSheetName("Sheet1", avanceStationsSheet); err != nil {
			t.Fatal(err)
		}
		for i, v := range []string{"Code", "Description", "Charge", "Station type", "Street", "City", "Post code", "Country"} {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellStr(avanceStationsSheet, cell, v)
		}
		for i, v := range []string{"ATHAP", "Athens International Airport", "", "Airport", "Desk", "Spata", "19019", "GR"} {
			cell, _ := excelize.CoordinatesToCellName(i+1, 2)
			f.SetCellStr(avanceStationsSheet, cell, v)
		}
		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			t.Fatal(err)
		}

		_, err := NewAvanceWithReader(&buf).GetLocationsPage("")
		if !errors.Is(err, ErrAvanceInvalidWorkbook) {
			t.Errorf("err = %v, want ErrAvanceInvalidWorkbook", err)
		}
	})

	t.Run("empty reader", func(t *testing.T) {
		_, err := NewAvanceWithReader(bytes.NewReader(nil)).GetLocationsPage("")
		if err == nil {
			t.Fatal("want an error")
		}
	})
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
