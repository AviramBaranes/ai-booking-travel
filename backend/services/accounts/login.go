package accounts

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
	"encore.dev/pubsub"
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

// SendCustomerLoginOTPParams defines the parameters required to send a login OTP to a customer.
type SendCustomerLoginOTPParams struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,israeli_phone"`
}

func (p SendCustomerLoginOTPParams) Validate() error {
	return validation.ValidateStruct(p)
}

// encore:api public method=POST path=/customer-login/send-otp
func (s *Service) SendCustomerLoginOTP(ctx context.Context, params SendCustomerLoginOTPParams) error {
	user, err := s.query.GetUserByPhone(ctx, &params.PhoneNumber)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return ErrInvalidCredentials
		}
		rlog.Error("failed to get user by phone number", "phone_number", params.PhoneNumber, "error", err)
		return api_errors.ErrInternalError
	}

	if user.Role != db.UserRoleCustomer {
		return ErrInvalidCredentials
	}

	otp, err := generateOTP(6)
	if err != nil {
		rlog.Error("generating otp failed", "error", err)
		return api_errors.ErrInternalError
	}

	err = s.query.SaveOTP(ctx, db.SaveOTPParams{
		ID:  user.ID,
		Otp: &otp,
	})

	if err != nil {
		rlog.Error("failed to save otp", "user_id", user.ID, "error", err)
		return api_errors.ErrInternalError
	}

	if _, err := CustomerLoginOTPRequestedTopic.Publish(ctx, &CustomerLoginOTPRequestedEvent{
		PhoneNumber: params.PhoneNumber,
		OTP:         otp,
	}); err != nil {
		rlog.Error("failed to publish customer login OTP requested event", "phone_number", params.PhoneNumber, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}

type CustomerLoginOTPRequestedEvent struct {
	PhoneNumber string `json:"phoneNumber"`
	OTP         string `json:"otp"`
}

var CustomerLoginOTPRequestedTopic = pubsub.NewTopic[*CustomerLoginOTPRequestedEvent](
	"customer-login-otp-requested",
	pubsub.TopicConfig{
		DeliveryGuarantee: pubsub.AtLeastOnce,
	},
)

func generateOTP(length int) (string, error) {
	charset := "0123456789"

	otp := make([]byte, length)
	randMax := int64(len(charset))
	for i := range otp {
		index, err := rand.Int(rand.Reader, big.NewInt(randMax))
		if err != nil {
			return "", err
		}

		otp[i] = charset[index.Int64()]
	}

	return string(otp), nil
}

type ValidateCustomerLoginOTPParams struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,israeli_phone"`
	OTP         string `json:"otp" validate:"required,len=6"`
}

func (p ValidateCustomerLoginOTPParams) Validate() error {
	return validation.ValidateStruct(p)
}

// encore:api public method=POST path=/customer-login/validate-otp
func (s *Service) ValidateCustomerLoginOTP(ctx context.Context, params ValidateCustomerLoginOTPParams) (*LoginResponse, error) {
	user, err := s.query.GetUserByPhone(ctx, &params.PhoneNumber)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get user by phone number", "phone_number", params.PhoneNumber, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if user.Otp == nil || *user.Otp != params.OTP {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user, nil)
	if err != nil {
		rlog.Error("failed to generate tokens in validate customer login OTP", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	// Clear OTP after successful login
	err = s.query.SaveOTP(ctx, db.SaveOTPParams{
		ID:  user.ID,
		Otp: nil,
	})
	if err != nil {
		rlog.Error("failed to clear OTP after successful login", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:           user.ID,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
		PhoneNumber:  ptrToStr(user.PhoneNumber),
	}, nil
}

func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
