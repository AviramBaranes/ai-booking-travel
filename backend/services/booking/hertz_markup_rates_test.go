package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	locations_mocks "encore.app/services/booking/mocks"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func validCreateParams() CreateHertzMarkupRateRequest {
	return CreateHertzMarkupRateRequest{
		Country: "US", Brand: "ZR", CarGroup: "E",
		PickupDateFrom: "2026-01-01", PickupDateTo: "2026-12-31",
		NumOfRentalDaysFrom: 1, NumOfRentalDaysTo: 7,
		MarkUpGross: 15.5, MarkUpNet: 10.0,
	}
}

func validUpdateParams() UpdateHertzMarkupRateRequest {
	return UpdateHertzMarkupRateRequest{
		Country: "US", Brand: "ZR", CarGroup: "E",
		PickupDateFrom: "2026-01-01", PickupDateTo: "2026-12-31",
		NumOfRentalDaysFrom: 1, NumOfRentalDaysTo: 7,
		MarkUpGross: 20.0, MarkUpNet: 12.0,
	}
}

func mockService(t *testing.T) (*locations_mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := locations_mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func invalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

// --- Tests grouped by endpoint ---

func TestListHertzMarkupRates(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing page", func(t *testing.T) {
		api_errors.AssertApiError(t, invalidValueErr("page"), (ListHertzMarkupRatesRequest{Limit: 10}).Validate())
	})

	t.Run("validation rejects limit exceeds max", func(t *testing.T) {
		api_errors.AssertApiError(t, invalidValueErr("limit"), (ListHertzMarkupRatesRequest{Page: 1, Limit: 101}).Validate())
	})

	t.Run("validation rejects invalid sort direction", func(t *testing.T) {
		api_errors.AssertApiError(t, invalidValueErr("sortDir"), (ListHertzMarkupRatesRequest{Page: 1, Limit: 10, SortDir: "up"}).Validate())
	})

	t.Run("validation rejects invalid sort field", func(t *testing.T) {
		api_errors.AssertApiError(t, api_errors.ErrInvalidValue, (ListHertzMarkupRatesRequest{Page: 1, Limit: 10, SortBy: "invalid_field"}).Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := (ListHertzMarkupRatesRequest{Country: "US", SortBy: "country", SortDir: "desc", Page: 1, Limit: 50}).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns rates successfully", func(t *testing.T) {
		created, err := s.CreateHertzMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}

		resp, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "US", Brand: "ZR", Page: 1, Limit: 100})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found := false
		for _, r := range resp.Rates {
			if r.ID == created.ID {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("created rate not found in list")
		}
	})

	t.Run("returns empty list when no rates match filters", func(t *testing.T) {
		resp, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "XX", Page: 1, Limit: 10})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Rates) != 0 {
			t.Fatalf("expected 0 rates, got %d", len(resp.Rates))
		}
	})

	t.Run("filters by country", func(t *testing.T) {
		// Seed 3 DE rates and 2 IT rates under unique brand to isolate
		for i := 0; i < 3; i++ {
			p := validCreateParams()
			p.Country = "DE"
			p.Brand = "FC"
			p.NumOfRentalDaysFrom = int32(i + 1)
			p.NumOfRentalDaysTo = int32(i + 10)
			if _, err := s.CreateHertzMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed DE: %v", err)
			}
		}
		for i := 0; i < 2; i++ {
			p := validCreateParams()
			p.Country = "IT"
			p.Brand = "FC"
			p.NumOfRentalDaysFrom = int32(i + 1)
			p.NumOfRentalDaysTo = int32(i + 10)
			if _, err := s.CreateHertzMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed IT: %v", err)
			}
		}

		resp, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "DE", Brand: "FC", Page: 1, Limit: 100})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Rates) != 3 {
			t.Fatalf("expected 3 DE rates, got %d", len(resp.Rates))
		}
		for _, r := range resp.Rates {
			if r.Country != "DE" {
				t.Fatalf("expected country DE, got %s", r.Country)
			}
		}
	})

	t.Run("filters by brand and car group", func(t *testing.T) {
		// Seed 2 with brand=FX carGroup=F, and 2 with brand=FX carGroup=G
		for i := 0; i < 2; i++ {
			p := validCreateParams()
			p.Country = "BF"
			p.Brand = "FX"
			p.CarGroup = "F"
			p.NumOfRentalDaysFrom = int32(i + 1)
			p.NumOfRentalDaysTo = int32(i + 10)
			if _, err := s.CreateHertzMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed FX/F: %v", err)
			}
		}
		for i := 0; i < 2; i++ {
			p := validCreateParams()
			p.Country = "BF"
			p.Brand = "FX"
			p.CarGroup = "G"
			p.NumOfRentalDaysFrom = int32(i + 1)
			p.NumOfRentalDaysTo = int32(i + 10)
			if _, err := s.CreateHertzMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed FX/G: %v", err)
			}
		}

		resp, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "BF", Brand: "FX", CarGroup: "F", Page: 1, Limit: 100})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Rates) != 2 {
			t.Fatalf("expected 2 rates with brand=FX carGroup=F, got %d", len(resp.Rates))
		}
		for _, r := range resp.Rates {
			if r.Brand != "FX" || r.CarGroup != "F" {
				t.Fatalf("expected brand=FX carGroup=F, got brand=%s carGroup=%s", r.Brand, r.CarGroup)
			}
		}
	})

	t.Run("paginates results", func(t *testing.T) {
		// Seed 5 rates under unique filter to test pagination across 3 pages
		for i := 0; i < 5; i++ {
			p := validCreateParams()
			p.Country = "PG"
			p.Brand = "TT"
			p.NumOfRentalDaysFrom = int32(i + 1)
			p.NumOfRentalDaysTo = int32(i + 10)
			if _, err := s.CreateHertzMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed: %v", err)
			}
		}

		page1, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "PG", Brand: "TT", Page: 1, Limit: 2})
		if err != nil {
			t.Fatalf("page 1 error: %v", err)
		}
		if len(page1.Rates) != 2 {
			t.Fatalf("expected 2 rates on page 1, got %d", len(page1.Rates))
		}

		page2, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "PG", Brand: "TT", Page: 2, Limit: 2})
		if err != nil {
			t.Fatalf("page 2 error: %v", err)
		}
		if len(page2.Rates) != 2 {
			t.Fatalf("expected 2 rates on page 2, got %d", len(page2.Rates))
		}

		page3, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "PG", Brand: "TT", Page: 3, Limit: 2})
		if err != nil {
			t.Fatalf("page 3 error: %v", err)
		}
		if len(page3.Rates) != 1 {
			t.Fatalf("expected 1 rate on page 3, got %d", len(page3.Rates))
		}

		// Verify no overlap between pages
		seen := make(map[int64]bool)
		for _, r := range page1.Rates {
			seen[r.ID] = true
		}
		for _, r := range page2.Rates {
			if seen[r.ID] {
				t.Fatalf("page 2 contains rate %d already seen on page 1", r.ID)
			}
			seen[r.ID] = true
		}
		for _, r := range page3.Rates {
			if seen[r.ID] {
				t.Fatalf("page 3 contains rate %d already seen on earlier page", r.ID)
			}
		}
	})

	t.Run("sorts by car_group descending", func(t *testing.T) {
		for _, cg := range []string{"A", "C", "B"} {
			p := validCreateParams()
			p.Country = "SR"
			p.CarGroup = cg
			if _, err := s.CreateHertzMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed: %v", err)
			}
		}

		resp, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Country: "SR", SortBy: "car_group", SortDir: "desc", Page: 1, Limit: 100})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for i := 1; i < len(resp.Rates); i++ {
			if resp.Rates[i].CarGroup > resp.Rates[i-1].CarGroup {
				t.Fatalf("expected descending order, got %s after %s", resp.Rates[i].CarGroup, resp.Rates[i-1].CarGroup)
			}
		}
	})

	t.Run("defaults sort to country ascending", func(t *testing.T) {
		for _, c := range []string{"FR", "CA", "GB"} {
			p := validCreateParams()
			p.Country = c
			p.Brand = "DF"
			if _, err := s.CreateHertzMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed: %v", err)
			}
		}

		resp, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Brand: "DF", Page: 1, Limit: 100})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for i := 1; i < len(resp.Rates); i++ {
			if resp.Rates[i].Country < resp.Rates[i-1].Country {
				t.Fatalf("expected ascending country order, got %s after %s", resp.Rates[i].Country, resp.Rates[i-1].Country)
			}
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().ListHertzMarkupRates(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListHertzMarkupRates(ctx, ListHertzMarkupRatesRequest{Page: 1, Limit: 10})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateHertzMarkupRate(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing country", func(t *testing.T) {
		p := validCreateParams()
		p.Country = ""
		api_errors.AssertApiError(t, invalidValueErr("country"), p.Validate())
	})

	t.Run("validation rejects missing brand", func(t *testing.T) {
		p := validCreateParams()
		p.Brand = ""
		api_errors.AssertApiError(t, invalidValueErr("brand"), p.Validate())
	})

	t.Run("validation rejects invalid pickup date", func(t *testing.T) {
		p := validCreateParams()
		p.PickupDateFrom = "not-a-date"
		api_errors.AssertApiError(t, invalidValueErr("pickupDateFrom"), p.Validate())
	})

	t.Run("validation rejects missing car group", func(t *testing.T) {
		p := validCreateParams()
		p.CarGroup = ""
		api_errors.AssertApiError(t, invalidValueErr("carGroup"), p.Validate())
	})

	t.Run("validation rejects rental days from zero", func(t *testing.T) {
		p := validCreateParams()
		p.NumOfRentalDaysFrom = 0
		api_errors.AssertApiError(t, invalidValueErr("numOfRentalDaysFrom"), p.Validate())
	})

	t.Run("validation rejects rental days to less than from", func(t *testing.T) {
		p := validCreateParams()
		p.NumOfRentalDaysFrom = 5
		p.NumOfRentalDaysTo = 3
		api_errors.AssertApiError(t, invalidValueErr("numOfRentalDaysTo"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validCreateParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("creates rate successfully", func(t *testing.T) {
		resp, err := s.CreateHertzMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if resp.Country != "US" {
			t.Fatalf("expected country US, got %s", resp.Country)
		}
		if resp.PickupDateFrom != "2026-01-01" {
			t.Fatalf("expected pickup date from 2026-01-01, got %s", resp.PickupDateFrom)
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().InsertHertzMarkupRate(gomock.Any(), gomock.Any()).Return(db.HertzMarkupRate{}, errors.New("db error"))

		_, err := s.CreateHertzMarkupRate(ctx, validCreateParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestUpdateHertzMarkupRate(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing country", func(t *testing.T) {
		p := validUpdateParams()
		p.Country = ""
		api_errors.AssertApiError(t, invalidValueErr("country"), p.Validate())
	})

	t.Run("validation rejects invalid pickup date", func(t *testing.T) {
		p := validUpdateParams()
		p.PickupDateFrom = "bad"
		api_errors.AssertApiError(t, invalidValueErr("pickupDateFrom"), p.Validate())
	})

	t.Run("validation rejects rental days to less than from", func(t *testing.T) {
		p := validUpdateParams()
		p.NumOfRentalDaysFrom = 10
		p.NumOfRentalDaysTo = 5
		api_errors.AssertApiError(t, invalidValueErr("numOfRentalDaysTo"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validUpdateParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("updates rate successfully", func(t *testing.T) {
		created, err := s.CreateHertzMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}

		resp, err := s.UpdateHertzMarkupRate(ctx, created.ID, validUpdateParams())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.MarkUpGross != 20.0 {
			t.Fatalf("expected markup gross 20.0, got %f", resp.MarkUpGross)
		}
		if resp.MarkUpNet != 12.0 {
			t.Fatalf("expected markup net 12.0, got %f", resp.MarkUpNet)
		}
	})

	t.Run("returns not found when rate does not exist", func(t *testing.T) {
		_, err := s.UpdateHertzMarkupRate(ctx, 999999, validUpdateParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().UpdateHertzMarkupRate(gomock.Any(), gomock.Any()).Return(db.HertzMarkupRate{}, errors.New("db error"))

		_, err := s.UpdateHertzMarkupRate(ctx, 1, validUpdateParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestDeleteHertzMarkupRate(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("deletes rate successfully", func(t *testing.T) {
		created, err := s.CreateHertzMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}

		if err := s.DeleteHertzMarkupRate(ctx, created.ID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's gone
		_, err = s.UpdateHertzMarkupRate(ctx, created.ID, validUpdateParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns not found when rate does not exist", func(t *testing.T) {
		api_errors.AssertApiError(t, api_errors.ErrNotFound, s.DeleteHertzMarkupRate(ctx, 999999))
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteHertzMarkupRate(gomock.Any(), int64(1)).Return(int64(0), errors.New("db error"))

		api_errors.AssertApiError(t, api_errors.ErrInternalError, s.DeleteHertzMarkupRate(ctx, 1))
	})
}
