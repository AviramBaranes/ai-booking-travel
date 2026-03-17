package booking

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	locations_mocks "encore.app/services/booking/mocks"
	"go.uber.org/mock/gomock"
)

func TestSearchLocations(t *testing.T) {
	ctx := context.Background()

	t.Run("validate the search query is not missing or too short", func(t *testing.T) {
		cases := []struct {
			name            string
			params          SearchLocationParams
			isExpectedError bool
		}{
			{name: "missing search query", params: SearchLocationParams{}, isExpectedError: true},
			{name: "search query too short", params: SearchLocationParams{Search: "ab"}, isExpectedError: true},
			{name: "valid search query", params: SearchLocationParams{Search: "abc"}, isExpectedError: false},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.params.Validate()
				if tc.isExpectedError && err == nil {
					t.Fatalf("expected validation error, got nil")
				}
				if !tc.isExpectedError && err != nil {
					t.Fatalf("expected no validation error, got %v", err)
				}
			})
		}

	})

	t.Run("returns empty list when no locations match the search query", func(t *testing.T) {
		res, err := SearchLocations(ctx, SearchLocationParams{Search: "nonexistent"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(res.Locations) != 0 {
			t.Fatalf("expected empty list of locations, got %v", res.Locations)
		}
	})

	t.Run("returns error when database query fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := locations_mocks.NewMockQuerier(ctrl)
		q.EXPECT().SearchLocations(gomock.Any(), "error").Return(nil, errors.New("database error")).Times(1)

		s := &Service{query: q}
		_, err := s.SearchLocations(ctx, SearchLocationParams{Search: "error"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns list of locations that match the search query", func(t *testing.T) {
		search := "TES"
		ids, err := query.InsertManyLocation(ctx, db.InsertManyLocationParams{
			Countries:    []string{search, "Thailand", "Israel", "USA", "Germany"},
			Cities:       []string{"London", search, "Tel Aviv", "New York", "Berlin"},
			Names:        []string{"London Airport", "Bangkok Airport", search, "NY Airport", "Berlin Airport"},
			CountryCodes: []string{"EN", "TH", "IL", "US", "DE"},
			Iatas:        []string{"LON", "BKK", "TLV", "NYC", search},
		})
		if err != nil {
			t.Fatalf("failed to insert test locations: %v", err)
		}

		for i, id := range ids {
			_, err = query.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
				LocationID:       id,
				Broker:           db.BrokerFlex,
				BrokerLocationID: fmt.Sprintf("loc-%d", i),
			})
			if err != nil {
				t.Fatalf("failed to insert broker code for location %d: %v", id, err)
			}

			if i == 1 {
				err = query.DisableLocationBrokerCode(ctx, id)
				if err != nil {
					t.Fatalf("failed to disable broker code for location %d: %v", id, err)
				}
			}
		}

		// changing search case to prove search is case-insensitive
		search = "tEs"
		locs, err := SearchLocations(ctx, SearchLocationParams{Search: search})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(locs.Locations) != 4 {
			t.Fatalf("expected 4 matching locations, got %d", len(locs.Locations))
		}

		// ordered by IATA ASC: BKK, LON, TES, TLV
		expectedLocs := []struct {
			id   int64
			name string
			iata string
		}{
			{id: ids[1], name: "Bangkok Airport", iata: "BKK"},
			{id: ids[0], name: "London Airport", iata: "LON"},
			{id: ids[4], name: "Berlin Airport", iata: "TES"},
			{id: ids[2], name: "TES", iata: "TLV"},
		}

		for i, expected := range expectedLocs {
			loc := locs.Locations[i]
			if loc.ID != expected.id {
				t.Errorf("expected location [%d] to have ID %d, got %d", i, expected.id, loc.ID)
			}
			if loc.Name != expected.name {
				t.Errorf("expected location [%d] to have Name '%s', got '%s'", i, expected.name, loc.Name)
			}
			if loc.Iata == nil || *loc.Iata != expected.iata {
				t.Errorf("expected location [%d] to have Iata '%s', got '%v'", i, expected.iata, loc.Iata)
			}
		}
	})
}
