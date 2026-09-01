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
