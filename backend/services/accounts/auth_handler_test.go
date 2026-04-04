package accounts

import (
	"context"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	jwtgo "github.com/golang-jwt/jwt/v4"
)

func TestAuthHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid token", func(t *testing.T) {
		office_id := int32(123)
		user := db.User{
			ID:       123,
			Role:     db.UserRoleAgent,
			Email:    "test@test.com",
			OfficeID: &office_id,
		}

		token, err := jwt.SignAccessToken(user)
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		uid, authData, err := AuthHandler(ctx, token)
		if err != nil {
			t.Fatalf("AuthHandler failed: %v", err)
		}

		if string(uid) != "123" {
			t.Errorf("Expected UID '123', got '%s'", uid)
		}

		if authData.UserID != user.ID {
			t.Errorf("Expected UserID %d, got %d", user.ID, authData.UserID)
		}
		if authData.Role != UserRoleAgent {
			t.Errorf("Expected Role %s, got %s", UserRoleAgent, authData.Role)
		}
		if authData.OfficeID != office_id {
			t.Errorf("Expected OfficeID %d, got %d", office_id, authData.OfficeID)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, _, err := AuthHandler(ctx, "invalid-token")
		api_errors.AssertApiError(t, api_errors.ErrUnauthenticated, err)
	})

	t.Run("Expired token", func(t *testing.T) {
		office_id := int32(456)
		user := db.User{
			ID:       456,
			Role:     db.UserRoleAgent,
			Email:    "test@test.com",
			OfficeID: &office_id,
		}

		token, err := jwt.SignAccessToken(user)
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}

		// Fast forward time to simulate expiration
		// AccessTokenTTL is 15 minutes, so we add 16 minutes
		originalTimeFunc := jwtgo.TimeFunc
		defer func() { jwtgo.TimeFunc = originalTimeFunc }()

		jwtgo.TimeFunc = func() time.Time {
			return time.Now().Add(20 * time.Minute)
		}

		_, _, err = AuthHandler(ctx, token)
		api_errors.AssertApiError(t, api_errors.ErrUnauthenticated, err)
	})
}
