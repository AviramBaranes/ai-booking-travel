package jwt

import (
	"encore.app/services/accounts/db"
	"github.com/golang-jwt/jwt/v4"
)

// AccessTokenClaims represents the claims for an access token.
type AccessTokenClaims struct {
	Role       db.UserRole `json:"role"`
	UserID     int64       `json:"userId"`
	OfficeID   *int64      `json:"officeId,omitempty"`
	AdminRefID *int64      `json:"adminRefId,omitempty"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents the claims for a refresh token.
type RefreshTokenClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}
