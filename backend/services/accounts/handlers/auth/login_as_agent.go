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

// LoginAsAgentParams defines the parameters required to login as an agent.
type LoginAsAgentParams struct {
	AgentID int64 `json:"agentId" validate:"required"`
}

func (p LoginAsAgentParams) Validate() error {
	return validation.ValidateStruct(p)
}

// LoginAsAgent logs in as an agent on behalf of an admin.
// adminID is extracted from the caller's auth context by the thin wrapper.
func (s *AuthService) LoginAsAgent(ctx context.Context, params LoginAsAgentParams, adminID int64) (*LoginResponse, error) {
	agent, err := s.query.GetUserById(ctx, params.AgentID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get agent by ID", "agent_id", params.AgentID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	data := accessTokenDataFromUser(agent, &adminID)
	accessToken, refreshToken, err := s.generateTokens(ctx, data)

	if err != nil {
		rlog.Error("failed to generate tokens in login as agent", "user_id", agent.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:           agent.ID,
		Role:         agent.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        agent.Email,
		PhoneNumber:  ptrToStr(agent.PhoneNumber),
		OfficeID:     agent.OfficeID,
	}, nil
}

// accessTokenDataFromUser converts a db.GetUserByIdRow to jwt.AccessTokenData, including the admin ref ID for login as agent.
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
