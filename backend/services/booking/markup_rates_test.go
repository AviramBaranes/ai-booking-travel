package booking

import (
	"context"
	"testing"

	"encore.app/internal/api_errors"
	markup_rate "encore.app/services/booking/handlers/markup_rate"
)

// --- Helpers ---

func validCreateParams() markup_rate.CreateMarkupRateParams {
	return markup_rate.CreateMarkupRateParams{
		CountryCode: "US",
		Broker:      "flex",
		MarkUpGross: 15.5,
		MarkUpNet:   10.0,
	}
}

func validUpdateParams() markup_rate.UpdateMarkupRateParams {
	return markup_rate.UpdateMarkupRateParams{
		CountryCode: "US",
		Broker:      "flex",
		MarkUpGross: 20.0,
		MarkUpNet:   12.0,
	}
}

// --- Tests grouped by endpoint ---

func TestListMarkupRates(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing page", func(t *testing.T) {
		api_errors.AssertApiError(t, invalidValueErr("page"), (markup_rate.ListMarkupRatesParams{}).Validate())
	})

	t.Run("validation rejects invalid sort direction", func(t *testing.T) {
		api_errors.AssertApiError(t, invalidValueErr("sortDir"), (markup_rate.ListMarkupRatesParams{Page: 1, SortDir: "up"}).Validate())
	})

	t.Run("validation rejects invalid sort field", func(t *testing.T) {
		api_errors.AssertApiError(t, api_errors.ErrInvalidValue, (markup_rate.ListMarkupRatesParams{Page: 1, SortBy: "invalid_field"}).Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := (markup_rate.ListMarkupRatesParams{Country: "US", SortBy: "country", SortDir: "desc", Page: 1}).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("returns rates successfully", func(t *testing.T) {
		created, err := s.CreateMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}

		resp, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Country: "US", Page: 1})
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
		resp, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Country: "XX", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Rates) != 0 {
			t.Fatalf("expected 0 rates, got %d", len(resp.Rates))
		}
	})

	t.Run("filters by country", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			p := validCreateParams()
			p.CountryCode = "DE"
			p.MarkUpGross = float64(10 + i)
			if _, err := s.CreateMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed DE: %v", err)
			}
		}
		for i := 0; i < 2; i++ {
			p := validCreateParams()
			p.CountryCode = "IT"
			p.MarkUpGross = float64(10 + i)
			if _, err := s.CreateMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed IT: %v", err)
			}
		}

		resp, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Country: "DE", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Rates) != 3 {
			t.Fatalf("expected 3 DE rates, got %d", len(resp.Rates))
		}
		for _, r := range resp.Rates {
			if r.CountryCode != "DE" {
				t.Fatalf("expected countryCode DE, got %s", r.CountryCode)
			}
		}
	})

	t.Run("filters by broker", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			p := validCreateParams()
			p.CountryCode = "BF"
			p.Broker = "flex"
			p.MarkUpGross = float64(10 + i)
			if _, err := s.CreateMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed flex: %v", err)
			}
		}
		for i := 0; i < 2; i++ {
			p := validCreateParams()
			p.CountryCode = "BF"
			p.Broker = "hertz"
			p.MarkUpGross = float64(10 + i)
			if _, err := s.CreateMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed hertz: %v", err)
			}
		}

		resp, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Country: "BF", Broker: "flex", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Rates) != 2 {
			t.Fatalf("expected 2 flex rates, got %d", len(resp.Rates))
		}
		for _, r := range resp.Rates {
			if r.Broker != "flex" {
				t.Fatalf("expected broker=flex, got %s", r.Broker)
			}
		}
	})

	t.Run("paginates results", func(t *testing.T) {
		// Seed 16 rates under unique filter so page 1 = 15, page 2 = 1
		for i := 0; i < 16; i++ {
			p := validCreateParams()
			p.CountryCode = "PG"
			p.MarkUpGross = float64(10 + i)
			if _, err := s.CreateMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed: %v", err)
			}
		}

		page1, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Country: "PG", Page: 1})
		if err != nil {
			t.Fatalf("page 1 error: %v", err)
		}
		if len(page1.Rates) != 15 {
			t.Fatalf("expected 15 rates on page 1, got %d", len(page1.Rates))
		}

		page2, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Country: "PG", Page: 2})
		if err != nil {
			t.Fatalf("page 2 error: %v", err)
		}
		if len(page2.Rates) != 1 {
			t.Fatalf("expected 1 rate on page 2, got %d", len(page2.Rates))
		}

		page3, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Country: "PG", Page: 3})
		if err != nil {
			t.Fatalf("page 3 error: %v", err)
		}
		if len(page3.Rates) != 0 {
			t.Fatalf("expected 0 rates on page 3, got %d", len(page3.Rates))
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
		}
	})

	t.Run("defaults sort to country ascending", func(t *testing.T) {
		for _, c := range []string{"FR", "CA", "GB"} {
			p := validCreateParams()
			p.CountryCode = c
			p.Broker = "hertz"
			if _, err := s.CreateMarkupRate(ctx, p); err != nil {
				t.Fatalf("failed to seed: %v", err)
			}
		}

		resp, err := s.ListMarkupRates(ctx, markup_rate.ListMarkupRatesParams{Broker: "hertz", Page: 1})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for i := 1; i < len(resp.Rates); i++ {
			if resp.Rates[i].CountryCode < resp.Rates[i-1].CountryCode {
				t.Fatalf("expected ascending country order, got %s after %s", resp.Rates[i].CountryCode, resp.Rates[i-1].CountryCode)
			}
		}
	})

}

func TestCreateMarkupRate(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing country code", func(t *testing.T) {
		p := validCreateParams()
		p.CountryCode = ""
		api_errors.AssertApiError(t, invalidValueErr("countryCode"), p.Validate())
	})

	t.Run("validation rejects missing broker", func(t *testing.T) {
		p := validCreateParams()
		p.Broker = ""
		api_errors.AssertApiError(t, invalidValueErr("broker"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validCreateParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("creates rate successfully", func(t *testing.T) {
		resp, err := s.CreateMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if resp.CountryCode != "US" {
			t.Fatalf("expected countryCode US, got %s", resp.CountryCode)
		}
		if resp.Broker != "flex" {
			t.Fatalf("expected broker flex, got %s", resp.Broker)
		}
	})

}

func TestUpdateMarkupRate(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing country", func(t *testing.T) {
		p := validUpdateParams()
		p.CountryCode = ""
		api_errors.AssertApiError(t, invalidValueErr("country"), p.Validate())
	})

	t.Run("validation rejects missing broker", func(t *testing.T) {
		p := validUpdateParams()
		p.Broker = ""
		api_errors.AssertApiError(t, invalidValueErr("broker"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validUpdateParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("updates rate successfully", func(t *testing.T) {
		created, err := s.CreateMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}

		resp, err := s.UpdateMarkupRate(ctx, created.ID, validUpdateParams())
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
		_, err := s.UpdateMarkupRate(ctx, 999999, validUpdateParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

}

func TestDeleteMarkupRate(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("deletes rate successfully", func(t *testing.T) {
		created, err := s.CreateMarkupRate(ctx, validCreateParams())
		if err != nil {
			t.Fatalf("failed to seed: %v", err)
		}

		if err := s.DeleteMarkupRate(ctx, created.ID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's gone
		_, err = s.UpdateMarkupRate(ctx, created.ID, validUpdateParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns not found when rate does not exist", func(t *testing.T) {
		api_errors.AssertApiError(t, api_errors.ErrNotFound, s.DeleteMarkupRate(ctx, 999999))
	})

}
