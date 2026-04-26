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

func TestCreateFirstAdmin(t *testing.T) {
	ctx := context.Background()

	// Save original secrets and restore after tests
	originalUsername := secrets.FirstAdminEmail
	originalPassword := secrets.FirstAdminPassword
	defer func() {
		secrets.FirstAdminEmail = originalUsername
		secrets.FirstAdminPassword = originalPassword
	}()

	t.Run("Secrets not set", func(t *testing.T) {
		secrets.FirstAdminEmail = ""
		secrets.FirstAdminPassword = ""

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		createFirstAdmin(query)
	})

	t.Run("Success And User already exists", func(t *testing.T) {
		email := "admin@example.com"
		secrets.FirstAdminEmail = email
		secrets.FirstAdminPassword = "password123"

		// validate admin not exists yet:
		admin, err := query.GetUserByEmail(ctx, email)
		if err != nil && !errors.Is(err, db.ErrNoRows) {
			t.Fatalf("failed to get user by email: %v", err)
		}
		if err == nil {
			t.Fatalf("expected no admin user, but found one: %v", admin)
		}

		// success, should create admin user
		createFirstAdmin(query)

		// validate admin was created:
		admin, err = query.GetUserByEmail(ctx, email)
		if err != nil {
			t.Fatalf("failed to get user by email after creation: %v", err)
		}
		if admin.Email != email {
			t.Errorf("expected email %s, got %s", email, admin.Email)
		}

		// user already exists, should not panic or create another user
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The code panicked when it should not have")
			}
		}()

		createFirstAdmin(query)
	})

	t.Run("Database error checking user", func(t *testing.T) {
		secrets.FirstAdminEmail = "admin@example.com"
		secrets.FirstAdminPassword = "password123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		expectedErr := errors.New("db error")
		q.EXPECT().
			CheckUserExists(gomock.Any(), secrets.FirstAdminEmail).
			Return(int32(0), expectedErr)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		createFirstAdmin(q)
	})
}

// --- Helpers ---

func adminMockService(t *testing.T) (*mocks.MockQuerier, *Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	q := mocks.NewMockQuerier(ctrl)
	return q, &Service{query: q}
}

func createTestAdmin(t *testing.T, s *Service, email string) *CreateAdminResponse {
	t.Helper()
	resp, err := s.CreateAdmin(context.Background(), CreateAdminRequest{
		FirstName: "Test",
		LastName:  "Admin",
		Email:     email,
		Password:  "ValidPass123!",
	})
	if err != nil {
		t.Fatalf("failed to create admin %s: %v", email, err)
	}
	return resp
}

// --- Tests ---

func TestListAdmins(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("returns created admin in list", func(t *testing.T) {
		t.Parallel()
		admin := createTestAdmin(t, s, "list_admin_1@test.com")
		defer query.DeleteUser(ctx, admin.ID)

		resp, err := s.ListAdmins(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found := false
		for _, a := range resp.Admins {
			if a.ID == admin.ID {
				found = true
				if a.Email != "list_admin_1@test.com" {
					t.Fatalf("expected email list_admin_1@test.com, got %s", a.Email)
				}
			}
		}
		if !found {
			t.Fatal("created admin not found in list")
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		t.Parallel()
		q, s := adminMockService(t)
		q.EXPECT().ListAdmins(gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListAdmins(ctx)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateAdmin(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: query}

	t.Run("creates admin successfully", func(t *testing.T) {
		t.Parallel()
		resp, err := s.CreateAdmin(ctx, CreateAdminRequest{
			FirstName: "Create",
			LastName:  "Ok",
			Email:     "create_admin_ok@test.com",
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
		admin := createTestAdmin(t, s, "dup_admin@test.com")
		defer query.DeleteUser(ctx, admin.ID)

		_, err := s.CreateAdmin(ctx, CreateAdminRequest{
			FirstName: "Dup",
			LastName:  "Admin",
			Email:     "dup_admin@test.com",
			Password:  "ValidPass123!",
		})
		api_errors.AssertApiError(t, ErrEmailAlreadyExists, err)
	})

	t.Run("validation rejects empty firstName", func(t *testing.T) {
		t.Parallel()
		p := CreateAdminRequest{FirstName: "", LastName: "Admin", Email: "admin@test.com", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("firstName"), p.Validate())
	})

	t.Run("validation rejects empty lastName", func(t *testing.T) {
		t.Parallel()
		p := CreateAdminRequest{FirstName: "Test", LastName: "", Email: "admin@test.com", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("lastName"), p.Validate())
	})

	t.Run("validation rejects invalid email", func(t *testing.T) {
		t.Parallel()
		p := CreateAdminRequest{FirstName: "Test", LastName: "Admin", Email: "not-an-email", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("validation rejects empty email", func(t *testing.T) {
		t.Parallel()
		p := CreateAdminRequest{FirstName: "Test", LastName: "Admin", Email: "", Password: "ValidPass123!"}
		api_errors.AssertApiError(t, invalidValueErr("email"), p.Validate())
	})

	t.Run("validation rejects weak password", func(t *testing.T) {
		t.Parallel()
		p := CreateAdminRequest{FirstName: "Test", LastName: "Admin", Email: "weak_pw@test.com", Password: "short"}
		api_errors.AssertApiError(t, ErrPasswordTooShort, p.Validate())
	})

	t.Run("validation rejects password without uppercase", func(t *testing.T) {
		t.Parallel()
		p := CreateAdminRequest{FirstName: "Test", LastName: "Admin", Email: "no_upper@test.com", Password: "validpass123!"}
		api_errors.AssertApiError(t, ErrPasswordNoUpper, p.Validate())
	})

	t.Run("returns error when check exists db fails", func(t *testing.T) {
		t.Parallel()
		q, s := adminMockService(t)
		q.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Return(int32(0), errors.New("db error"))

		_, err := s.CreateAdmin(ctx, CreateAdminRequest{
			FirstName: "DB",
			LastName:  "Fail",
			Email:     "db_fail@test.com",
			Password:  "ValidPass123!",
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when create db fails", func(t *testing.T) {
		t.Parallel()
		q, s := adminMockService(t)
		q.EXPECT().CheckUserExists(gomock.Any(), gomock.Any()).Return(int32(0), db.ErrNoRows)
		q.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(db.CreateAdminRow{}, errors.New("db error"))

		_, err := s.CreateAdmin(ctx, CreateAdminRequest{
			FirstName: "DB",
			LastName:  "CreateFail",
			Email:     "db_create_fail@test.com",
			Password:  "ValidPass123!",
		})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
