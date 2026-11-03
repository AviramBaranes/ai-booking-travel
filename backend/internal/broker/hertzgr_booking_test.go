package broker

import (
	"errors"
	"strings"
	"testing"
)

func hertzGRTestBooking() BookingParams {
	return BookingParams{
		RateQualifier:   "ITRHABT:A",
		SupplierCode:    "ZE",
		PlanID:          "1",
		PickupLocation:  "GR030",
		DropoffLocation: "GR030",
		QuotedPrice:     164.01,
		DriverFirstName: "JOHN",
		DriverLastName:  "LAST",
		FlightNumber:    "A3 960",
		DriverAge:       "30",
		PickupDate:      "2026-12-10",
		PickupTime:      "10:00",
		DropoffDate:     "2026-12-17",
		DropoffTime:     "10:00",
	}
}

// hertzGRBookingServer answers reservations with reservation and cancellations with cancel.
func hertzGRBookingServer(t *testing.T, reservation, cancel string) (*HertzGR, func() []string) {
	t.Helper()
	h, _, bodies := hertzGRTestServer(t, func(body string) string {
		if strings.Contains(body, "OTA_VehCancelRQ") {
			return cancel
		}
		return reservation
	})
	return h, bodies
}

func cancelCalls(bodies []string) []string {
	var out []string
	for _, b := range bodies {
		if strings.Contains(b, "OTA_VehCancelRQ") {
			out = append(out, b)
		}
	}
	return out
}

func TestHertzGRBook(t *testing.T) {
	t.Run("books the plan's group and returns Autohellas' number", func(t *testing.T) {
		h, bodies := hertzGRBookingServer(t, hertzGRReservationFixture, hertzGRCancelFixture)

		res, err := h.Book(hertzGRTestBooking())
		if err != nil {
			t.Fatalf("Book: %v", err)
		}
		if res.ConfirmationNumber != "GR400454846" {
			t.Errorf("confirmation = %q, want GR400454846", res.ConfirmationNumber)
		}
		if n := len(cancelCalls(bodies())); n != 0 {
			t.Errorf("cancel calls = %d, want none", n)
		}

		rq := bodies()[0]
		for _, want := range []string{
			`<VehicleReservationRequest xmlns="http://www.hertz.gr/xml"><OTA_VehResRQ`,
			`<RequestorID Type="4" ID="CCN1"><CompanyName Code="VC1"></CompanyName></RequestorID>`,
			`<RequestorID Type="16" ID="AIBT`,
			`<VehResRQCore Status="All">`,
			`<PickUpLocation LocationCode="GR030">`,
			`<PersonName><GivenName>JOHN</GivenName><Surname>LAST</Surname></PersonName>`,
			`<VehPref VendorCarType="A"></VehPref>`,
			`<DriverType Age="30"></DriverType>`,
			`<RateQualifier RateQualifier="ITRHABT"></RateQualifier>`,
			`<ArrivalDetails TransportationCode="14" Number="A3 960"></ArrivalDetails>`,
		} {
			if !strings.Contains(rq, want) {
				t.Errorf("request lacks %s:\n%s", want, rq)
			}
		}
	})

	t.Run("cancels a reservation priced above the quote", func(t *testing.T) {
		h, bodies := hertzGRBookingServer(t, hertzGRReservationFixture, hertzGRCancelFixture)
		p := hertzGRTestBooking()
		p.QuotedPrice = 150

		_, err := h.Book(p)
		if err == nil || errors.Is(err, ErrBookingLeftOpen) || !strings.Contains(err.Error(), "above the 150.00 quoted") {
			t.Errorf("err = %v, want a cancelled-for-price error", err)
		}
		cancels := cancelCalls(bodies())
		if len(cancels) != 1 || !strings.Contains(cancels[0], `<UniqueID Type="14" ID="GR400454846">`) {
			t.Errorf("cancel calls = %v, want one for GR400454846", cancels)
		}
	})

	t.Run("cancels a reservation that is not confirmed", func(t *testing.T) {
		onHold := strings.Replace(hertzGRReservationFixture, `ReservationStatus="Committed"`, `ReservationStatus="On Hold"`, 1)
		h, bodies := hertzGRBookingServer(t, onHold, hertzGRCancelFixture)

		if _, err := h.Book(hertzGRTestBooking()); err == nil || !strings.Contains(err.Error(), `"On Hold"`) {
			t.Errorf("err = %v, want an on-hold error", err)
		}
		if n := len(cancelCalls(bodies())); n != 1 {
			t.Errorf("cancel calls = %d, want 1", n)
		}
	})

	t.Run("reports a reservation it could not cancel as left open", func(t *testing.T) {
		onHold := strings.Replace(hertzGRReservationFixture, `ReservationStatus="Committed"`, `ReservationStatus="On Hold"`, 1)
		refused := strings.NewReplacer("CODE", "97", "TEXT", "Booking not found").Replace(hertzGRCancelErrorFixture)
		h, _ := hertzGRBookingServer(t, onHold, refused)

		if _, err := h.Book(hertzGRTestBooking()); !errors.Is(err, ErrBookingLeftOpen) {
			t.Errorf("err = %v, want ErrBookingLeftOpen", err)
		}
	})

	t.Run("fails on an OTA error without cancelling", func(t *testing.T) {
		failed := strings.NewReplacer("VehicleAvailRateResult", "VehicleReservationResult", "OTA_VehAvailRateRS", "OTA_VehResRS",
			"CODE", "280", "TEXT", "Requested car category or rate code not available").Replace(hertzGRErrorFixture)
		h, bodies := hertzGRBookingServer(t, failed, hertzGRCancelFixture)

		if _, err := h.Book(hertzGRTestBooking()); err == nil || !strings.Contains(err.Error(), "280") {
			t.Errorf("err = %v, want the 280 error", err)
		}
		if n := len(cancelCalls(bodies())); n != 0 {
			t.Errorf("cancel calls = %d, want none", n)
		}
	})

	t.Run("rejects a malformed rate qualifier", func(t *testing.T) {
		p := hertzGRTestBooking()
		p.RateQualifier = "ITRHABT"
		if _, err := NewHertzGR(HertzGRConfig{}).Book(p); err == nil {
			t.Error("err = nil, want an error")
		}
	})
}

func TestHertzGRCancel(t *testing.T) {
	t.Run("cancels by reservation number and surname", func(t *testing.T) {
		h, bodies := hertzGRBookingServer(t, "", hertzGRCancelFixture)

		if err := h.Cancel("GR400454846", "LAST", "ZE"); err != nil {
			t.Fatalf("Cancel: %v", err)
		}
		rq := bodies()[0]
		for _, want := range []string{
			`<VehicleCancelReservationRequest xmlns="http://www.hertz.gr/xml"><OTA_VehCancelRQ`,
			`<VehCancelRQCore CancelType="Cancel">`,
			`<UniqueID Type="14" ID="GR400454846"></UniqueID>`,
			`<PersonName><Surname>LAST</Surname></PersonName>`,
		} {
			if !strings.Contains(rq, want) {
				t.Errorf("request lacks %s:\n%s", want, rq)
			}
		}
	})

	t.Run("treats an already cancelled reservation as cancelled", func(t *testing.T) {
		already := strings.NewReplacer("CODE", "95", "TEXT", "Booking already cancelled").Replace(hertzGRCancelErrorFixture)
		h, _ := hertzGRBookingServer(t, "", already)
		if err := h.Cancel("GR400454846", "LAST", "ZE"); err != nil {
			t.Errorf("err = %v, want nil", err)
		}
	})

	t.Run("reports any other refusal", func(t *testing.T) {
		notFound := strings.NewReplacer("CODE", "97", "TEXT", "Booking not found").Replace(hertzGRCancelErrorFixture)
		h, _ := hertzGRBookingServer(t, "", notFound)
		if err := h.Cancel("GR400454846", "LAST", "ZE"); err == nil || !strings.Contains(err.Error(), "Booking not found") {
			t.Errorf("err = %v, want the refusal", err)
		}
	})
}

func TestHertzGRReference(t *testing.T) {
	a, err := hertzGRReference()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := hertzGRReference()
	if len(a) != 20 || !strings.HasPrefix(a, "AIBT") || a == b {
		t.Errorf("references %q, %q: want two different 20-character AIBT references", a, b)
	}
}
