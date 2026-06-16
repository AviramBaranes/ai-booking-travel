package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/location"
)

func TestDeleteLocation(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes broker code and keeps location when other broker codes exist", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		cityTLV := "Tel Aviv"
		iataTLV := "TLV"
		loc, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", City: &cityTLV, Name: "Ben Gurion Airport", Iata: &iataTLV,
			},
			db.BrokerFlex, "flex-tlv-del-keep",
		)

		// Add a second broker code to the same location
		_, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
			LocationID:       loc.ID,
			Broker:           db.BrokerHertz,
			BrokerLocationID: "hertz-tlv-del-keep",
		})
		if err != nil {
			t.Fatalf("failed to insert second broker code: %v", err)
		}

		// Delete the first broker code
		err = s.DeleteLocation(ctx, lbc.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Location should still exist
		_, err = q.GetLocationById(ctx, loc.ID)
		if err != nil {
			t.Fatalf("location should still exist, got error: %v", err)
		}
	})

	t.Run("deletes broker code and location when no other broker codes exist", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		cityParis := "Paris"
		iataCDG := "CDG"
		loc, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "France", CountryCode: "FR", City: &cityParis, Name: "Charles de Gaulle", Iata: &iataCDG,
			},
			db.BrokerFlex, "flex-cdg-del-orphan",
		)

		err := s.DeleteLocation(ctx, lbc.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Location should be deleted
		_, err = q.GetLocationById(ctx, loc.ID)
		if !errors.Is(err, db.ErrNoRows) {
			t.Fatalf("expected location to be deleted, got err: %v", err)
		}
	})

	t.Run("returns not found for non-existent broker code", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		err := s.DeleteLocation(ctx, 999999)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

}

func TestToggleLocation(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Helper to read back the enabled state of a broker code.
	getEnabled := func(t *testing.T, lbc db.LocationBrokerCode) bool {
		t.Helper()
		row, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker:           lbc.Broker,
			BrokerLocationID: lbc.BrokerLocationID,
			LocationID:       lbc.LocationID,
		})
		if err != nil {
			t.Fatalf("failed to fetch broker code: %v", err)
		}
		return row.Enabled
	}

	t.Run("disables an enabled location", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Disable Test",
			},
			db.BrokerFlex, "flex-toggle-disable",
		)

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to be enabled by default")
		}

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: false})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if getEnabled(t, lbc) {
			t.Fatal("expected location to be disabled after toggle")
		}
	})

	t.Run("enables a disabled location", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Enable Test",
			},
			db.BrokerFlex, "flex-toggle-enable",
		)

		// Disable first
		if err := q.ToggleLocationBrokerCode(ctx, db.ToggleLocationBrokerCodeParams{
			Enabled: false,
			ID:      lbc.ID,
		}); err != nil {
			t.Fatalf("failed to disable: %v", err)
		}
		if getEnabled(t, lbc) {
			t.Fatal("expected location to be disabled")
		}

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: true})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to be enabled after toggle")
		}
	})

	t.Run("enabling an already enabled location is idempotent", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Idempotent Test",
			},
			db.BrokerFlex, "flex-toggle-idempotent",
		)

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to be enabled by default")
		}

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: true})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to still be enabled")
		}
	})

	t.Run("disabling an already disabled location is idempotent", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Idempotent Disable",
			},
			db.BrokerFlex, "flex-toggle-idempotent-dis",
		)

		if err := q.ToggleLocationBrokerCode(ctx, db.ToggleLocationBrokerCodeParams{
			Enabled: false,
			ID:      lbc.ID,
		}); err != nil {
			t.Fatalf("failed to disable: %v", err)
		}

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: false})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if getEnabled(t, lbc) {
			t.Fatal("expected location to still be disabled")
		}
	})

}

func TestToggleLocationIsAirport(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Helper to read back the is_airport flag for a location.
	getIsAirport := func(t *testing.T, locID int64) bool {
		t.Helper()
		row, err := q.GetLocationById(ctx, locID)
		if err != nil {
			t.Fatalf("failed to fetch location: %v", err)
		}
		return row.IsAirport
	}

	t.Run("sets is_airport to true", func(t *testing.T) {
		loc, locBc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Airport Toggle True Test",
			},
			db.BrokerFlex, "flex-airport-toggle-true",
		)

		if getIsAirport(t, loc.ID) {
			t.Fatal("expected is_airport to be false by default")
		}

		err := s.ToggleLocationIsAirport(ctx, locBc.ID, location.ToggleLocationIsAirportParams{IsAirport: true})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !getIsAirport(t, loc.ID) {
			t.Fatal("expected is_airport to be true after toggle")
		}
	})

	t.Run("sets is_airport to false", func(t *testing.T) {
		loc, locBc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Airport Toggle False Test",
			},
			db.BrokerFlex, "flex-airport-toggle-false",
		)

		// Set to true first
		if err := q.ToggleIsAirport(ctx, db.ToggleIsAirportParams{
			IsAirport: true,
			ID:        loc.ID,
		}); err != nil {
			t.Fatalf("failed to set is_airport: %v", err)
		}
		if !getIsAirport(t, loc.ID) {
			t.Fatal("expected is_airport to be true")
		}

		err := s.ToggleLocationIsAirport(ctx, locBc.ID, location.ToggleLocationIsAirportParams{IsAirport: false})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if getIsAirport(t, loc.ID) {
			t.Fatal("expected is_airport to be false after toggle")
		}
	})

	t.Run("setting true when already true is idempotent", func(t *testing.T) {
		loc, locBc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Airport Idempotent True Test",
			},
			db.BrokerFlex, "flex-airport-idempotent-true",
		)

		if err := q.ToggleIsAirport(ctx, db.ToggleIsAirportParams{
			IsAirport: true,
			ID:        loc.ID,
		}); err != nil {
			t.Fatalf("failed to set is_airport: %v", err)
		}

		err := s.ToggleLocationIsAirport(ctx, locBc.ID, location.ToggleLocationIsAirportParams{IsAirport: true})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !getIsAirport(t, loc.ID) {
			t.Fatal("expected is_airport to still be true")
		}
	})

	t.Run("setting false when already false is idempotent", func(t *testing.T) {
		loc, locBc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Airport Idempotent False Test",
			},
			db.BrokerFlex, "flex-airport-idempotent-false",
		)

		if getIsAirport(t, loc.ID) {
			t.Fatal("expected is_airport to be false by default")
		}

		err := s.ToggleLocationIsAirport(ctx, locBc.ID, location.ToggleLocationIsAirportParams{IsAirport: false})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if getIsAirport(t, loc.ID) {
			t.Fatal("expected is_airport to still be false")
		}
	})
}
