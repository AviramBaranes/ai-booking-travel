package auth

import (
	"context"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/services/auth/db"
	"encore.dev/beta/auth"
)

type AuthData struct {
	UserID     int32
	Role       db.UserRole
	Username   string
	OfficeCode string
	AgentCode  string
}

// encore: authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, *AuthData, error) {
	claims, err := jwt.ValidateAccessToken(token)
	if err != nil {
		return "", nil, api_errors.ErrUnauthenticated
	}

	authData := &AuthData{
		UserID:     claims.UserID,
		Role:       claims.Role,
		Username:   claims.Username,
		OfficeCode: claims.OfficeCode,
		AgentCode:  claims.AgentCode,
	}

	uid := strconv.Itoa(int(authData.UserID))
	return auth.UID(uid), authData, nil
}
