package broker

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

var hertzGRTestRates = []HertzGRRate{
	{Code: "ITRH296", PlanName: "Standard", ErpDayCharge: 3},
	{Code: "ITRT296", PlanName: "Standard", ErpDayCharge: 3},
	{Code: "ITRD296", PlanName: "Standard", ErpDayCharge: 3},
	{Code: "ITRF296", PlanName: "Standard", ErpDayCharge: 3},
}

// Mykonos Airport holds Hertz's desk and a desk shared by Thrifty and Dollar; Athens Airport a desk
// shared by Firefly, Hertz and Thrifty, and Dollar's own. Firefly has no Mykonos desk.
var (
	hertzGRTestMykonos = []Station{{Code: "GR001", Brands: []string{"ZE"}}, {Code: "GR221", Brands: []string{"ZR", "ZT"}}}
	hertzGRTestAthens  = []Station{{Code: "GR030", Brands: []string{"FF", "ZE", "ZT"}}, {Code: "GR423", Brands: []string{"ZR"}}}
)

func hertzGRTestParams() SearchAvailabilityParams {
	return SearchAvailabilityParams{
		CountryCode:     "GR",
		PickupLocation:  "GR001",
		DropoffLocation: "GR001",
		PickupStations:  hertzGRTestMykonos,
		DropoffStations: hertzGRTestMykonos,
		PickupDate:      "2026-12-10",
		PickupTime:      "10:00",
		DropoffDate:     "2026-12-17",
		DropoffTime:     "10:00",
		DriverAge:       30,
	}
}

func hertzGRTestResponse(t *testing.T, fixture string) *hertzGRVehAvailRateRS {
	t.Helper()
	var rs hertzGRVehAvailRateRS
	if err := decodeHertzGRResponse([]byte(fixture), "VehicleAvailRateResult", &rs); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	return &rs
}

func hertzGRTestStations(t *testing.T) map[string]hertzGRStation {
	t.Helper()
	f, err := os.Open("testdata/hertzgr_locations.xml")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()

	stations, err := NewHertzGRWithReader(f).fetchStations()
	if err != nil {
		t.Fatalf("fetchStations: %v", err)
	}
	return stations
}

func TestHertzGRQuotes(t *testing.T) {
	t.Run("pairs each brand with its own desks", func(t *testing.T) {
		p := hertzGRTestParams()
		p.DropoffStations = hertzGRTestAthens

		got := hertzGRQuotes(hertzGRTestRates, p)
		want := []hertzGRQuote{
			{rate: hertzGRTestRates[0], planID: 1, brand: "ZE", pickup: "GR001", dropoff: "GR030"},
			{rate: hertzGRTestRates[1], planID: 2, brand: "ZT", pickup: "GR221", dropoff: "GR030"},
			{rate: hertzGRTestRates[2], planID: 3, brand: "ZR", pickup: "GR221", dropoff: "GR423"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("quotes = %+v, want %+v (Firefly has no Mykonos desk)", got, want)
		}
	})

	t.Run("skips a rate code whose brand letter is unknown", func(t *testing.T) {
		got := hertzGRQuotes([]HertzGRRate{{Code: "ITRX296"}, {Code: "IT"}}, hertzGRTestParams())
		if len(got) != 0 {
			t.Errorf("quotes = %+v, want none", got)
		}
	})
}

func TestBuildHertzGRAvailability(t *testing.T) {
	p := hertzGRTestParams()
	stations := hertzGRTestStations(t)
	quotes := hertzGRQuotes(hertzGRTestRates[:2], p) // Hertz at GR001, Thrifty at GR221

	sameStation := hertzGRTestResponse(t, hertzGRSameStationFixture)
	resp := buildHertzGRAvailability(p, 7, quotes, []*hertzGRVehAvailRateRS{sameStation, sameStation}, stations, nil)

	t.Run("sells only available cars, one vehicle per brand and group", func(t *testing.T) {
		var got []string
		for _, v := range resp.AvailableVehicles {
			got = append(got, v.CarDetails.SupplierName+":"+v.CarDetails.CarGroup)
		}
		want := []string{"Hertz:A", "Hertz:K", "Thrifty:A", "Thrifty:K"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("vehicles = %v, want %v", got, want)
		}
	})

	t.Run("maps the car details", func(t *testing.T) {
		want := CarDetails{
			Model:        "CHEVROLET MATIZ",
			CarGroup:     "A",
			ImageURL:     "https://images.autohellas.gr/CarImages/i10.gif",
			SupplierName: "Hertz",
			CarType:      "car",
			Acriss:       "MDMR",
			FullAcriss:   "MDMR",
			HasAC:        true,
			IsAutoGear:   false,
			IsElectric:   false,
			Seats:        4,
			Bags:         2,
			Doors:        4,
		}
		if got := resp.AvailableVehicles[0].CarDetails; got != want {
			t.Errorf("car details = %+v, want %+v", got, want)
		}
		if !resp.AvailableVehicles[1].CarDetails.IsAutoGear {
			t.Error("the automatic car is not marked automatic")
		}
		if got := resp.AvailableVehicles[0].LocationDetails.LocationType; got != "Airport" {
			t.Errorf("location type = %q, want Airport", got)
		}
	})

	t.Run("prices a plan per rate code at its brand's desk", func(t *testing.T) {
		got := resp.AvailableVehicles[2].Plans // Thrifty, group A
		want := []Plan{{
			PlanID:                 2,
			PlanName:               "Standard",
			Price:                  123,
			ChargedErpPriceWithVat: 21,
			RateQualifier:          "ITRT296:A",
			SupplierName:           "Thrifty",
			SupplierCode:           "ZT",
			PickupLocationCode:     "GR221",
			DropoffLocationCode:    "GR221",
		}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("plans = %+v, want %+v", got, want)
		}
		if resp.AvailableVehicles[2].PriceDetails.Currency != "EUR" {
			t.Errorf("currency = %q, want EUR", resp.AvailableVehicles[2].PriceDetails.Currency)
		}
	})

	t.Run("identifies every plan uniquely", func(t *testing.T) {
		seen := make(map[string]bool)
		for _, v := range resp.AvailableVehicles {
			for _, pl := range v.Plans {
				key := pl.RateQualifier + "|" + pl.SupplierCode + "|" + pl.PlanName
				if seen[key] {
					t.Errorf("duplicate plan identity %s", key)
				}
				seen[key] = true
			}
		}
	})

	t.Run("describes each brand once, named like its plans", func(t *testing.T) {
		if len(resp.SuppliersInfo) != 2 {
			t.Fatalf("suppliers = %d, want 2", len(resp.SuppliersInfo))
		}
		hertz := resp.SuppliersInfo[0]
		if hertz.Name != "Hertz" || resp.SuppliersInfo[1].Name != "Thrifty" {
			t.Errorf("supplier names = %q, %q", hertz.Name, resp.SuppliersInfo[1].Name)
		}
		wantIncs := []Inclusions{{ProductName: "Standard", ProductInclusions: []string{"Airport Service", "Unlimited mileage"}}}
		if !reflect.DeepEqual(hertz.Inclusions, wantIncs) {
			t.Errorf("inclusions = %+v, want %+v", hertz.Inclusions, wantIncs)
		}
		if len(hertz.AddOns) != 0 {
			t.Errorf("add-ons = %+v, want none", hertz.AddOns)
		}
		if hertz.PickupDetails.Address != "Airport, Mykonos" || hertz.PickupDetails.PhoneNumber != "(+30) 22890 27511" {
			t.Errorf("pickup details = %+v, want Mykonos Airport's desk", hertz.PickupDetails)
		}
	})

	t.Run("a failed quote leaves the other brands", func(t *testing.T) {
		resp := buildHertzGRAvailability(p, 7, quotes, []*hertzGRVehAvailRateRS{nil, sameStation}, stations, nil)
		for _, v := range resp.AvailableVehicles {
			if v.CarDetails.SupplierName != "Thrifty" {
				t.Errorf("vehicle of %s returned, want only Thrifty", v.CarDetails.SupplierName)
			}
		}
		if len(resp.AvailableVehicles) != 2 {
			t.Errorf("vehicles = %d, want 2", len(resp.AvailableVehicles))
		}
	})
}

func TestBuildHertzGRAvailabilityOneWay(t *testing.T) {
	p := hertzGRTestParams()
	p.DropoffStations = hertzGRTestAthens
	quotes := hertzGRQuotes(hertzGRTestRates[:1], p)
	rs := hertzGRTestResponse(t, strings.ReplaceAll(hertzGROneWayFixture, `Status="OnRequest"`, `Status="Available"`))

	resp := buildHertzGRAvailability(p, 2, quotes, []*hertzGRVehAvailRateRS{rs}, nil, nil)
	if len(resp.AvailableVehicles) != 1 {
		t.Fatalf("vehicles = %d, want 1", len(resp.AvailableVehicles))
	}
	v := resp.AvailableVehicles[0]

	if got := v.Plans[0].Price; got != 84 {
		t.Errorf("price = %v, want the 234.06 total less the 150.06 paid on arrival", got)
	}
	if v.PriceDetails.Fees.DropCharge != 150 || v.PriceDetails.Fees.DropChargeCurrency != "EUR" {
		t.Errorf("fees = %+v, want a 150 EUR drop charge", v.PriceDetails.Fees)
	}
	if len(v.Plans[0].Info) != 0 {
		t.Errorf("info = %v, want none: the drop charge is all that is paid on arrival", v.Plans[0].Info)
	}
	if v.Plans[0].DropoffLocationCode != "GR030" {
		t.Errorf("dropoff = %q, want GR030", v.Plans[0].DropoffLocationCode)
	}
	if v.LocationDetails.LocationType != "City" {
		t.Errorf("location type = %q, want City for a response without a park location", v.LocationDetails.LocationType)
	}
}

func TestDecodeHertzGRResponse(t *testing.T) {
	t.Run("reports a SOAP fault", func(t *testing.T) {
		var rs hertzGRVehAvailRateRS
		err := decodeHertzGRResponse([]byte(hertzGRFaultFixture), "VehicleAvailRateResult", &rs)
		if err == nil || !strings.Contains(err.Error(), "Authentication failed") {
			t.Errorf("err = %v, want the fault", err)
		}
	})

	t.Run("reads a result sent as escaped text", func(t *testing.T) {
		start := strings.Index(hertzGRSameStationFixture, "<OTA_VehAvailRateRS")
		end := strings.Index(hertzGRSameStationFixture, "</VehicleAvailRateResult>")
		var escaped strings.Builder
		if err := xml.EscapeText(&escaped, []byte(hertzGRSameStationFixture[start:end])); err != nil {
			t.Fatal(err)
		}
		body := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><VehicleAvailRateResult xmlns="http://www.hertz.gr/xml">` +
			escaped.String() + `</VehicleAvailRateResult></s:Body></s:Envelope>`

		var rs hertzGRVehAvailRateRS
		if err := decodeHertzGRResponse([]byte(body), "VehicleAvailRateResult", &rs); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(rs.VendorAvails) != 1 || len(rs.VendorAvails[0].Avails) != 3 {
			t.Errorf("decoded %+v, want the three cars", rs.VendorAvails)
		}
	})

	t.Run("reports a response without the result", func(t *testing.T) {
		var rs hertzGRVehAvailRateRS
		if err := decodeHertzGRResponse([]byte(`<s:Envelope xmlns:s="x"><s:Body/></s:Envelope>`), "VehicleAvailRateResult", &rs); err == nil {
			t.Error("err = nil, want an error")
		}
	})
}

// hertzGRTestServer stands in for Autohellas, answering every call with respond and recording the
// requests. Calls come from concurrent quotes, hence the mutex.
func hertzGRTestServer(t *testing.T, respond func(body string) string) (*HertzGR, func() []*http.Request, func() []string) {
	t.Helper()

	savedCreds, savedErr := hertzGRCreds, hertzGRCredsErr
	hertzGRCreds = hertzGRCredentials{userName: "user", password: "secret", ccn: "CCN1", vc: "VC1"}
	hertzGRCredsErr = nil
	t.Cleanup(func() { hertzGRCreds, hertzGRCredsErr = savedCreds, savedErr })

	// Searches decorate results from the station feed, which the tests must not fetch.
	c := &hertzGRStationsCache
	c.mu.Lock()
	savedStations, savedExpires := c.stations, c.expires
	c.stations, c.expires = hertzGRTestStations(t), time.Now().Add(time.Hour)
	c.mu.Unlock()
	t.Cleanup(func() {
		c.mu.Lock()
		c.stations, c.expires = savedStations, savedExpires
		c.mu.Unlock()
	})

	fc := &hertzGRFleetCache
	fc.mu.Lock()
	savedImages, savedFleetExpires := fc.images, fc.expires
	fc.images, fc.expires = map[string]string{}, time.Now().Add(time.Hour)
	fc.mu.Unlock()
	t.Cleanup(func() {
		fc.mu.Lock()
		fc.images, fc.expires = savedImages, savedFleetExpires
		fc.mu.Unlock()
	})

	var (
		mu     sync.Mutex
		reqs   []*http.Request
		bodies []string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		reqs = append(reqs, r)
		bodies = append(bodies, string(b))
		mu.Unlock()
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		_, _ = io.WriteString(w, respond(string(b)))
	}))
	t.Cleanup(srv.Close)

	h := NewHertzGR(HertzGRConfig{BaseURL: srv.URL, Rates: hertzGRTestRates})
	return h,
		func() []*http.Request { mu.Lock(); defer mu.Unlock(); return reqs },
		func() []string { mu.Lock(); defer mu.Unlock(); return bodies }
}

func TestHertzGRSearchAvailability(t *testing.T) {
	t.Run("sends one SOAP call per brand with our credentials", func(t *testing.T) {
		h, reqs, bodies := hertzGRTestServer(t, func(string) string { return hertzGRSameStationFixture })

		resp, err := h.SearchAvailability(hertzGRTestParams())
		if err != nil {
			t.Fatalf("SearchAvailability: %v", err)
		}
		// Hertz, Thrifty and Dollar have Mykonos desks; Firefly is not quoted.
		if got := len(reqs()); got != 3 {
			t.Fatalf("requests = %d, want 3", got)
		}
		if len(resp.AvailableVehicles) != 6 {
			t.Errorf("vehicles = %d, want two cars for each of three brands", len(resp.AvailableVehicles))
		}

		r := reqs()[0]
		if r.Method != http.MethodPost || r.Header.Get("SOAPAction") != "VehicleAvailRate" || !strings.HasPrefix(r.Header.Get("Content-Type"), "text/xml") {
			t.Errorf("request = %s, SOAPAction %q, Content-Type %q", r.Method, r.Header.Get("SOAPAction"), r.Header.Get("Content-Type"))
		}

		var thrifty string
		for _, b := range bodies() {
			if strings.Contains(b, `RateQualifier="ITRT296"`) {
				thrifty = b
			}
		}
		for _, want := range []string{
			`<AuthenticationHeader xmlns="http://www.hertz.gr/xml"><UserName>user</UserName><Password>secret</Password></AuthenticationHeader>`,
			`<VehicleAvailRateRequest xmlns="http://www.hertz.gr/xml"><OTA_VehAvailRateRQ xmlns="http://www.opentravel.org/OTA/2003/05"`,
			`<RequestorID Type="4" ID="CCN1"><CompanyName Code="VC1"></CompanyName></RequestorID>`,
			`<VehAvailRQCore Status="Available">`,
			`PickUpDateTime="2026-12-10T10:00:00" ReturnDateTime="2026-12-17T10:00:00"`,
			`<PickUpLocation LocationCode="GR221">`,
			`<DriverType Age="30"></DriverType>`,
		} {
			if !strings.Contains(thrifty, want) {
				t.Errorf("Thrifty request lacks %s:\n%s", want, thrifty)
			}
		}
	})

	t.Run("a response with nothing to sell is empty, not an error", func(t *testing.T) {
		noCars := strings.NewReplacer("CODE", "301", "TEXT", "No cars available").Replace(hertzGRErrorFixture)
		h, _, _ := hertzGRTestServer(t, func(string) string { return noCars })

		resp, err := h.SearchAvailability(hertzGRTestParams())
		if err != nil || len(resp.AvailableVehicles) != 0 {
			t.Errorf("resp = %+v, err = %v, want no vehicles and no error", resp, err)
		}
	})

	t.Run("fails when every brand fails", func(t *testing.T) {
		badAgent := strings.NewReplacer("CODE", "168", "TEXT", "Agent not recognized").Replace(hertzGRErrorFixture)
		h, _, _ := hertzGRTestServer(t, func(string) string { return badAgent })

		if _, err := h.SearchAvailability(hertzGRTestParams()); err == nil || !strings.Contains(err.Error(), "Agent not recognized") {
			t.Errorf("err = %v, want the agent error", err)
		}
	})

	t.Run("one brand failing leaves the others", func(t *testing.T) {
		badAgent := strings.NewReplacer("CODE", "168", "TEXT", "Agent not recognized").Replace(hertzGRErrorFixture)
		h, _, _ := hertzGRTestServer(t, func(body string) string {
			if strings.Contains(body, "ITRD296") {
				return badAgent
			}
			return hertzGRSameStationFixture
		})

		resp, err := h.SearchAvailability(hertzGRTestParams())
		if err != nil || len(resp.AvailableVehicles) != 4 {
			t.Errorf("vehicles = %d, err = %v, want Hertz and Thrifty's four", len(resp.AvailableVehicles), err)
		}
	})

	t.Run("returns nothing without rate codes, without calling Autohellas", func(t *testing.T) {
		h, reqs, _ := hertzGRTestServer(t, func(string) string { return hertzGRSameStationFixture })
		h.rates = nil

		resp, err := h.SearchAvailability(hertzGRTestParams())
		if err != nil || len(resp.AvailableVehicles) != 0 || len(reqs()) != 0 {
			t.Errorf("vehicles = %d, requests = %d, err = %v", len(resp.AvailableVehicles), len(reqs()), err)
		}
	})
}

func TestParseHertzGRCredentials(t *testing.T) {
	got, err := parseHertzGRCredentials(" user , pass,CCN1,VC1 ")
	if err != nil || got != (hertzGRCredentials{userName: "user", password: "pass", ccn: "CCN1", vc: "VC1"}) {
		t.Errorf("parse = %+v, %v", got, err)
	}
	for _, bad := range []string{"", "user,pass,CCN1", "user,,CCN1,VC1"} {
		if _, err := parseHertzGRCredentials(bad); err == nil {
			t.Errorf("parse(%q) err = nil, want an error", bad)
		}
	}
}

func TestHertzGRStationOpeningHours(t *testing.T) {
	stations := hertzGRTestStations(t)
	kefalonia := stations["GR002"]

	t.Run("lists the summer season's hours", func(t *testing.T) {
		hours := kefalonia.openingHours("2026-07-01")
		if len(hours) != 7 || hours[0] != (OpeningHoursItem{Day: "Sunday", OpenTime: "08:30", CloseTime: "21:00"}) {
			t.Errorf("hours = %+v", hours)
		}
	})

	t.Run("lists a split day twice in the winter season", func(t *testing.T) {
		hours := kefalonia.openingHours("2026-12-10")
		want := []OpeningHoursItem{
			{Day: "Sunday", OpenTime: "09:00", CloseTime: "14:00"},
			{Day: "Sunday", OpenTime: "18:00", CloseTime: "20:00"},
		}
		if len(hours) != 14 || !reflect.DeepEqual(hours[:2], want) {
			t.Errorf("hours = %+v, want two rows a day", hours)
		}
	})

	t.Run("marks a day the season doesn't open as closed", func(t *testing.T) {
		summer := kefalonia.Schedules[0]
		summer.Times = summer.Times[1:] // drop Sunday
		hours := hertzGRStation{Schedules: []hertzGRSchedule{summer}}.openingHours("2026-07-01")
		if hours[0] != (OpeningHoursItem{Day: "Sunday", OpenTime: "Closed"}) {
			t.Errorf("Sunday = %+v, want closed", hours[0])
		}
	})

	t.Run("is empty outside every season", func(t *testing.T) {
		if hours := kefalonia.openingHours("2030-01-01"); hours != nil {
			t.Errorf("hours = %+v, want none", hours)
		}
	})
}

func TestHertzGRFleetImages(t *testing.T) {
	var feed hertzGRFleetFeed
	err := xml.Unmarshal([]byte(`<Fleet><Country Code="gr"><Brand Code="ZE"><Agent Code="B2B">
		<Car Group = "A" AcrissCode="MDMR" SmallImage="CarImages/NewWebsiteGR/2025/hyundai-i10.jpg" />
	</Agent></Brand></Country></Fleet>`), &feed)
	if err != nil {
		t.Fatal(err)
	}
	images := hertzGRFleetImages(feed)

	if got := images[hertzGRFleetKey("GR", "ZE", "A")]; got != "https://images.autohellas.gr/CarImages/NewWebsiteGR/2025/hyundai-i10.jpg" {
		t.Errorf("image = %q", got)
	}

	t.Run("fills a car whose response has no picture", func(t *testing.T) {
		p := hertzGRTestParams()
		quotes := hertzGRQuotes(hertzGRTestRates[:1], p)
		noPicture := strings.ReplaceAll(hertzGRSameStationFixture, "<PictureURL>http://images.autohellas.gr/CarImages/i10.gif</PictureURL>", "")
		resp := buildHertzGRAvailability(p, 7, quotes, []*hertzGRVehAvailRateRS{hertzGRTestResponse(t, noPicture)}, nil, images)

		if got := resp.AvailableVehicles[0].CarDetails.ImageURL; got != images[hertzGRFleetKey("GR", "ZE", "A")] {
			t.Errorf("image = %q, want the fleet picture", got)
		}
	})
}
