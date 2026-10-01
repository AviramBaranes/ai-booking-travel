package broker

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// avanceTestServer answers each Wheelsys page with the given body and records the queries it got.
func avanceTestServer(t *testing.T, pages map[string]string) (*Avance, *[]url.URL) {
	t.Helper()

	savedCreds, savedErr := avanceCreds, avanceCredsErr
	avanceCreds = avanceCredentials{accountNo: "10268", linkCode: "TBAVA", agentCode: "ABTAV"}
	avanceCredsErr = nil
	t.Cleanup(func() { avanceCreds, avanceCredsErr = savedCreds, savedErr })

	var calls []url.URL
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, *r.URL)
		for page, body := range pages {
			if strings.HasSuffix(r.URL.Path, "/"+page+"_TBAVA.html") {
				_, _ = w.Write([]byte(body))
				return
			}
		}
		t.Errorf("unexpected request %s", r.URL)
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)

	return NewAvance(AvanceConfig{BaseURL: srv.URL}), &calls
}

var avanceTestBooking = BookingParams{
	RateQualifier:   "01286b5b-ad6a-49ae-b74e-ae39cfdf3162:E-GQ",
	SupplierCode:    avanceSupplierCode,
	PlanID:          "2",
	PickupLocation:  "022",
	DropoffLocation: "029",
	SelectedAddOns:  []SelectAddOn{{ID: 9001, Quantity: 1}, {ID: 9004, Quantity: 2}, {ID: 9002, Quantity: 0}},
	DriverFirstName: "KONSTANTINOS ALEXANDROS MICHAEL",
	DriverLastName:  "PAPADOPOULOS-KARAMANLIS",
	FlightNumber:    "LY541",
	DriverAge:       "35",
	PickupDate:      "2026-10-15",
	DropoffDate:     "2026-10-22",
	PickupTime:      "10:00",
	DropoffTime:     "09:30",
}

func TestAvanceBook(t *testing.T) {
	t.Run("confirmed", func(t *testing.T) {
		a, calls := avanceTestServer(t, map[string]string{"new-res": avanceFixtureNewResConfirmed})

		res, err := a.Book(avanceTestBooking)
		if err != nil {
			t.Fatalf("Book: %v", err)
		}
		if res.ConfirmationNumber != "5RLDDC" {
			t.Errorf("ConfirmationNumber = %q, want 5RLDDC", res.ConfirmationNumber)
		}
		if len(*calls) != 1 {
			t.Fatalf("calls = %d, want 1", len(*calls))
		}

		call := (*calls)[0]
		if want := "/10268/link/v3/new-res_TBAVA.html"; call.Path != want {
			t.Errorf("path = %s, want %s", call.Path, want)
		}

		q := call.Query()
		for param, want := range map[string]string{
			"agent":          "ABTAV",
			"DATE_FROM":      "15/10/2026",
			"TIME_FROM":      "10:00",
			"DATE_TO":        "22/10/2026",
			"TIME_TO":        "09:30",
			"PICKUP_STATION": "022",
			"RETURN_STATION": "029",
			"CDP":            avanceCDPZeroExcess,
			"GROUP":          "E-GQ",
			"QUOTEREF":       "01286b5b-ad6a-49ae-b74e-ae39cfdf3162",
			"VOUCHERNO":      avanceVoucherNo(avanceTestBooking.RateQualifier),
			"CUSTFIRST_NAME": "KONSTANTINOS ALEXANDROS MICHAE",
			"CUSTLAST_NAME":  "PAPADOPOULOS-KARAMANLIS",
			"PICKUP_INFO":    "LY541",
			"DRIVERAGE":      "35",
			"CS":             "1",
			// Two paid drivers plus the free first one.
			"ADD": "3",
		} {
			if got := q.Get(param); got != want {
				t.Errorf("%s = %q, want %q", param, got, want)
			}
		}
		if q.Has("BS") {
			t.Errorf("BS = %q, want an add-on with no quantity left out", q.Get("BS"))
		}
	})

	t.Run("on request is cancelled", func(t *testing.T) {
		a, calls := avanceTestServer(t, map[string]string{
			"new-res":    avanceFixtureNewResOnRequest,
			"cancel-res": avanceFixtureCancelRes,
		})

		if _, err := a.Book(avanceTestBooking); err == nil {
			t.Fatal("Book succeeded, want an error for an on-request reservation")
		}
		if len(*calls) != 2 || !strings.Contains((*calls)[1].Path, "cancel-res") {
			t.Fatalf("calls = %v, want new-res then cancel-res", *calls)
		}
		if got := (*calls)[1].Query().Get("irn"); got != "5NS061" {
			t.Errorf("cancelled irn = %q, want 5NS061", got)
		}
	})

	t.Run("on request that can't be cancelled is left open", func(t *testing.T) {
		a, _ := avanceTestServer(t, map[string]string{
			"new-res":    avanceFixtureNewResOnRequest,
			"cancel-res": `<response><reservation irn="5NS061" status="ERR/999" /></response>`,
		})

		_, err := a.Book(avanceTestBooking)
		if !errors.Is(err, ErrBookingLeftOpen) || !strings.Contains(err.Error(), "5NS061") {
			t.Errorf("err = %v, want ErrBookingLeftOpen naming 5NS061", err)
		}
	})

	for name, fixture := range map[string]string{
		"error in status":   avanceFixtureResStatusError,
		"error in a block":  avanceFixtureResErrorsBlock,
		"unparsable answer": "Service Unavailable",
	} {
		t.Run(name, func(t *testing.T) {
			a, calls := avanceTestServer(t, map[string]string{"new-res": fixture})

			if _, err := a.Book(avanceTestBooking); err == nil {
				t.Fatal("Book succeeded, want an error")
			}
			if len(*calls) != 1 {
				t.Errorf("calls = %d, want no cancel-res without a reservation", len(*calls))
			}
		})
	}

	t.Run("bad snapshot values fail before calling Wheelsys", func(t *testing.T) {
		for name, p := range map[string]BookingParams{
			"rate qualifier": func() BookingParams { p := avanceTestBooking; p.RateQualifier = "no-group"; return p }(),
			"plan id":        func() BookingParams { p := avanceTestBooking; p.PlanID = "7"; return p }(),
			"add-on id": func() BookingParams {
				p := avanceTestBooking
				p.SelectedAddOns = []SelectAddOn{{ID: 1, Quantity: 1}}
				return p
			}(),
		} {
			if _, err := avanceBookingQuery(p); err == nil {
				t.Errorf("%s: avanceBookingQuery succeeded, want an error", name)
			}
		}
	})
}

func TestAvanceVoucherNo(t *testing.T) {
	a := avanceVoucherNo("quote-1:A")
	if len(a) > 25 || !strings.HasPrefix(a, "AIBT") {
		t.Errorf("voucher number %q, want AIBT and at most 25 characters", a)
	}
	if a != avanceVoucherNo("quote-1:A") {
		t.Error("voucher number is not deterministic")
	}
	if a == avanceVoucherNo("quote-1:B") {
		t.Error("two groups of one quote share a voucher number")
	}
}

func TestAvanceCancel(t *testing.T) {
	for name, tc := range map[string]struct {
		fixture string
		wantErr bool
	}{
		"cancelled":         {fixture: avanceFixtureCancelRes},
		"already cancelled": {fixture: avanceFixtureCancelAlreadyCancelled},
		"car collected":     {fixture: `<response><reservation irn="5RLDDC" status="ERR/105" /></response>`, wantErr: true},
	} {
		t.Run(name, func(t *testing.T) {
			a, calls := avanceTestServer(t, map[string]string{"cancel-res": tc.fixture})

			err := a.Cancel("5RLDDC", "", "")
			if (err != nil) != tc.wantErr {
				t.Errorf("Cancel error = %v, want error %v", err, tc.wantErr)
			}
			if got := (*calls)[0].Query().Get("irn"); got != "5RLDDC" {
				t.Errorf("irn = %q, want 5RLDDC", got)
			}
		})
	}
}

func TestAvanceGenerateVoucher(t *testing.T) {
	d := &VoucherData{
		BookingReferenceID:  "5RLDDC",
		CustomerName:        "Mr ISRAEL ISRAELI",
		FlightNumber:        "LY541",
		PickupLoc:           "Athens International Airport",
		PickupBranch:        "Athens International Airport, Arrivals",
		PickupPhone:         "+30 210 3538700",
		PickupInstructions:  "Meet our representative at the arrivals exit",
		PickupDate:          "2026-10-15",
		PickupTime:          "10:00",
		DropoffLoc:          "Athens International Airport",
		DropoffInstructions: "Leave the car at the P4 parking",
		DropoffDate:         "2026-10-22",
		DropoffTime:         "10:00",
		CarGroupDesc:        "E-GQ",
		LeadModel:           "Toyota Yaris",
		Passengers:          5,
		Suitcases:           2,
		PrepaidIncludes:     []string{"Collision Damage Waiver", "Unlimited mileage"},
		OptionalServices:    []string{"Baby Seat × 1"},
		PayAtPickup:         []string{"Young driver fee: 70 EUR"},
		Deposit:             186,
		DepositCurrency:     "EUR",
		Excess:              0,
		ExcessCurrency:      "EUR",
		TheftExcess:         744,
		TheftExcessCurrency: "EUR",
	}

	html, err := NewAvance(AvanceConfig{}).GenerateVoucher(d)
	if err != nil {
		t.Fatalf("GenerateVoucher: %v", err)
	}

	for _, want := range []string{
		"5RLDDC", "LY541", "Meet our representative", "210 3538700", "Baby Seat × 1",
		"Young driver fee: 70 EUR", "Leave the car at the P4 parking", "186 EUR", "744 EUR", "wheels", "Toyota Yaris",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("voucher is missing %q", want)
		}
	}

	d.Excess = 1240
	html, err = NewAvance(AvanceConfig{}).GenerateVoucher(d)
	if err != nil {
		t.Fatalf("GenerateVoucher: %v", err)
	}
	if strings.Contains(html, "wheels") {
		t.Error("voucher carries the zero-excess exclusions on a plan with an excess")
	}
}
