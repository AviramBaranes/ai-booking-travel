package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// LoginAsUserParams defines the parameters required to login as a user (agent or customer) on behalf of an admin.
type LoginAsUserParams struct {
	UserID int64 `json:"userId" validate:"required"`
}

func (p LoginAsUserParams) Validate() error {
	return validation.ValidateStruct(p)
}

// LoginAsUser logs in as an agent or customer on behalf of an admin.
// adminID is extracted from the caller's auth context by the thin wrapper.
func (s *AuthService) LoginAsUser(ctx context.Context, params LoginAsUserParams, adminID int64) (*LoginResponse, error) {
	user, err := s.query.GetUserById(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get user by ID", "user_id", params.UserID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if user.Role != db.UserRoleAgent && user.Role != db.UserRoleCustomer {
		return nil, api_errors.ErrUnauthorized
	}

	data := accessTokenDataFromUser(user, &adminID)
	tokens, err := s.generateTokens(ctx, data)

	if err != nil {
		rlog.Error("failed to generate tokens in login as user", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:                   user.ID,
		Role:                 user.Role,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		AccessToken:          tokens.AccessToken,
		AccessTokenExpiresAt: tokens.AccessTokenExpiresAt,
		Email:                user.Email,
		PhoneNumber:          ptrToStr(user.PhoneNumber),
		OfficeID:             user.OfficeID,
		SetCookies:           authCookies(tokens.RefreshToken),
	}, nil
}

// accessTokenDataFromUser converts a db.GetUserByIdRow to jwt.AccessTokenData.
// OrganizationContext is only populated for agents; customers have no office/org.
func accessTokenDataFromUser(agent db.GetUserByIdRow, adminID *int64) jwt.AccessTokenData {
	var orgCtx *jwt.OrganizationContext
	if agent.OfficeID != nil && agent.OrganizationID != nil && agent.IsOrganic != nil {
		orgCtx = &jwt.OrganizationContext{
			OfficeID:       *agent.OfficeID,
			OrganizationID: *agent.OrganizationID,
			IsOrganic:      *agent.IsOrganic,
		}
	}

	return jwt.AccessTokenData{
		UserID:              agent.ID,
		Role:                agent.Role,
		AdminRefID:          adminID,
		OrganizationContext: orgCtx,
	}
}
