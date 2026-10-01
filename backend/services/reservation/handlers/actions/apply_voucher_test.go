package actions

import (
	"testing"

	"encore.app/services/reservation/db"
)

func TestToVoucherDataBlanksUnreadableDetails(t *testing.T) {
	d, err := toVoucherData(db.Reservation{
		BrokerReservationID: "5RLDDC",
		CarDetails:          []byte(`{"model":"Toyota Yaris"}`),
		PayAtPickup:         []byte(`not json`),
		PickupDetails:       []byte(`{"address":`),
		DropoffDetails:      []byte(`{"address":"Athens International Airport","phoneNumber":"210 3538700"}`),
		Excess:              1240,
		ExcessCurrency:      "EUR",
	})
	if err != nil {
		t.Fatalf("toVoucherData: %v", err)
	}

	if d.Deposit != 0 || len(d.OptionalServices) != 0 || len(d.PayAtPickup) != 0 {
		t.Errorf("pay at pickup = deposit %d, services %v, fees %v, want blank", d.Deposit, d.OptionalServices, d.PayAtPickup)
	}
	if d.PickupBranch != "" {
		t.Errorf("PickupBranch = %q, want blank", d.PickupBranch)
	}
	if d.DropoffBranch != "Athens International Airport" || d.DropoffPhone != "210 3538700" {
		t.Errorf("dropoff = %q %q, want the readable details kept", d.DropoffBranch, d.DropoffPhone)
	}
	if d.LeadModel != "Toyota Yaris" || d.Excess != 1240 {
		t.Errorf("voucher = %+v, want the rest filled", d)
	}
}
