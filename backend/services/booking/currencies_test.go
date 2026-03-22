package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func validCreateCurrencyParams() CreateCurrencyRequest {
	return CreateCurrencyRequest{
		CurrencyCode:    "USD",
		CurrencyISOName: "US Dollar",
		Rate:            3.65,
	}
}

func validUpdateCurrencyParams() UpdateCurrencyRequest {
	code := "EUR"
	name := "Euro"
	rate := 4.12
	return UpdateCurrencyRequest{
		CurrencyCode:    &code,
		CurrencyISOName: &name,
		Rate:            &rate,
	}
}

func currencyInvalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

// createTestCurrency is a shorthand to seed a currency with a unique code and ISO name.
func createTestCurrency(t *testing.T, s *Service, code, isoName string) *CurrencyResponse {
	t.Helper()
	p := validCreateCurrencyParams()
	p.CurrencyCode = code
	p.CurrencyISOName = isoName
	resp, err := s.CreateCurrency(context.Background(), p)
	if err != nil {
		t.Fatalf("failed to seed currency %s: %v", code, err)
	}
	return resp
}

// --- Tests grouped by endpoint ---

func TestListCurrencies(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("returns currencies successfully", func(t *testing.T) {
		c1 := createTestCurrency(t, s, "LIST-A", "List Currency A")
		c2 := createTestCurrency(t, s, "LIST-B", "List Currency B")
		c3 := createTestCurrency(t, s, "LIST-C", "List Currency C")

		resp, err := s.ListCurrencies(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		ids := make(map[int32]bool)
		for _, c := range resp.Currencies {
			ids[c.ID] = true
		}
		for _, want := range []*CurrencyResponse{c1, c2, c3} {
			if !ids[want.ID] {
				t.Fatalf("expected currency %d (%s) in list", want.ID, want.CurrencyCode)
			}
		}
	})

	t.Run("returns empty list when no currencies exist", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().ListCurrencies(gomock.Any()).Return([]db.Currency{}, nil)

		resp, err := s.ListCurrencies(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Currencies) != 0 {
			t.Fatalf("expected 0 currencies, got %d", len(resp.Currencies))
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().ListCurrencies(gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListCurrencies(ctx)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateCurrency(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing currency code", func(t *testing.T) {
		p := validCreateCurrencyParams()
		p.CurrencyCode = ""
		api_errors.AssertApiError(t, currencyInvalidValueErr("currencyCode"), p.Validate())
	})

	t.Run("validation rejects blank currency code", func(t *testing.T) {
		p := validCreateCurrencyParams()
		p.CurrencyCode = "   "
		api_errors.AssertApiError(t, currencyInvalidValueErr("currencyCode"), p.Validate())
	})

	t.Run("validation rejects missing currency ISO name", func(t *testing.T) {
		p := validCreateCurrencyParams()
		p.CurrencyISOName = ""
		api_errors.AssertApiError(t, currencyInvalidValueErr("currencyISOName"), p.Validate())
	})

	t.Run("validation rejects blank currency ISO name", func(t *testing.T) {
		p := validCreateCurrencyParams()
		p.CurrencyISOName = "   "
		api_errors.AssertApiError(t, currencyInvalidValueErr("currencyISOName"), p.Validate())
	})

	t.Run("validation rejects zero rate", func(t *testing.T) {
		p := validCreateCurrencyParams()
		p.Rate = 0
		api_errors.AssertApiError(t, currencyInvalidValueErr("rate"), p.Validate())
	})

	t.Run("validation rejects negative rate", func(t *testing.T) {
		p := validCreateCurrencyParams()
		p.Rate = -1.5
		api_errors.AssertApiError(t, currencyInvalidValueErr("rate"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validCreateCurrencyParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("creates currency successfully", func(t *testing.T) {
		resp, err := s.CreateCurrency(ctx, CreateCurrencyRequest{
			CurrencyCode:    "CREATE-OK",
			CurrencyISOName: "Create OK Currency",
			Rate:            3.75,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if resp.CurrencyCode != "CREATE-OK" {
			t.Fatalf("expected currency code 'CREATE-OK', got %q", resp.CurrencyCode)
		}
		if resp.CurrencyISOName != "Create OK Currency" {
			t.Fatalf("expected currency ISO name 'Create OK Currency', got %q", resp.CurrencyISOName)
		}
		if resp.Rate != 3.75 {
			t.Fatalf("expected rate 3.75, got %f", resp.Rate)
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().CreateCurrency(gomock.Any(), gomock.Any()).Return(db.Currency{}, errors.New("db error"))

		_, err := s.CreateCurrency(ctx, validCreateCurrencyParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestUpdateCurrency(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects blank currency code", func(t *testing.T) {
		p := validUpdateCurrencyParams()
		blank := "   "
		p.CurrencyCode = &blank
		api_errors.AssertApiError(t, currencyInvalidValueErr("currencyCode"), p.Validate())
	})

	t.Run("validation rejects blank currency ISO name", func(t *testing.T) {
		p := validUpdateCurrencyParams()
		blank := "   "
		p.CurrencyISOName = &blank
		api_errors.AssertApiError(t, currencyInvalidValueErr("currencyISOName"), p.Validate())
	})

	t.Run("validation rejects zero rate", func(t *testing.T) {
		p := validUpdateCurrencyParams()
		zero := 0.0
		p.Rate = &zero
		api_errors.AssertApiError(t, currencyInvalidValueErr("rate"), p.Validate())
	})

	t.Run("validation rejects negative rate", func(t *testing.T) {
		p := validUpdateCurrencyParams()
		neg := -2.5
		p.Rate = &neg
		api_errors.AssertApiError(t, currencyInvalidValueErr("rate"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validUpdateCurrencyParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("validation accepts partial update with only code", func(t *testing.T) {
		code := "GBP"
		p := UpdateCurrencyRequest{CurrencyCode: &code}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("updates currency successfully", func(t *testing.T) {
		created := createTestCurrency(t, s, "UPD-BEFORE", "Update Before Currency")

		resp, err := s.UpdateCurrency(ctx, created.ID, validUpdateCurrencyParams())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.CurrencyCode != "EUR" {
			t.Fatalf("expected currency code 'EUR', got %q", resp.CurrencyCode)
		}
		if resp.CurrencyISOName != "Euro" {
			t.Fatalf("expected currency ISO name 'Euro', got %q", resp.CurrencyISOName)
		}
		if resp.Rate != 4.12 {
			t.Fatalf("expected rate 4.12, got %f", resp.Rate)
		}
	})

	t.Run("partial update only changes provided fields", func(t *testing.T) {
		created := createTestCurrency(t, s, "PRTL", "Partial Currency")

		newCode := "PRTL-NEW"
		resp, err := s.UpdateCurrency(ctx, created.ID, UpdateCurrencyRequest{CurrencyCode: &newCode})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.CurrencyCode != "PRTL-NEW" {
			t.Fatalf("expected currency code 'PRTL-NEW', got %q", resp.CurrencyCode)
		}
		// Other fields should remain unchanged
		if resp.CurrencyISOName != "Partial Currency" {
			t.Fatalf("expected currency ISO name unchanged 'Partial Currency', got %q", resp.CurrencyISOName)
		}
		if resp.Rate != created.Rate {
			t.Fatalf("expected rate unchanged %f, got %f", created.Rate, resp.Rate)
		}
	})

	t.Run("returns not found when currency does not exist", func(t *testing.T) {
		_, err := s.UpdateCurrency(ctx, 999999, validUpdateCurrencyParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().UpdateCurrency(gomock.Any(), gomock.Any()).Return(db.Currency{}, errors.New("db error"))

		_, err := s.UpdateCurrency(ctx, 1, validUpdateCurrencyParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestDeleteCurrency(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("deletes currency successfully", func(t *testing.T) {
		created := createTestCurrency(t, s, "DEL-OK", "Delete OK Currency")

		if err := s.DeleteCurrency(ctx, created.ID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's gone by trying to update it
		_, err := s.UpdateCurrency(ctx, created.ID, validUpdateCurrencyParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteCurrency(gomock.Any(), int32(1)).Return(errors.New("db error"))

		api_errors.AssertApiError(t, api_errors.ErrInternalError, s.DeleteCurrency(ctx, 1))
	})
}
