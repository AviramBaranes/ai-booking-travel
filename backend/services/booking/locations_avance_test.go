package booking

import (
	"context"
	"testing"

	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/location"
)

// TestInsertAvanceLocationsMergesAirports covers the two things the Avance import depends on that
// the parser tests cannot reach: that Postgres accepts the new 'avance' enum value, and that an
// Avance airport station attaches to the canonical location an existing broker already created
// rather than duplicating it.
func TestInsertAvanceLocationsMergesAirports(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()

	existing, _ := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{
			Country: "Greece", CountryCode: "GR", Name: "Avance Merge Test Airport",
			City: strPtr("Athens"), Iata: strPtr("ZZZ"),
		},
		db.BrokerFlex, "flex-avance-merge",
	)

	b := &mockBroker{
		name: broker.BrokerAvance,
		pages: []broker.LocationPage{{
			Locations: []broker.Location{
				{ID: "ZZZAP", Name: "Avance Merge Test Airport", Country: "Greece", CountryCode: "GR", City: "Athens", Iata: "ZZZ"},
				{ID: "ZZZDT", Name: "Avance Merge Test Downtown", Country: "Greece", CountryCode: "GR", City: "Athens"},
			},
		}},
	}

	if err := location.NewLocationService(q).InsertLocations(ctx, b, ""); err != nil {
		t.Fatalf("InsertLocations: %v", err)
	}

	t.Run("airport station reuses the existing canonical location", func(t *testing.T) {
		code, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker: db.BrokerAvance, BrokerLocationID: "ZZZAP", LocationID: existing.ID,
		})
		if err != nil {
			t.Fatalf("GetLocationBrokerCode: %v", err)
		}
		t.Cleanup(func() { _, _ = q.DeleteLocationBrokerCode(ctx, code.ID) })

		if code.LocationID != existing.ID {
			t.Errorf("LocationID = %d, want %d (should merge, not duplicate)", code.LocationID, existing.ID)
		}
		if code.Broker != db.BrokerAvance {
			t.Errorf("Broker = %q, want %q", code.Broker, db.BrokerAvance)
		}
	})

	t.Run("non-airport station gets its own location", func(t *testing.T) {
		loc, err := q.GetLocationByBrokerLocationID(ctx, "ZZZDT")
		if err != nil {
			t.Fatalf("GetLocationByBrokerLocationID: %v", err)
		}
		code, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker: db.BrokerAvance, BrokerLocationID: "ZZZDT", LocationID: loc.ID,
		})
		if err != nil {
			t.Fatalf("GetLocationBrokerCode: %v", err)
		}
		t.Cleanup(func() {
			_, _ = q.DeleteLocationBrokerCode(ctx, code.ID)
			_ = q.DeleteLocationByID(ctx, loc.ID)
		})

		if loc.ID == existing.ID {
			t.Error("downtown station reused the airport location")
		}
		if loc.Iata != nil {
			t.Errorf("Iata = %v, want nil", loc.Iata)
		}
	})
}

// TestInsertLocationsKeepsExistingLocations covers the import rules production depends on: an
// existing location is the source of truth, and a new airport is marked as one.
func TestInsertLocationsKeepsExistingLocations(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	ls := location.NewLocationService(q)

	importStations := func(t *testing.T, locs ...broker.Location) {
		t.Helper()
		b := &mockBroker{name: broker.BrokerAvance, pages: []broker.LocationPage{{Locations: locs}}}
		if err := ls.InsertLocations(ctx, b, ""); err != nil {
			t.Fatalf("InsertLocations: %v", err)
		}
	}
	locationOf := func(t *testing.T, code string) db.Location {
		t.Helper()
		loc, err := q.GetLocationByBrokerLocationID(ctx, code)
		if err != nil {
			t.Fatalf("GetLocationByBrokerLocationID(%s): %v", code, err)
		}
		return loc
	}
	deleteStation := func(code string) {
		if c, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{Broker: db.BrokerAvance, BrokerLocationID: code}); err == nil {
			_, _ = q.DeleteLocationBrokerCode(ctx, c.ID)
		}
	}
	cleanup := func(code string) {
		t.Cleanup(func() {
			loc, err := q.GetLocationByBrokerLocationID(ctx, code)
			if err != nil {
				return
			}
			deleteStation(code)
			_ = q.DeleteLocationByID(ctx, loc.ID)
		})
	}

	t.Run("an IATA match keeps the existing name and only fills blanks", func(t *testing.T) {
		existing, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Keep Test Airport (YYA)", Iata: strPtr("YYA")},
			db.BrokerFlex, "flex-keep-test",
		)
		importStations(t, broker.Location{ID: "KEEPAP", Name: "Supplier Name For YYA", Country: "Greece", CountryCode: "GR", City: "Keep City", Iata: "YYA"})
		t.Cleanup(func() { deleteStation("KEEPAP") })

		loc := locationOf(t, "KEEPAP")
		if loc.ID != existing.ID {
			t.Fatalf("station went to location %d, want the existing %d", loc.ID, existing.ID)
		}
		if loc.Name != "Keep Test Airport (YYA)" {
			t.Errorf("Name = %q, want the existing name kept", loc.Name)
		}
		if loc.City == nil || *loc.City != "Keep City" {
			t.Errorf("City = %v, want the empty city filled in", loc.City)
		}
	})

	t.Run("a new IATA location is an airport", func(t *testing.T) {
		cleanup("NEWAP")
		importStations(t, broker.Location{ID: "NEWAP", Name: "New Test Airport", Country: "Greece", CountryCode: "GR", Iata: "YYB"})

		loc := locationOf(t, "NEWAP")
		if loc.Iata == nil || *loc.Iata != "YYB" || !loc.IsAirport {
			t.Errorf("location = iata %v, airport %v, want YYB and an airport", loc.Iata, loc.IsAirport)
		}
	})
}
