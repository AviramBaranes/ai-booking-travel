package accounts

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/mocks"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func accountantMockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func createTestAccountant(t *testing.T, s *Service, email string) *CreateAccountantResponse {
	t.Helper()
	resp, err := s.CreateAccountant(context.Background(), CreateAccountantRequest{
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

	t.Run("returns error when db fails", func(t *testing.T) {
		t.Parallel()
		q, s := accountantMockService(t)
		q.EXPECT().ListStaffByRole(gomock.Any(), db.UserRoleAccountant).Return(nil, errors.New("db error"))

		_, err := s.ListAccountants(ctx)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateAccountant(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("creates accountant successfully", func(t *testing.T) {
		t.Parallel()
		resp, err := s.CreateAccountant(ctx, CreateAccountantRequest{
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

		_, err := s.CreateAccountant(ctx, CreateAccountantRequest{
			FirstName: "Dup",
			LastName:  "Accountant",
			Email:     "dup_accountant@test.com",
			Password:  "ValidPass123!",
		})
		api_errors.AssertApiError(t, ErrEmailAlreadyExists, err)
	})

	t.Run("validation rejects empty firstName", func(t *testing.T) {
		t.Parallel()
		p := CreateAccountantRequest{FirstName: "", LastName: "Accountant", Email: "acc@test.com", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("firstName"), p.Validate())
	})

	t.Run("validation rejects empty lastName", func(t *testing.T) {
		t.Parallel()
		p := CreateAccountantRequest{FirstName: "Test", LastName: "", Email: "acc@test.com", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("lastName"), p.Validate())
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		p := CreateAccountantRequest{FirstName: "Test", LastName: "Accountant", Email: "not-an-email", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("validation rejects weak password", func(t *testing.T) {
		t.Parallel()
		p := CreateAccountantRequest{FirstName: "Test", LastName: "Accountant", Email: "acc@test.com", Password: "short"}
		api_errors.AssertApiError(t, ErrPasswordTooShort, p.Validate())
	})

	t.Run("validation rejects password without uppercase", func(t *testing.T) {
		t.Parallel()
		p := CreateAccountantRequest{FirstName: "Test", LastName: "Accountant", Email: "acc@test.com", Password: "validpass123!"}
		api_errors.AssertApiError(t, ErrPasswordNoUpper, p.Validate())
	})

	t.Run("returns error when check exists db fails", func(t *testing.T) {
		t.Parallel()
		q, s := accountantMockService(t)
		q.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Return(int32(0), errors.New("db error"))

		_, err := s.CreateAccountant(ctx, CreateAccountantRequest{
			FirstName: "DB",
			LastName:  "Fail",
			Email:     "db_fail_acc@test.com",
			Password:  "ValidPass123!",
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when create db fails", func(t *testing.T) {
		t.Parallel()
		q, s := accountantMockService(t)
		q.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Return(int32(0), db.ErrNoRows)
		q.EXPECT().CreateStaffUser(gomock.Any(), gomock.Any()).Return(db.CreateStaffUserRow{}, errors.New("db error"))

		_, err := s.CreateAccountant(ctx, CreateAccountantRequest{
			FirstName: "DB",
			LastName:  "CreateFail",
			Email:     "db_create_fail_acc@test.com",
			Password:  "ValidPass123!",
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
