package jwt

import (
	"encore.app/services/accounts/db"
	"github.com/golang-jwt/jwt/v4"
)

// AccessTokenClaims represents the claims for an access token.
type AccessTokenClaims struct {
	Role       db.UserRole `json:"role"`
	UserID     int32       `json:"userId"`
	OfficeID   *int32      `json:"officeId,omitempty"`
	AdminRefID *int32      `json:"adminRefId,omitempty"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents the claims for a refresh token.
type RefreshTokenClaims struct {
	UserID int32 `json:"userId"`
	jwt.RegisteredClaims
}
