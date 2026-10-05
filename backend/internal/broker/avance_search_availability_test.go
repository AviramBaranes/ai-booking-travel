package broker

import (
	"slices"
	"strings"
	"testing"
)

// The fixtures in avance_fixtures_test.go are live Wheelsys responses for these parameters.
var avanceTestParams = SearchAvailabilityParams{
	CountryCode:     "GR",
	PickupLocation:  "022",
	DropoffLocation: "022",
	PickupDate:      "2026-11-15",
	DropoffDate:     "2026-11-22",
	PickupTime:      "10:00",
	DropoffTime:     "10:00",
	DriverAge:       35,
}

// avanceTestErpDayCharges mirrors booking/config.cue.
var avanceTestErpDayCharges = map[string]float64{avanceCDPStandard: 9, avanceCDPZeroExcess: 3}

func avanceTestQuote(t *testing.T, fixture string) *avanceQuoteXML {
	t.Helper()
	q, err := parseAvanceQuote([]byte(fixture))
	if err != nil {
		t.Fatalf("parse quote: %v", err)
	}
	return q
}

func avanceTestReference(t *testing.T) avanceReference {
	t.Helper()
	stations, err := parseAvanceStations([]byte(avanceFixtureStations))
	if err != nil {
		t.Fatalf("parse stations: %v", err)
	}
	options, err := parseAvanceOptions([]byte(avanceFixtureOptions))
	if err != nil {
		t.Fatalf("parse options: %v", err)
	}
	return avanceReference{stations: stations, options: options}
}

func avancePlanByName(t *testing.T, v AvailableVehicle, name string) Plan {
	t.Helper()
	for _, p := range v.Plans {
		if p.PlanName == name {
			return p
		}
	}
	t.Fatalf("group %s has no %s plan", v.CarDetails.CarGroup, name)
	return Plan{}
}

func TestAvanceBuildAvailability(t *testing.T) {
	vch := avanceTestQuote(t, avanceFixtureQuoteVCH)
	fdw0 := avanceTestQuote(t, avanceFixtureQuoteFDW0)
	resp := buildAvanceAvailability(avanceTestParams, 7, avanceTestErpDayCharges, []*avanceQuoteXML{vch, fdw0}, avanceTestReference(t))

	byGroup := make(map[string]AvailableVehicle, len(resp.AvailableVehicles))
	groups := make([]string, 0, len(resp.AvailableVehicles))
	for _, v := range resp.AvailableVehicles {
		byGroup[v.CarDetails.CarGroup] = v
		groups = append(groups, v.CarDetails.CarGroup)
	}

	t.Run("sells only available groups", func(t *testing.T) {
		if want := []string{"A", "D", "E-GQ", "GM3"}; !slices.Equal(groups, want) {
			t.Errorf("groups = %v, want %v (FAD is on request)", groups, want)
		}
	})

	t.Run("every vehicle carries a standard and a gold plan", func(t *testing.T) {
		for _, v := range resp.AvailableVehicles {
			if len(v.Plans) != 2 || v.Plans[0].PlanID != 1 || v.Plans[1].PlanID != 2 {
				t.Errorf("group %s plans = %+v, want Standard (1) then Gold (2)", v.CarDetails.CarGroup, v.Plans)
			}
			if v.Broker != BrokerAvance {
				t.Errorf("group %s broker = %q", v.CarDetails.CarGroup, v.Broker)
			}
		}
	})

	t.Run("prices the standard plan from the prepaid amount", func(t *testing.T) {
		got := avancePlanByName(t, byGroup["A"], "Standard")
		want := Plan{
			PlanID:                 1,
			PlanName:               "Standard",
			Price:                  77.66,
			ChargedErpPriceWithVat: 63, // 7 days x 9
			RateQualifier:          vch.ID + ":A",
			SupplierName:           "Avance",
			SupplierCode:           "AVANCE",
			Deposit:                868, // 700 + 24% VAT
			DepositCurrency:        "EUR",
			Excess:                 868,
			ExcessCurrency:         "EUR",
			TheftExcess:            868,
			TheftExcessCurrency:    "EUR",
		}
		got.Info = nil
		if !planEqual(got, want) {
			t.Errorf("A Standard = %+v\nwant          %+v", got, want)
		}
	})

	t.Run("prices the zero excess plan", func(t *testing.T) {
		got := avancePlanByName(t, byGroup["A"], "Gold")
		if got.Price != 201.40 || got.Excess != 0 || got.Deposit != 186 || got.TheftExcess != 868 {
			t.Errorf("A Gold price/excess/deposit/theft = %v/%d/%d/%d, want 201.4/0/186/868",
				got.Price, got.Excess, got.Deposit, got.TheftExcess)
		}
		if got.ChargedErpPriceWithVat != 21 {
			t.Errorf("A Gold ERP = %v, want 21 (7 days x 3)", got.ChargedErpPriceWithVat)
		}
		if got.RateQualifier != fdw0.ID+":A" {
			t.Errorf("A Gold RateQualifier = %q", got.RateQualifier)
		}
	})

	t.Run("finds the FDW0 deposit when the waiver has no excess attribute", func(t *testing.T) {
		if got := avancePlanByName(t, byGroup["GM3"], "Gold").Deposit; got != 186 {
			t.Errorf("GM3 Gold deposit = %d, want 186", got)
		}
	})

	t.Run("ignores the zero deposit on the theft waiver", func(t *testing.T) {
		gold := avancePlanByName(t, byGroup["E-GQ"], "Gold")
		if gold.Deposit != 186 || gold.TheftExcess != 1178 {
			t.Errorf("E-GQ Gold deposit/theft = %d/%d, want 186/1178", gold.Deposit, gold.TheftExcess)
		}
		std := avancePlanByName(t, byGroup["E-GQ"], "Standard")
		if std.Deposit != 1178 || std.Excess != 1178 {
			t.Errorf("E-GQ Standard deposit/excess = %d/%d, want 1178/1178", std.Deposit, std.Excess)
		}
	})

	t.Run("rate qualifiers are unique per car and plan", func(t *testing.T) {
		seen := make(map[string]bool)
		for _, v := range resp.AvailableVehicles {
			for _, p := range v.Plans {
				key := p.RateQualifier + "|" + p.SupplierCode + "|" + p.PlanName
				if seen[key] {
					t.Errorf("duplicate plan identity %q: booking would reserve the wrong car", key)
				}
				seen[key] = true
			}
		}
	})

	t.Run("maps car details", func(t *testing.T) {
		got := byGroup["E-GQ"].CarDetails
		want := CarDetails{
			Model: "Geely EX5 Electric", CarGroup: "E-GQ", ImageURL: got.ImageURL, SupplierName: "Avance",
			CarType: "SUV", Acriss: "DFAE", FullAcriss: "DFAE",
			HasAC: true, IsAutoGear: true, IsElectric: true, Seats: 5, Bags: 2, Doors: 5,
		}
		if got != want {
			t.Errorf("E-GQ = %+v\nwant  %+v", got, want)
		}
		if byGroup["A"].CarDetails.IsAutoGear || !byGroup["GM3"].CarDetails.IsAutoGear {
			t.Error("gearbox: want A manual (MBMR) and GM3 automatic (MDAR)")
		}
	})

	t.Run("badges the pickup station type", func(t *testing.T) {
		if got := byGroup["A"].LocationDetails.LocationType; got != "Airport" {
			t.Errorf("LocationType = %q, want Airport", got)
		}
	})

	t.Run("adds no station fees for a driver of standard age", func(t *testing.T) {
		if got := byGroup["A"].PriceDetails; got.Fees != (Fees{}) || got.Currency != "EUR" {
			t.Errorf("PriceDetails = %+v, want EUR and no fees", got)
		}
	})

	t.Run("spells out the excess in the plan terms", func(t *testing.T) {
		std := avancePlanByName(t, byGroup["A"], "Standard").Info
		if !slices.Contains(std, "Excess: 868 EUR") || strings.Contains(strings.Join(std, "|"), "Theft excess") {
			t.Errorf("Standard info = %q, want the excess and no separate theft line", std)
		}
		gold := avancePlanByName(t, byGroup["A"], "Gold").Info
		for _, want := range []string{
			"Excess: 0 EUR",
			"Theft excess: 868 EUR",
			"Damage to the underside, interior, wheels and glass is not covered",
		} {
			if !slices.Contains(gold, want) {
				t.Errorf("Gold info = %q, missing %q", gold, want)
			}
		}
	})

	t.Run("describes the supplier once", func(t *testing.T) {
		if len(resp.SuppliersInfo) != 1 || resp.SuppliersInfo[0].Name != "Avance" {
			t.Fatalf("SuppliersInfo = %+v, want one Avance entry", resp.SuppliersInfo)
		}
		pickup := resp.SuppliersInfo[0].PickupDetails
		if pickup.Address != "Desk in Arrivals Hall 30, Spata1, 19019" || pickup.PhoneNumber != "+30 2103533088" {
			t.Errorf("pickup = %+v", pickup)
		}
		if len(pickup.OpeningHours) != 7 || pickup.OpeningHours[0] != (OpeningHoursItem{Day: "Monday", OpenTime: "00:00", CloseTime: "23:59"}) {
			t.Errorf("opening hours = %+v", pickup.OpeningHours)
		}
	})

	t.Run("offers only seats, snow chains and additional drivers", func(t *testing.T) {
		want := []AddOn{
			{ID: 9001, Name: "Baby Seat", Price: 49, Currency: "EUR", AllowedQuantity: 3, Period: "Per Rental"},
			{ID: 9002, Name: "Booster Seat", Price: 49, Currency: "EUR", AllowedQuantity: 3, Period: "Per Rental"},
			{ID: 9003, Name: "Snow Chains", Price: 50, Currency: "EUR", AllowedQuantity: 1, Period: "Per Rental"},
			{ID: 9004, Name: "Additional Driver", Price: 28, Currency: "EUR", AllowedQuantity: 2, Period: "Per Rental"},
		}
		if got := resp.SuppliersInfo[0].AddOns; !slices.Equal(got, want) {
			t.Errorf("AddOns = %+v\nwant    %+v", got, want)
		}
	})

	t.Run("lists what each plan includes", func(t *testing.T) {
		incs := resp.SuppliersInfo[0].Inclusions
		if len(incs) != 2 || incs[0].ProductName != "Standard" || incs[1].ProductName != "Gold" {
			t.Fatalf("Inclusions = %+v, want Standard then Gold", incs)
		}
		for _, inc := range incs {
			for _, want := range []string{"Collision Damage Waiver", "Theft Waiver with Excess", "Unlimited mileage", "First additional driver free"} {
				if !slices.Contains(inc.ProductInclusions, want) {
					t.Errorf("%s inclusions = %q, missing %q", inc.ProductName, inc.ProductInclusions, want)
				}
			}
		}
		if slices.Contains(incs[0].ProductInclusions, "Full Damage Waiver") || !slices.Contains(incs[1].ProductInclusions, "Full Damage Waiver") {
			t.Errorf("want Full Damage Waiver in Gold only: %+v", incs)
		}
	})
}

func TestAvanceYoungDriverFee(t *testing.T) {
	young := avanceTestQuote(t, avanceFixtureQuoteYoung)
	resp := buildAvanceAvailability(avanceTestParams, 7, avanceTestErpDayCharges, []*avanceQuoteXML{young, nil}, avanceTestReference(t))

	a := resp.AvailableVehicles[0]
	if a.CarDetails.CarGroup != "A" {
		t.Fatalf("first group = %s, want A", a.CarDetails.CarGroup)
	}
	// The young driver fee is settled at the station, so it stays out of the prepaid price.
	if a.PriceDetails.Fees.YoungDriverFee != 70 || a.PriceDetails.Fees.YoungDriverFeeCurrency != "EUR" {
		t.Errorf("Fees = %+v, want a 70 EUR young driver fee", a.PriceDetails.Fees)
	}
	if got := a.Plans[0].Price; got != 77.66 {
		t.Errorf("Price = %v, want 77.66 unchanged by the young driver fee", got)
	}
}

func TestAvanceBuildAvailabilityPartial(t *testing.T) {
	ref := avanceTestReference(t)

	t.Run("a failed product leaves the other plan", func(t *testing.T) {
		vch := avanceTestQuote(t, avanceFixtureQuoteVCH)
		resp := buildAvanceAvailability(avanceTestParams, 7, avanceTestErpDayCharges, []*avanceQuoteXML{vch, nil}, ref)
		for _, v := range resp.AvailableVehicles {
			if len(v.Plans) != 1 || v.Plans[0].PlanName != "Standard" {
				t.Errorf("group %s plans = %+v, want Standard only", v.CarDetails.CarGroup, v.Plans)
			}
		}
		if incs := resp.SuppliersInfo[0].Inclusions; len(incs) != 1 || incs[0].ProductName != "Standard" {
			t.Errorf("Inclusions = %+v, want Standard only", incs)
		}
	})

	t.Run("nothing available returns an empty response", func(t *testing.T) {
		resp := buildAvanceAvailability(avanceTestParams, 7, avanceTestErpDayCharges, []*avanceQuoteXML{nil, nil}, ref)
		if len(resp.AvailableVehicles) != 0 || len(resp.SuppliersInfo) != 0 {
			t.Errorf("resp = %+v, want empty", resp)
		}
	})
}

func TestParseAvanceQuoteError(t *testing.T) {
	_, err := parseAvanceQuote([]byte(avanceFixtureQuoteError))
	if err == nil || !strings.Contains(err.Error(), "ERR/102") {
		t.Errorf("err = %v, want the ERR/102 error returned inside the 200 response", err)
	}
}

func TestAcrissHasAC(t *testing.T) {
	for code, want := range map[string]bool{
		"MBMR": true,  // unspecified fuel, air
		"SVMD": true,  // diesel, air
		"DFAE": true,  // electric, air
		"EDMN": false, // unspecified fuel, no air
		"CDMQ": false, // diesel, no air
		"ABC":  false, // too short to tell
	} {
		if got := acrissHasAC(code); got != want {
			t.Errorf("acrissHasAC(%q) = %v, want %v", code, got, want)
		}
	}
}

// planEqual compares plans field by field; Plan holds a slice so it is not comparable with ==.
func planEqual(a, b Plan) bool {
	return a.PlanID == b.PlanID && a.PlanName == b.PlanName && a.Price == b.Price &&
		a.BrokerErpPrice == b.BrokerErpPrice && a.ChargedErpPriceWithVat == b.ChargedErpPriceWithVat &&
		a.RateQualifier == b.RateQualifier && a.SupplierName == b.SupplierName && a.SupplierCode == b.SupplierCode &&
		a.Deposit == b.Deposit && a.DepositCurrency == b.DepositCurrency &&
		a.Excess == b.Excess && a.ExcessCurrency == b.ExcessCurrency &&
		a.TheftExcess == b.TheftExcess && a.TheftExcessCurrency == b.TheftExcessCurrency &&
		slices.Equal(a.Info, b.Info)
}

func TestAvanceAdditionalDriverNotFirstFree(t *testing.T) {
	// Booking adds the free first driver to the paid ones, so a quote without it offers none.
	fixture := strings.ReplaceAll(avanceFixtureQuoteVCH, `code="ADD" rate="2800" firstfree="true"`, `code="ADD" rate="2800" firstfree="false"`)
	resp := buildAvanceAvailability(avanceTestParams, 7, avanceTestErpDayCharges, []*avanceQuoteXML{avanceTestQuote(t, fixture), nil}, avanceTestReference(t))

	for _, ao := range resp.SuppliersInfo[0].AddOns {
		if ao.ID == avanceAddOnIDs[avanceAdditionalDriver] {
			t.Errorf("AddOns has %+v, want no additional driver", ao)
		}
	}
	if slices.Contains(resp.SuppliersInfo[0].Inclusions[0].ProductInclusions, "First additional driver free") {
		t.Error("Inclusions has a free additional driver the quote doesn't give")
	}
}

func TestAvanceTermsReturned(t *testing.T) {
	resp := buildAvanceAvailability(avanceTestParams, 7, avanceTestErpDayCharges, []*avanceQuoteXML{avanceTestQuote(t, avanceFixtureQuoteVCH), nil}, avanceTestReference(t))

	terms := resp.SuppliersInfo[0].TermsAndConditions
	if len(terms) != len(avanceTerms) {
		t.Fatalf("terms = %d, want all %d from the Terms sheet", len(terms), len(avanceTerms))
	}
	for _, term := range terms {
		if term.Title == "" || term.HtmlContent == "" {
			t.Errorf("term %+v has an empty title or content", term)
		}
	}
}
