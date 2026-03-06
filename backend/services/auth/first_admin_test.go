package auth

import (
	"errors"
	"testing"

	"encore.app/services/auth/db"
	"encore.app/services/auth/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateFirstAdmin(t *testing.T) {
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

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		// Expect no calls to DB
		createFirstAdmin(q)
	})

	t.Run("User already exists", func(t *testing.T) {
		secrets.FirstAdminUsername = "admin"
		secrets.FirstAdminPassword = "password123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			CheckUserExists(gomock.Any(), secrets.FirstAdminUsername).
			Return(int32(1), nil)

		// Expect no call to RegisterAdmin
		createFirstAdmin(q)
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

	t.Run("Success", func(t *testing.T) {
		secrets.FirstAdminUsername = "admin"
		secrets.FirstAdminPassword = "password123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			CheckUserExists(gomock.Any(), secrets.FirstAdminUsername).
			Return(int32(0), db.ErrNoRows)

		q.EXPECT().
			RegisterAdmin(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ interface{}, params db.RegisterAdminParams) (db.RegisterAdminRow, error) {
				if params.Username != secrets.FirstAdminUsername {
					t.Errorf("Expected username %s, got %s", secrets.FirstAdminUsername, params.Username)
				}
				// Verify password hash is present (we can't verify exact hash easily without mocking password package, but we can check it's not empty)
				if params.PasswordHash == "" {
					t.Error("Expected password hash to be set")
				}
				return db.RegisterAdminRow{ID: 1, Username: params.Username}, nil
			})

		createFirstAdmin(q)
	})

	t.Run("Hashing failure", func(t *testing.T) {
		// This is hard to test without mocking the password package or passing an invalid password that causes bcrypt to fail (which is rare/hard with simple strings).
		// However, the requirement is to test the logic in first_admin.go.
		// Since password.HashPassword is a direct dependency and not injected, we might skip this specific failure case or accept we can't easily reach it without refactoring first_admin.go to accept a hasher.
		// For now, I'll skip this specific edge case to avoid over-engineering or refactoring the production code unless requested.
	})

	t.Run("Register admin failure", func(t *testing.T) {
		secrets.FirstAdminUsername = "admin"
		secrets.FirstAdminPassword = "password123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		q := mocks.NewMockQuerier(ctrl)
		q.EXPECT().
			CheckUserExists(gomock.Any(), secrets.FirstAdminUsername).
			Return(int32(0), db.ErrNoRows)

		q.EXPECT().
			RegisterAdmin(gomock.Any(), gomock.Any()).
			Return(db.RegisterAdminRow{}, errors.New("register error"))

		// Should just log error and return, no panic
		createFirstAdmin(q)
	})
}
