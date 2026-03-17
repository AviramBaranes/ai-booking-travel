package auth

import (
	"context"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/services/auth/db"
	jwtgo "github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestAuthHandler(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid token", func(t *testing.T) {
		user := db.User{
			ID:         123,
			Role:       db.UserRoleAdmin,
			Username:   "testuser",
			OfficeCode: pgtype.Text{String: "office1", Valid: true},
			AgentCode:  pgtype.Text{String: "agent1", Valid: true},
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
		if authData.Role != user.Role {
			t.Errorf("Expected Role %s, got %s", user.Role, authData.Role)
		}
		if authData.Username != user.Username {
			t.Errorf("Expected Username %s, got %s", user.Username, authData.Username)
		}
		if authData.OfficeCode != user.OfficeCode.String {
			t.Errorf("Expected OfficeCode %s, got %s", user.OfficeCode.String, authData.OfficeCode)
		}
		if authData.AgentCode != user.AgentCode.String {
			t.Errorf("Expected AgentCode %s, got %s", user.AgentCode.String, authData.AgentCode)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, _, err := AuthHandler(ctx, "invalid-token")
		api_errors.AssertApiError(t, api_errors.ErrUnauthenticated, err)
	})

	t.Run("Expired token", func(t *testing.T) {
		user := db.User{
			ID:       456,
			Role:     db.UserRoleAgent,
			Username: "expireduser",
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
