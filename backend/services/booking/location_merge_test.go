package booking

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/location"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func mergeTestService() (*Service, *db.Queries, *pgxpool.Pool) {
	pool := sqldb.Driver[*pgxpool.Pool](bookingsDB)
	q := db.New(pool)
	return &Service{query: q, inTx: db.NewTxRunner(pool)}, q, pool
}

// seedMergeOffer inserts a price offer picking up at one location and dropping off at another.
func seedMergeOffer(t *testing.T, pool *pgxpool.Pool, pickupID, dropoffID int64) int64 {
	t.Helper()
	ctx := context.Background()
	var id int64
	err := pool.QueryRow(ctx, `
		INSERT INTO price_offers (
			agent_id, name, pickup_location_id, dropoff_location_id, pickup_date, dropoff_date,
			pickup_time, dropoff_time, driver_age, rental_days, plan_id, broker, rate_qualifier,
			supplier_code, car_details, currency_code, currency_rate, purchase_price,
			markup_percentage, broker_erp_price, bt_erp_price, total_price, offered_currency_code,
			offered_price
		) VALUES (
			1, 'Merge Test Offer', $1, $2, '2030-01-01', '2030-01-05', '10:00', '10:00', '30', 4,
			'plan', 'flex', 'rq', 'sup', '{}', 'EUR', 1, 100, 10, 0, 0, 110, 'EUR', 110
		) RETURNING id`, pickupID, dropoffID).Scan(&id)
	if err != nil {
		t.Fatalf("failed to seed price offer: %v", err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM price_offers WHERE id = $1`, id) })
	return id
}

func mergeSuffix() string {
	return fmt.Sprintf("%d %d", time.Now().UnixNano(), nextLocSeq())
}

func TestListLocationMergeSuggestions(t *testing.T) {
	ctx := context.Background()
	s, q, _ := mergeTestService()
	suffix := mergeSuffix()

	flex, _ := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Suggest Downtown " + suffix},
		db.BrokerFlex, "flex-suggest-"+suffix,
	)
	avance, _ := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "merge suggest - downtown " + suffix},
		db.BrokerAvance, "avance-suggest-"+suffix,
	)
	otherCountry, _ := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{Country: "Cyprus", CountryCode: "CY", Name: "Merge Suggest Downtown " + suffix},
		db.BrokerHertz, "hertz-suggest-"+suffix,
	)
	terminal1, _ := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Suggest Terminal 1 " + suffix},
		db.BrokerFlex, "flex-suggest-t1-"+suffix,
	)
	terminal2, _ := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Suggest Terminal 2 " + suffix},
		db.BrokerAvance, "avance-suggest-t2-"+suffix,
	)

	resp, err := s.ListLocationMergeSuggestions(ctx)
	if err != nil {
		t.Fatalf("ListLocationMergeSuggestions: %v", err)
	}

	groupOf := func(id int64) []int64 {
		for _, g := range resp.Groups {
			var ids []int64
			for _, l := range g.Locations {
				ids = append(ids, l.ID)
			}
			if slices.Contains(ids, id) {
				return ids
			}
		}
		return nil
	}

	t.Run("names that differ in punctuation and case are grouped", func(t *testing.T) {
		ids := groupOf(flex.ID)
		if !slices.Contains(ids, avance.ID) {
			t.Fatalf("group of %d = %v, want it to contain %d", flex.ID, ids, avance.ID)
		}
		if slices.Contains(ids, otherCountry.ID) {
			t.Errorf("group %v contains the same name from another country", ids)
		}
	})

	t.Run("digits are part of the name", func(t *testing.T) {
		if ids := groupOf(terminal1.ID); slices.Contains(ids, terminal2.ID) {
			t.Errorf("Terminal 1 and Terminal 2 were grouped: %v", ids)
		}
	})
}

func TestMergeLocations(t *testing.T) {
	ctx := context.Background()
	s, q, pool := mergeTestService()

	t.Run("moves codes, aliases and offers into the kept location", func(t *testing.T) {
		suffix := mergeSuffix()
		keepName := "Merge Downtown " + suffix
		removeName := "merge - downtown " + suffix
		keep, keepCode := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: keepName, City: strPtr("Athens")},
			db.BrokerFlex, "flex-merge-"+suffix,
		)
		remove, removeCode := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: removeName, Iata: strPtr("ZYX")},
			db.BrokerAvance, "avance-merge-"+suffix,
		)
		if err := q.InsertManyLocationAliases(ctx, db.InsertManyLocationAliasesParams{
			LocationIds: []int64{remove.ID}, Aliases: []string{"Old Alias " + suffix},
		}); err != nil {
			t.Fatalf("seed alias: %v", err)
		}
		offerID := seedMergeOffer(t, pool, remove.ID, remove.ID)

		resp, err := s.MergeLocations(ctx, location.MergeLocationsParams{
			KeepID: keep.ID, RemoveID: remove.ID,
			Name: keepName, Country: "Greece", CountryCode: "gr",
			City: "Athens", Iata: "zyx", IsAirport: false,
		})
		if err != nil {
			t.Fatalf("MergeLocations: %v", err)
		}

		merged := resp.Location
		if merged.ID != keep.ID || merged.Name != keepName || merged.CountryCode != "GR" {
			t.Errorf("merged = %+v, want id %d named %q in GR", merged, keep.ID, keepName)
		}
		if merged.Iata == nil || *merged.Iata != "ZYX" {
			t.Errorf("Iata = %v, want the removed location's ZYX", merged.Iata)
		}

		var codeIDs []int64
		for _, c := range merged.BrokerCodes {
			codeIDs = append(codeIDs, c.ID)
		}
		if !slices.Contains(codeIDs, keepCode.ID) || !slices.Contains(codeIDs, removeCode.ID) {
			t.Errorf("broker codes = %v, want both %d and %d", codeIDs, keepCode.ID, removeCode.ID)
		}

		if !slices.Contains(merged.Aliases, "Old Alias "+suffix) {
			t.Errorf("aliases = %v, want the removed location's alias", merged.Aliases)
		}
		if !slices.Contains(merged.Aliases, removeName) {
			t.Errorf("aliases = %v, want the discarded name %q", merged.Aliases, removeName)
		}

		if _, err := q.GetLocationById(ctx, remove.ID); !errors.Is(err, db.ErrNoRows) {
			t.Errorf("removed location still exists (err %v)", err)
		}

		var pickupID, dropoffID int64
		if err := pool.QueryRow(ctx, `SELECT pickup_location_id, dropoff_location_id FROM price_offers WHERE id = $1`, offerID).
			Scan(&pickupID, &dropoffID); err != nil {
			t.Fatalf("read offer: %v", err)
		}
		if pickupID != keep.ID || dropoffID != keep.ID {
			t.Errorf("offer locations = %d/%d, want %d", pickupID, dropoffID, keep.ID)
		}
	})

	t.Run("the kept location can take the removed location's name", func(t *testing.T) {
		suffix := mergeSuffix()
		keepName := "Merge Rename " + suffix
		removeName := "Merge - Rename " + suffix
		keep, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: keepName},
			db.BrokerFlex, "flex-rename-"+suffix,
		)
		remove, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: removeName},
			db.BrokerAvance, "avance-rename-"+suffix,
		)

		resp, err := s.MergeLocations(ctx, location.MergeLocationsParams{
			KeepID: keep.ID, RemoveID: remove.ID, Name: removeName, Country: "Greece", CountryCode: "GR",
		})
		if err != nil {
			t.Fatalf("MergeLocations: %v", err)
		}
		if resp.Location.Name != removeName {
			t.Errorf("Name = %q, want %q", resp.Location.Name, removeName)
		}
		if !slices.Contains(resp.Location.Aliases, keepName) {
			t.Errorf("aliases = %v, want the discarded name %q", resp.Location.Aliases, keepName)
		}
	})

	t.Run("locations with a code of the same broker are rejected", func(t *testing.T) {
		suffix := mergeSuffix()
		keep, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Shared A " + suffix},
			db.BrokerFlex, "flex-shared-a-"+suffix,
		)
		remove, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Shared B " + suffix},
			db.BrokerFlex, "flex-shared-b-"+suffix,
		)

		_, err := s.MergeLocations(ctx, location.MergeLocationsParams{
			KeepID: keep.ID, RemoveID: remove.ID, Name: "Merge Shared A " + suffix, Country: "Greece", CountryCode: "GR",
		})
		api_errors.AssertApiError(t, api_errors.NewErrorWithDetail(errs.FailedPrecondition, "both locations have a code of the same broker", api_errors.ErrorDetails{
			Code: api_errors.CodeLocationsShareBroker,
		}), err)

		if _, err := q.GetLocationById(ctx, remove.ID); err != nil {
			t.Errorf("rejected merge deleted the location: %v", err)
		}
	})

	t.Run("a name taken by a third location is rejected and nothing changes", func(t *testing.T) {
		suffix := mergeSuffix()
		taken := "Merge Taken " + suffix
		seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: taken},
			db.BrokerHertz, "hertz-taken-"+suffix,
		)
		keep, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Conflict A " + suffix},
			db.BrokerFlex, "flex-conflict-"+suffix,
		)
		remove, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Conflict B " + suffix},
			db.BrokerAvance, "avance-conflict-"+suffix,
		)

		_, err := s.MergeLocations(ctx, location.MergeLocationsParams{
			KeepID: keep.ID, RemoveID: remove.ID, Name: taken, Country: "Greece", CountryCode: "GR",
		})
		api_errors.AssertApiError(t, api_errors.NewErrorWithDetail(errs.FailedPrecondition, "another location already has this name or iata", api_errors.ErrorDetails{
			Code: api_errors.CodeLocationMergeConflict,
		}), err)

		if _, err := q.GetLocationById(ctx, remove.ID); err != nil {
			t.Errorf("rolled back merge deleted the location: %v", err)
		}
	})

	t.Run("a missing location is not found", func(t *testing.T) {
		suffix := mergeSuffix()
		keep, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: "Merge Missing " + suffix},
			db.BrokerFlex, "flex-missing-"+suffix,
		)

		_, err := s.MergeLocations(ctx, location.MergeLocationsParams{
			KeepID: keep.ID, RemoveID: -1, Name: "Merge Missing " + suffix, Country: "Greece", CountryCode: "GR",
		})
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("a re-import does not recreate the removed location", func(t *testing.T) {
		suffix := mergeSuffix()
		keepName := "Merge Reimport " + suffix
		removeName := "Merge - Reimport " + suffix
		keep, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: keepName},
			db.BrokerFlex, "flex-reimport-"+suffix,
		)
		remove, _ := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{Country: "Greece", CountryCode: "GR", Name: removeName},
			db.BrokerAvance, "avance-reimport-"+suffix,
		)

		if _, err := s.MergeLocations(ctx, location.MergeLocationsParams{
			KeepID: keep.ID, RemoveID: remove.ID, Name: keepName, Country: "Greece", CountryCode: "GR",
		}); err != nil {
			t.Fatalf("MergeLocations: %v", err)
		}

		b := &mockBroker{name: broker.BrokerAvance, pages: []broker.LocationPage{{Locations: []broker.Location{
			{ID: "avance-reimport-" + suffix, Name: removeName, Country: "Greece", CountryCode: "GR"},
		}}}}
		if err := location.NewLocationService(q).InsertLocations(ctx, b, ""); err != nil {
			t.Fatalf("InsertLocations: %v", err)
		}

		var count int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM locations WHERE country_code = 'GR' AND lower(name) = lower($1)`, removeName).
			Scan(&count); err != nil {
			t.Fatalf("count locations: %v", err)
		}
		if count != 0 {
			t.Errorf("re-import recreated %q", removeName)
		}
	})
}

func TestMergeLocationsValidation(t *testing.T) {
	valid := location.MergeLocationsParams{KeepID: 1, RemoveID: 2, Name: "A", Country: "Greece", CountryCode: "GR"}

	t.Run("accepts a valid merge", func(t *testing.T) {
		if err := valid.Validate(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("rejects merging a location into itself", func(t *testing.T) {
		p := valid
		p.RemoveID = p.KeepID
		api_errors.AssertApiError(t, invalidValueErr("keep_id"), p.Validate())
	})

	t.Run("rejects a blank name", func(t *testing.T) {
		p := valid
		p.Name = "  "
		api_errors.AssertApiError(t, invalidValueErr("name"), p.Validate())
	})

	t.Run("rejects an iata that is not three letters", func(t *testing.T) {
		p := valid
		p.Iata = "AB"
		api_errors.AssertApiError(t, invalidValueErr("iata"), p.Validate())
	})
}
