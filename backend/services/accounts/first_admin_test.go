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
	originalUsername := secrets.FirstAdminUsername
	originalPassword := secrets.FirstAdminPassword
	defer func() {
		secrets.FirstAdminUsername = originalUsername
		secrets.FirstAdminPassword = originalPassword
	}()

	t.Run("Secrets not set", func(t *testing.T) {
		secrets.FirstAdminUsername = ""
		secrets.FirstAdminPassword = ""

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		createFirstAdmin(query)
	})

	t.Run("Success And User already exists", func(t *testing.T) {
		username := "admin"
		secrets.FirstAdminUsername = username
		secrets.FirstAdminPassword = "password123"

		// validate admin not exists yet:
		admin, err := query.GetUserByUsername(ctx, username)
		if err != nil && !errors.Is(err, db.ErrNoRows) {
			t.Fatalf("failed to get user by username: %v", err)
		}
		if err == nil {
			t.Fatalf("expected no admin user, but found one: %v", admin)
		}

		// success, should create admin user
		createFirstAdmin(query)

		// validate admin was created:
		admin, err = query.GetUserByUsername(ctx, username)
		if err != nil {
			t.Fatalf("failed to get user by username after creation: %v", err)
		}
		if admin.Username != username {
			t.Errorf("expected username %s, got %s", username, admin.Username)
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
		secrets.FirstAdminUsername = "admin"
		secrets.FirstAdminPassword = "password123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		expectedErr := errors.New("db error")
		q.EXPECT().
			CheckUserExists(gomock.Any(), secrets.FirstAdminUsername).
			Return(int32(0), expectedErr)

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		createFirstAdmin(q)
	})
}
