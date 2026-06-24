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
		officeID := int64(1)
		organizationID := int64(2)
		isOrganic := true
		data := jwt.AccessTokenData{
			UserID: 123,
			Role:   db.UserRoleAgent,
			OrganizationContext: &jwt.OrganizationContext{
				OfficeID:       officeID,
				OrganizationID: organizationID,
				IsOrganic:      isOrganic,
			},
		}

		token, _, err := jwt.SignAccessToken(data)
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

		if authData.UserID != data.UserID {
			t.Errorf("Expected UserID %d, got %d", data.UserID, authData.UserID)
		}
		if authData.Role != UserRole(data.Role) {
			t.Errorf("Expected Role %s, got %s", UserRole(data.Role), authData.Role)
		}
		if authData.OrganizationContext.OfficeID != officeID {
			t.Errorf("Expected OfficeID %d, got %d", officeID, authData.OrganizationContext.OfficeID)
		}
		if authData.OrganizationContext.OrganizationID != organizationID {
			t.Errorf("Expected OrganizationID %d, got %d", organizationID, authData.OrganizationContext.OrganizationID)
		}
		if authData.OrganizationContext.IsOrganic != isOrganic {
			t.Errorf("Expected IsOrganic %v, got %v", isOrganic, authData.OrganizationContext.IsOrganic)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, _, err := AuthHandler(ctx, "invalid-token")
		api_errors.AssertApiError(t, api_errors.ErrUnauthenticated, err)
	})

	t.Run("Expired token", func(t *testing.T) {
		officeID := int64(456)
		organizationID := int64(789)
		isOrganic := false
		data := jwt.AccessTokenData{
			UserID: 456,
			Role:   db.UserRoleAgent,
			OrganizationContext: &jwt.OrganizationContext{
				OfficeID:       officeID,
				OrganizationID: organizationID,
				IsOrganic:      isOrganic,
			},
		}

		token, _, err := jwt.SignAccessToken(data)
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
