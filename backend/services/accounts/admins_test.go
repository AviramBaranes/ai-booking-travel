package accounts

import (
	"context"
	"errors"
	"testing"

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
