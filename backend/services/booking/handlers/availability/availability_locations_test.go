package availability

import (
	"reflect"
	"testing"

	"encore.app/internal/broker"
	"encore.app/services/booking/db"
)

func TestBuildAvailabilityLocations(t *testing.T) {
	const mykonos, athens = 1, 2
	row := func(location int64, b db.Broker, code string, brands ...string) db.GetAllLocationBrokerCodesByLocationIDsRow {
		return db.GetAllLocationBrokerCodesByLocationIDsRow{
			LocationID: location, Broker: b, BrokerLocationID: code, Brands: brands, Enabled: true,
			LocationCountryCode: "GR", LocationName: map[int64]string{mykonos: "Mykonos Airport", athens: "Athens Airport"}[location],
		}
	}
	disabled := row(mykonos, db.BrokerAvance, "AV-MYK")
	disabled.Enabled = false

	rows := []db.GetAllLocationBrokerCodesByLocationIDsRow{
		row(mykonos, db.BrokerHertzGr, "GR221", "ZR", "ZT"),
		row(mykonos, db.BrokerFlex, "JMK01"),
		row(mykonos, db.BrokerHertzGr, "GR001", "ZE"),
		disabled,
		row(athens, db.BrokerHertzGr, "GR423", "ZR"),
		row(athens, db.BrokerHertzGr, "GR030", "FF", "ZE", "ZT"),
	}

	pickups, dropoffs := createBrokersMap(rows, SearchAvailabilityParams{PickupLocationID: mykonos}, athens)
	al := buildAvailabilityLocations(pickups, dropoffs)

	t.Run("keeps every station of a broker with brands, ordered by code", func(t *testing.T) {
		got := al[broker.BrokerHertzGR]
		wantPickup := []broker.Station{{Code: "GR001", Brands: []string{"ZE"}}, {Code: "GR221", Brands: []string{"ZR", "ZT"}}}
		wantDropoff := []broker.Station{{Code: "GR030", Brands: []string{"FF", "ZE", "ZT"}}, {Code: "GR423", Brands: []string{"ZR"}}}
		if !reflect.DeepEqual(got.pickupStations, wantPickup) || !reflect.DeepEqual(got.dropoffStations, wantDropoff) {
			t.Errorf("stations = %+v / %+v", got.pickupStations, got.dropoffStations)
		}
		if got.pickupBrokerLocationID != "GR001" || got.dropoffBrokerLocationID != "GR030" {
			t.Errorf("first stations = %s / %s, want GR001 / GR030", got.pickupBrokerLocationID, got.dropoffBrokerLocationID)
		}
	})

	t.Run("leaves out brokers missing at either end and disabled codes", func(t *testing.T) {
		if _, ok := al[broker.BrokerFlex]; ok {
			t.Error("flex has no Athens code but was kept")
		}
		if _, ok := pickups[db.BrokerAvance]; ok {
			t.Error("a disabled code was kept")
		}
	})
}
