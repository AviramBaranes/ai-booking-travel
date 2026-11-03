package booking_handlers

import (
	"encoding/json"
	"testing"

	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	availability "encore.app/services/booking/handlers/availability"
)

func TestFindSupplierInfo(t *testing.T) {
	snapshotOf := func(t *testing.T, suppliers ...broker.SupplierInfo) db.AvailablePlansSnapshot {
		t.Helper()
		b, err := json.Marshal(suppliers)
		if err != nil {
			t.Fatal(err)
		}
		return db.AvailablePlansSnapshot{SuppliersInfo: b}
	}
	flexHertz := broker.SupplierInfo{Name: "Hertz", Broker: broker.BrokerFlex, PickupDetails: broker.StationInfo{Address: "flex"}}
	ourHertz := broker.SupplierInfo{Name: "Hertz", Broker: broker.BrokerHertzGR, PickupDetails: broker.StationInfo{Address: "hertz_gr"}}

	t.Run("tells apart suppliers of the same name by broker", func(t *testing.T) {
		got := FindSupplierInfo(snapshotOf(t, flexHertz, ourHertz), availability.PlanPriceDetails{SupplierName: "Hertz", Broker: broker.BrokerHertzGR})
		if got.PickupDetails.Address != "hertz_gr" {
			t.Errorf("found %+v, want Hertz GR's", got)
		}
	})

	t.Run("matches a snapshot stored before suppliers recorded their broker", func(t *testing.T) {
		legacy := broker.SupplierInfo{Name: "Avance"}
		got := FindSupplierInfo(snapshotOf(t, legacy), availability.PlanPriceDetails{SupplierName: "Avance", Broker: broker.BrokerAvance})
		if got.Name != "Avance" {
			t.Errorf("found %+v, want the legacy entry", got)
		}
	})

	t.Run("returns nothing for another broker's supplier", func(t *testing.T) {
		got := FindSupplierInfo(snapshotOf(t, flexHertz), availability.PlanPriceDetails{SupplierName: "Hertz", Broker: broker.BrokerHertzGR})
		if got.Name != "" {
			t.Errorf("found %+v, want none", got)
		}
	})
}
