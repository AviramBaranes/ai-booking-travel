package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// LoginParams defines the parameters required for user login.
type LoginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required" encore:"sensitive"`
}

func (p LoginParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *AuthService) Login(ctx context.Context, p LoginParams) (*LoginResponse, error) {
	user, err := s.query.GetUserByEmail(ctx, p.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get user by email", "email", p.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if !password.ComparePassword(user.PasswordHash, p.Password) {
		return nil, ErrInvalidCredentials
	}

	var orgCtx *jwt.OrganizationContext
	if user.OfficeID != nil && user.OrganizationID != nil && user.IsOrganic != nil {
		orgCtx = &jwt.OrganizationContext{
			OfficeID:       *user.OfficeID,
			OrganizationID: *user.OrganizationID,
			IsOrganic:      *user.IsOrganic,
		}
	}

	tokens, err := s.generateTokens(ctx, jwt.AccessTokenData{
		UserID:              user.ID,
		Role:                user.Role,
		OrganizationContext: orgCtx,
	})
	if err != nil {
		rlog.Error("failed to generate tokens", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:                   user.ID,
		Role:                 user.Role,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		AccessToken:          tokens.AccessToken,
		RefreshToken:         tokens.RefreshToken,
		AccessTokenExpiresAt: tokens.AccessTokenExpiresAt,
		Email:                user.Email,
		PhoneNumber:          ptrToStr(user.PhoneNumber),
		OfficeID:             user.OfficeID,
	}, nil
}
