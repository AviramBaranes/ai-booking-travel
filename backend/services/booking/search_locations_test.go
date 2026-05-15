package booking

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

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
		// Use a per-run token to avoid matching seeded/shared test data.
		search := strings.ToUpper(strconv.FormatInt(time.Now().UnixNano(), 36))
		if len(search) > 3 {
			search = search[len(search)-3:]
		}
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

		t.Cleanup(func() {
			for _, id := range ids {
				_ = query.DeleteLocationByID(ctx, id)
			}
		})

		for i, id := range ids {
			lbc, err := query.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
				LocationID:       id,
				Broker:           db.BrokerFlex,
				BrokerLocationID: fmt.Sprintf("loc-%d", i),
			})
			if err != nil {
				t.Fatalf("failed to insert broker code for location %d: %v", id, err)
			}

			if i == 3 {
				err = query.DisableLocationBrokerCode(ctx, lbc.ID)
				if err != nil {
					t.Fatalf("failed to disable broker code for location %d: %v", id, err)
				}
			}
		}

		searchToken := search
		// Change search case to prove search is case-insensitive.
		search = strings.ToLower(search)
		locs, err := SearchLocations(ctx, SearchLocationParams{Search: search})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expectedByID := map[int64]struct {
			name string
			iata string
		}{
			ids[1]: {name: "Bangkok Airport", iata: "BKK"},
			ids[0]: {name: "London Airport", iata: "LON"},
			ids[4]: {name: "Berlin Airport", iata: searchToken},
			ids[2]: {name: searchToken, iata: "TLV"},
		}

		matched := map[int64]LocationResult{}
		for _, loc := range locs.Locations {
			if _, ok := expectedByID[loc.ID]; ok {
				matched[loc.ID] = loc
			}
		}

		if len(matched) != len(expectedByID) {
			t.Fatalf("expected %d seeded locations to match, got %d", len(expectedByID), len(matched))
		}

		for id, expected := range expectedByID {
			loc := matched[id]
			if loc.Name != expected.name {
				t.Errorf("expected location ID %d to have Name '%s', got '%s'", id, expected.name, loc.Name)
			}
			if loc.Iata == nil || *loc.Iata != expected.iata {
				t.Errorf("expected location ID %d to have Iata '%s', got '%v'", id, expected.iata, loc.Iata)
			}
		}
	})
}
