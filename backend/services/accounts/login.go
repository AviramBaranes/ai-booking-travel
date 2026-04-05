package accounts

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
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

// LoginResponse defines the response structure for user login.
type LoginResponse struct {
	ID           int32       `json:"id"`
	Email        string      `json:"email,omitempty"`
	Role         db.UserRole `json:"role,omitempty"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	PhoneNumber  string      `json:"phoneNumber,omitempty"`
	OfficeID     *int32      `json:"officeId,omitempty"`
}

// encore:api public path=/login method=POST
func (s *Service) Login(ctx context.Context, p LoginParams) (*LoginResponse, error) {
	row, err := s.query.GetUserByEmail(ctx, p.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get user by email", "email", p.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if !password.ComparePassword(row.PasswordHash, p.Password) {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, row, nil)
	if err != nil {
		rlog.Error("failed to generate tokens", "user_id", row.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:           row.ID,
		Role:         row.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        row.Email,
		PhoneNumber:  ptrToStr(row.PhoneNumber),
		OfficeID:     row.OfficeID,
	}, nil
}

func (s *Service) generateTokens(ctx context.Context, user db.User, adminRefID *int32) (string, string, error) {
	accessToken, err := jwt.SignAccessToken(user, adminRefID)
	if err != nil {
		return "", "", errs.Wrap(err, "failed to sign access token")
	}

	refreshToken, jti, exp, err := jwt.SignRefreshToken(user.ID)
	if err != nil {
		return "", "", errs.Wrap(err, "failed to sign refresh token")
	}

	err = s.query.SaveRefreshToken(ctx, db.SaveRefreshTokenParams{
		Jti:        jti,
		UserID:     user.ID,
		AdminRefID: adminRefID,
		ExpiresAt:  db.DBTime(exp),
	})
	if err != nil {
		return "", "", errs.Wrap(err, "failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

type LoginAsAgentParams struct {
	AgentID int32 `json:"agentId" validate:"required"`
}

func (p LoginAsAgentParams) Validate() error {
	return validation.ValidateStruct(p)
}

// encore:api auth method=POST path=/login/as-agent tag:admin
func (s *Service) LoginAsAgent(ctx context.Context, params LoginAsAgentParams) (*LoginResponse, error) {
	agent, err := s.query.GetUserById(ctx, params.AgentID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get agent by ID", "agent_id", params.AgentID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	authData := GetAuthData()
	accessToken, refreshToken, err := s.generateTokens(ctx, agent, &authData.UserID)
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

// encore:api public method=POST path=/login/back-to-admin tag:agent
func (s *Service) LoginBackToAdmin(ctx context.Context) (*LoginResponse, error) {
	authData := GetAuthData()

	var adminID int32
	if authData.AdminRefID == nil {
		return nil, ErrInvalidCredentials
	}

	adminID = *authData.AdminRefID

	admin, err := s.query.GetUserById(ctx, adminID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get admin by ID in login back to admin", "admin_id", adminID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, admin, nil)
	if err != nil {
		rlog.Error("failed to generate tokens in login back to admin", "user_id", admin.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:           admin.ID,
		Role:         admin.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        admin.Email,
	}, nil
}

func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
