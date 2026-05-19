package jwt

import (
	"encore.app/services/accounts/db"
	"github.com/golang-jwt/jwt/v4"
)

// AccessTokenClaims represents the claims for an access token.
type AccessTokenClaims struct {
	Role                db.UserRole          `json:"role"`
	UserID              int64                `json:"userId"`
	AdminRefID          *int64               `json:"adminRefId,omitempty"`
	OrganizationContext *OrganizationContext `json:"organizationContext,omitempty"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents the claims for a refresh token.
type RefreshTokenClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}
