package accounts

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	ch "encore.app/services/accounts/handlers/customer"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
)

// --- Helpers ---

func createTestCustomer(t *testing.T, ctx context.Context, suffix string) ch.CustomerResponse {
	t.Helper()
	email := fmt.Sprintf("%s@test.com", suffix)
	phone := randomIsraeliPhoneNumber()
	row, err := query.CreateCustomer(ctx, db.CreateCustomerParams{
		FirstName:    "Fname_" + suffix,
		LastName:     "Lname_" + suffix,
		Email:        email,
		PhoneNumber:  &phone,
		PasswordHash: "nohash",
	})
	if err != nil {
		t.Fatalf("createTestCustomer: %v", err)
	}
	t.Cleanup(func() { query.DeleteUser(ctx, row.ID) })
	return ch.CustomerResponse{
		ID:        row.ID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
	}
}

// --- Tests ---

func TestListCustomers(t *testing.T) {
	ctx := context.Background()

	// ── Validation ──

	t.Run("validation: page required and >=1", func(t *testing.T) {
		t.Parallel()
		wantErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
			Code: api_errors.CodeInvalidValue, Field: "page",
		})

		p := ch.ListCustomersParams{Search: "", Page: 0}
		err := p.Validate()
		api_errors.AssertApiError(t, wantErr, err)
	})

	t.Run("validation: page 1 passes", func(t *testing.T) {
		t.Parallel()
		p := ch.ListCustomersParams{Page: 1}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no validation error, got %v", err)
		}
	})

	// ── No rows ──

	t.Run("returns empty list when no customers match search", func(t *testing.T) {
		t.Parallel()

		resp, err := ListCustomers(ctx, ch.ListCustomersParams{
			Search: "zzznomatch_xyzabc_unique",
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Customers) != 0 {
			t.Fatalf("expected 0 customers, got %d", len(resp.Customers))
		}
		if resp.Total != 0 {
			t.Fatalf("expected total 0, got %d", resp.Total)
		}
	})

	// ── Filtered rows ──

	t.Run("returns only matching customers when search is applied", func(t *testing.T) {
		t.Parallel()

		uniqueTag := fmt.Sprintf("filtercust_%d", time.Now().UnixNano())
		c := createTestCustomer(t, ctx, uniqueTag)

		// Create a second customer that should NOT appear.
		createTestCustomer(t, ctx, fmt.Sprintf("other_%d", time.Now().UnixNano()))

		resp, err := ListCustomers(ctx, ch.ListCustomersParams{
			Search: c.Email,
			Page:   1,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Customers) != 1 {
			t.Fatalf("expected 1 matching customer, got %d", len(resp.Customers))
		}
		if resp.Customers[0].ID != c.ID {
			t.Fatalf("expected customer ID %d, got %d", c.ID, resp.Customers[0].ID)
		}
		if resp.Total != 1 {
			t.Fatalf("expected total 1, got %d", resp.Total)
		}
	})

	// ── All rows / pagination ──

	t.Run("paginates: returns first page of 15 when more than 15 exist", func(t *testing.T) {
		t.Parallel()

		paginationTag := fmt.Sprintf("pagcust_%d", time.Now().UnixNano())
		for i := 0; i < 18; i++ {
			createTestCustomer(t, ctx, fmt.Sprintf("%s_%02d", paginationTag, i))
		}

		page1, err := ListCustomers(ctx, ch.ListCustomersParams{
			Search: paginationTag,
			Page:   1,
		})
		if err != nil {
			t.Fatalf("page 1: %v", err)
		}
		if len(page1.Customers) != 15 {
			t.Fatalf("expected 15 customers on page 1, got %d", len(page1.Customers))
		}
		if page1.Total != 18 {
			t.Fatalf("expected total 18, got %d", page1.Total)
		}

		page2, err := ListCustomers(ctx, ch.ListCustomersParams{
			Search: paginationTag,
			Page:   2,
		})
		if err != nil {
			t.Fatalf("page 2: %v", err)
		}
		if len(page2.Customers) != 3 {
			t.Fatalf("expected 3 customers on page 2, got %d", len(page2.Customers))
		}
	})
}
func customerAuthContext(userID int64) context.Context {
	uid := auth.UID(strconv.FormatInt(userID, 10))
	return auth.WithContext(context.Background(), uid, &AuthData{
		UserID: userID,
		Role:   UserRoleCustomer,
	})
}

func TestUpdateCustomer(t *testing.T) {
	ctx := context.Background()

	// ── Validation ──

	validationCases := []struct {
		name   string
		params ch.UpdateCustomerParams
		field  string
	}{
		{"firstName required", ch.UpdateCustomerParams{LastName: "Last", Email: "a@b.com"}, "firstName"},
		{"lastName required", ch.UpdateCustomerParams{FirstName: "First", Email: "a@b.com"}, "lastName"},
		{"email required", ch.UpdateCustomerParams{FirstName: "First", LastName: "Last"}, "email"},
		{"email invalid format", ch.UpdateCustomerParams{FirstName: "First", LastName: "Last", Email: "notanemail"}, "email"},
	}
	for _, tc := range validationCases {
		tc := tc
		t.Run("validation: "+tc.name, func(t *testing.T) {
			t.Parallel()
			wantErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
				Code: api_errors.CodeInvalidValue, Field: tc.field,
			})
			api_errors.AssertApiError(t, wantErr, tc.params.Validate())
		})
	}

	// ── Success ──

	t.Run("updates customer fields successfully", func(t *testing.T) {
		t.Parallel()
		suffix := fmt.Sprintf("upd_%d", time.Now().UnixNano())
		c := createTestCustomer(t, ctx, suffix)

		authCtx := customerAuthContext(c.ID)
		newEmail := fmt.Sprintf("updated_%d@test.com", time.Now().UnixNano())
		err := UpdateCustomer(authCtx, ch.UpdateCustomerParams{
			FirstName: "NewFirst",
			LastName:  "NewLast",
			Email:     newEmail,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// verify the change persisted
		row, dbErr := query.GetCustomerByID(ctx, c.ID)
		if dbErr != nil {
			t.Fatalf("failed to fetch updated customer: %v", dbErr)
		}
		if row.FirstName != "NewFirst" {
			t.Fatalf("expected FirstName 'NewFirst', got %q", row.FirstName)
		}
		if row.LastName != "NewLast" {
			t.Fatalf("expected LastName 'NewLast', got %q", row.LastName)
		}
		if row.Email != newEmail {
			t.Fatalf("expected Email %q, got %q", newEmail, row.Email)
		}
	})

	// ── Not found ──

	t.Run("returns internal error when user id does not exist", func(t *testing.T) {
		t.Parallel()
		authCtx := customerAuthContext(999999999)
		err := UpdateCustomer(authCtx, ch.UpdateCustomerParams{
			FirstName: "Ghost",
			LastName:  "User",
			Email:     "ghost@test.com",
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
