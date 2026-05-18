package accounts

import (
	"context"
	"testing"

	"encore.app/internal/api_errors"
	user "encore.app/services/accounts/handlers/user"
)

// --- Helpers ---

func createTestAccountant(t *testing.T, s *Service, email string) *user.CreateAccountantResponse {
	t.Helper()
	resp, err := s.CreateAccountant(context.Background(), user.CreateAccountantParams{
		FirstName: "Test",
		LastName:  "Accountant",
		Email:     email,
		Password:  "ValidPass123!",
	})
	if err != nil {
		t.Fatalf("failed to create accountant %s: %v", email, err)
	}
	return resp
}

// --- Tests ---

func TestListAccountants(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("returns created accountant in list", func(t *testing.T) {
		t.Parallel()
		acc := createTestAccountant(t, s, "list_accountant_1@test.com")
		defer query.DeleteUser(ctx, acc.ID)

		resp, err := s.ListAccountants(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found := false
		for _, a := range resp.Accountants {
			if a.ID == acc.ID {
				found = true
				if a.Email != "list_accountant_1@test.com" {
					t.Fatalf("expected email list_accountant_1@test.com, got %s", a.Email)
				}
			}
		}
		if !found {
			t.Fatal("created accountant not found in list")
		}
	})

}

func TestCreateAccountant(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("creates accountant successfully", func(t *testing.T) {
		t.Parallel()
		resp, err := s.CreateAccountant(ctx, user.CreateAccountantParams{
			FirstName: "Create",
			LastName:  "Ok",
			Email:     "create_accountant_ok@test.com",
			Password:  "ValidPass123!",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer query.DeleteUser(ctx, resp.ID)

		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
	})

	t.Run("returns error on duplicate email", func(t *testing.T) {
		t.Parallel()
		acc := createTestAccountant(t, s, "dup_accountant@test.com")
		defer query.DeleteUser(ctx, acc.ID)

		_, err := s.CreateAccountant(ctx, user.CreateAccountantParams{
			FirstName: "Dup",
			LastName:  "Accountant",
			Email:     "dup_accountant@test.com",
			Password:  "ValidPass123!",
		})
		api_errors.AssertApiError(t, user.ErrEmailAlreadyExists, err)
	})

	t.Run("validation rejects empty firstName", func(t *testing.T) {
		t.Parallel()
		p := user.CreateAccountantParams{FirstName: "", LastName: "Accountant", Email: "acc@test.com", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("firstName"), p.Validate())
	})

	t.Run("validation rejects empty lastName", func(t *testing.T) {
		t.Parallel()
		p := user.CreateAccountantParams{FirstName: "Test", LastName: "", Email: "acc@test.com", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("lastName"), p.Validate())
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		p := user.CreateAccountantParams{FirstName: "Test", LastName: "Accountant", Email: "not-an-email", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("validation rejects weak password", func(t *testing.T) {
		t.Parallel()
		p := user.CreateAccountantParams{FirstName: "Test", LastName: "Accountant", Email: "acc@test.com", Password: "short"}
		api_errors.AssertApiError(t, user.ErrPasswordTooShort, p.Validate())
	})

	t.Run("validation rejects password without uppercase", func(t *testing.T) {
		t.Parallel()
		p := user.CreateAccountantParams{FirstName: "Test", LastName: "Accountant", Email: "acc@test.com", Password: "validpass123!"}
		api_errors.AssertApiError(t, user.ErrPasswordNoUpper, p.Validate())
	})

}
