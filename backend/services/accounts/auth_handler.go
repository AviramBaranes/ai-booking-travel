package accounts

import (
	"context"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.dev/beta/auth"
)

type UserRole string

const (
	UserRoleAdmin      UserRole = "admin"
	UserRoleAgent      UserRole = "agent"
	UserRoleCustomer   UserRole = "customer"
	UserRoleAccountant UserRole = "accountant"
)

type AuthData struct {
	UserID     int32
	Role       UserRole
	OfficeID   int32
	AdminRefID *int32
}

// encore: authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, *AuthData, error) {
	claims, err := jwt.ValidateAccessToken(token)
	if err != nil {
		return "", nil, api_errors.ErrUnauthenticated
	}

	var officeID int32
	if claims.OfficeID != nil {
		officeID = *claims.OfficeID
	}
	authData := &AuthData{
		UserID:     claims.UserID,
		Role:       UserRole(claims.Role),
		OfficeID:   officeID,
		AdminRefID: claims.AdminRefID,
	}

	uid := strconv.Itoa(int(authData.UserID))
	return auth.UID(uid), authData, nil
}

// GetAuthData is a helper function to retrieve the authentication data of the currently authenticated user from the context.
func GetAuthData() *AuthData {
	authData, ok := auth.Data().(*AuthData)
	if !ok {
		return nil
	}

	return authData
}
