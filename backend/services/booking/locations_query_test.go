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
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/location"
	locations_mocks "encore.app/services/booking/mocks"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func listLocationsInvalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

// seedLocationWithBrokerCode inserts a location and a broker code pointing to it.
// It registers a t.Cleanup to delete the broker code and location when the test ends.
func seedLocationWithBrokerCode(t *testing.T, q *db.Queries, loc db.InsertLocationParams, broker db.Broker, brokerLocID string) (db.Location, db.LocationBrokerCode) {
	t.Helper()
	ctx := context.Background()

	l, err := q.InsertLocation(ctx, loc)
	if err != nil {
		t.Fatalf("failed to seed location: %v", err)
	}

	lbc, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
		LocationID:       l.ID,
		Broker:           broker,
		BrokerLocationID: brokerLocID,
	})
	if err != nil {
		t.Fatalf("failed to seed broker code: %v", err)
	}

	t.Cleanup(func() {
		// Ignore errors — row may already be deleted by the test itself.
		_, _ = q.DeleteLocationBrokerCode(ctx, lbc.ID)
		_ = q.DeleteLocationByID(ctx, l.ID)
	})

	return l, lbc
}

// --- Tests ---

func TestListLocationsValidation(t *testing.T) {
	t.Run("rejects page 0", func(t *testing.T) {
		p := location.ListLocationsParams{Page: 0}
		api_errors.AssertApiError(t, listLocationsInvalidValueErr("page"), p.Validate())
	})

	t.Run("rejects negative page", func(t *testing.T) {
		p := location.ListLocationsParams{Page: -1}
		api_errors.AssertApiError(t, listLocationsInvalidValueErr("page"), p.Validate())
	})

	t.Run("accepts valid params", func(t *testing.T) {
		p := location.ListLocationsParams{Name: "test", Page: 1}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("accepts no filters", func(t *testing.T) {
		p := location.ListLocationsParams{Page: 1}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestListLocations(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Seed two distinct locations in different countries
	cityTLV := "Tel Aviv"
	iataTLV := "TLV"
	_, lbcIL := seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{
			Country: "Israel", CountryCode: "IL", City: &cityTLV, Name: "Ben Gurion Airport", Iata: &iataTLV,
		},
		db.BrokerFlex, "flex-tlv-list",
	)

	cityLHR := "London"
	iataLHR := "LHR"
	seedLocationWithBrokerCode(t, q,
		db.InsertLocationParams{
			Country: "United Kingdom", CountryCode: "GB", City: &cityLHR, Name: "Heathrow Airport", Iata: &iataLHR,
		},
		db.BrokerFlex, "flex-lhr-list",
	)

	t.Run("returns seeded location with correct fields", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Name: "Ben Gurion", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found := false
		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-tlv-list" {
				found = true
				if r.ID != lbcIL.ID {
					t.Fatalf("expected broker code ID %d, got %d", lbcIL.ID, r.ID)
				}
				if r.Name != "Ben Gurion Airport" {
					t.Fatalf("expected name 'Ben Gurion Airport', got %q", r.Name)
				}
				if r.CountryCode != "IL" {
					t.Fatalf("expected country_code 'IL', got %q", r.CountryCode)
				}
				if r.Country != "Israel" {
					t.Fatalf("expected country 'Israel', got %q", r.Country)
				}
				if r.City == nil || *r.City != "Tel Aviv" {
					t.Fatalf("expected city 'Tel Aviv', got %v", r.City)
				}
				if r.Iata == nil || *r.Iata != "TLV" {
					t.Fatalf("expected iata 'TLV', got %v", r.Iata)
				}
				break
			}
		}
		if !found {
			t.Fatal("seeded location not found in results")
		}
	})

	t.Run("filter by country code returns only matching locations", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{CountryCode: "IL", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-lhr-list" {
				t.Fatal("GB location should not appear when filtering by IL")
			}
		}

		found := false
		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-tlv-list" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected IL location in results")
		}
	})

	t.Run("filter by IATA returns only matching location", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Iata: "LHR", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-tlv-list" {
				t.Fatal("TLV location should not appear when filtering by LHR")
			}
		}

		found := false
		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-lhr-list" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected LHR location in results")
		}
	})

	t.Run("filter by name returns only matching location", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Name: "Heathrow", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-tlv-list" {
				t.Fatal("TLV location should not appear when filtering by Heathrow")
			}
		}

		found := false
		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-lhr-list" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected Heathrow location in results")
		}
	})

	t.Run("filter by enabled returns only matching locations", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Enabled: "true", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, r := range resp.Locations {
			if !r.Enabled {
				t.Fatal("disabled location should not appear when filtering by enabled=true")
			}
		}
	})

	t.Run("no filters returns all seeded locations", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		foundTLV, foundLHR := false, false
		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-tlv-list" {
				foundTLV = true
			}
			if r.BrokerLocationID == "flex-lhr-list" {
				foundLHR = true
			}
		}
		if !foundTLV || !foundLHR {
			t.Fatal("expected both locations when no filters applied")
		}
	})

	t.Run("total reflects filtered count", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{CountryCode: "IL", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != int64(len(resp.Locations)) {
			t.Fatalf("expected total %d to match location count %d", resp.Total, len(resp.Locations))
		}
	})

	t.Run("non-matching filter returns empty", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Name: "zzzznonexistent", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Locations) != 0 {
			t.Fatalf("expected 0 locations, got %d", len(resp.Locations))
		}
		if resp.Total != 0 {
			t.Fatalf("expected total 0, got %d", resp.Total)
		}
	})

	t.Run("filter by broker returns only matching broker codes", func(t *testing.T) {
		// Add a hertz broker code to the IL location
		_, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
			LocationID:       lbcIL.LocationID,
			Broker:           db.BrokerHertz,
			BrokerLocationID: "hertz-tlv-list",
		})
		if err != nil {
			t.Fatalf("failed to insert second broker code: %v", err)
		}

		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Broker: "hertz", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-tlv-list" || r.BrokerLocationID == "flex-lhr-list" {
				t.Fatal("flex broker codes should not appear when filtering by hertz")
			}
		}

		found := false
		for _, r := range resp.Locations {
			if r.BrokerLocationID == "hertz-tlv-list" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected hertz broker code in results")
		}
	})

	t.Run("multiple filters combine with AND", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{
			CountryCode: "IL",
			Name:        "Ben Gurion",
			Page:        1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		for _, r := range resp.Locations {
			if r.BrokerLocationID == "flex-lhr-list" {
				t.Fatal("GB location should not appear with IL+Ben Gurion filters")
			}
		}

		found := false
		for _, r := range resp.Locations {
			if r.Name == "Ben Gurion Airport" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected Ben Gurion location with combined filters")
		}
	})

	t.Run("returns error when count query fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().CountLocationBrokerCodesWithLocation(gomock.Any(), gomock.Any()).
			Return(int64(0), errors.New("db error"))

		_, err := s.ListLocations(ctx, location.ListLocationsParams{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when list query fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().CountLocationBrokerCodesWithLocation(gomock.Any(), gomock.Any()).
			Return(int64(5), nil)
		q.EXPECT().ListLocationBrokerCodesWithLocation(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("db error"))

		_, err := s.ListLocations(ctx, location.ListLocationsParams{Page: 1})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestListLocationsPagination(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Seed 16 locations with a shared prefix so we can filter for them specifically.
	// LocationsLimit is 15, so page 1 should have 15, page 2 should have 1.
	prefix := "PagTest"
	for i := 1; i <= 16; i++ {
		name := fmt.Sprintf("%s Location %02d", prefix, i)
		cc := fmt.Sprintf("P%d", i)
		brokerID := fmt.Sprintf("flex-pag-%02d", i)
		seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Pagland", CountryCode: cc, Name: name,
			},
			db.BrokerFlex, brokerID,
		)
	}

	t.Run("page 1 returns exactly 15 results", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Name: prefix, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Locations) != location.LocationsLimit {
			t.Fatalf("expected %d locations on page 1, got %d", location.LocationsLimit, len(resp.Locations))
		}
	})

	t.Run("total reflects all matching rows across pages", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Name: prefix, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Total != 16 {
			t.Fatalf("expected total 16, got %d", resp.Total)
		}
	})

	t.Run("page 2 returns the remaining result", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Name: prefix, Page: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Locations) != 1 {
			t.Fatalf("expected 1 location on page 2, got %d", len(resp.Locations))
		}
	})

	t.Run("page 2 does not repeat page 1 results", func(t *testing.T) {
		page1, err := s.ListLocations(ctx, location.ListLocationsParams{Name: prefix, Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		page2, err := s.ListLocations(ctx, location.ListLocationsParams{Name: prefix, Page: 2})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		page1IDs := map[int64]bool{}
		for _, r := range page1.Locations {
			page1IDs[r.ID] = true
		}
		for _, r := range page2.Locations {
			if page1IDs[r.ID] {
				t.Fatalf("location ID %d appeared on both page 1 and page 2", r.ID)
			}
		}
	})

	t.Run("page 3 returns empty", func(t *testing.T) {
		resp, err := s.ListLocations(ctx, location.ListLocationsParams{Name: prefix, Page: 3})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Locations) != 0 {
			t.Fatalf("expected 0 locations on page 3, got %d", len(resp.Locations))
		}
	})
}

func TestSearchLocations(t *testing.T) {
	ctx := context.Background()

	t.Run("validate the search query is not missing or too short", func(t *testing.T) {
		cases := []struct {
			name            string
			params          location.SearchLocationParams
			isExpectedError bool
		}{
			{name: "missing search query", params: location.SearchLocationParams{}, isExpectedError: true},
			{name: "search query too short", params: location.SearchLocationParams{Search: "ab"}, isExpectedError: true},
			{name: "valid search query", params: location.SearchLocationParams{Search: "abc"}, isExpectedError: false},
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
		res, err := SearchLocations(ctx, location.SearchLocationParams{Search: "nonexistent"})
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
		_, err := s.SearchLocations(ctx, location.SearchLocationParams{Search: "error"})
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
		locs, err := SearchLocations(ctx, location.SearchLocationParams{Search: search})
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

		matched := map[int64]location.LocationResult{}
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
