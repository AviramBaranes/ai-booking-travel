package accounts

import (
	"context"

	auth_handler "encore.app/services/accounts/handlers/auth_handler"
)

// encore:api public path=/login method=POST
func (s *Service) Login(ctx context.Context, p auth_handler.LoginParams) (*auth_handler.LoginResponse, error) {
	h := auth_handler.NewAuthService(s.query)
	return h.Login(ctx, p)
}

// encore:api auth method=POST path=/login/as-agent tag:admin
func (s *Service) LoginAsAgent(ctx context.Context, params auth_handler.LoginAsAgentParams) (*auth_handler.LoginResponse, error) {
	authData := GetAuthData()
	h := auth_handler.NewAuthService(s.query)
	return h.LoginAsAgent(ctx, params, authData.UserID)
}

// encore:api public method=POST path=/login/back-to-admin tag:agent
func (s *Service) LoginBackToAdmin(ctx context.Context) (*auth_handler.LoginResponse, error) {
	authData := GetAuthData()
	h := auth_handler.NewAuthService(s.query)
	return h.LoginBackToAdmin(ctx, authData.AdminRefID)
}

// encore:api public method=POST path=/customer-login/send-otp
func (s *Service) SendCustomerLoginOTP(ctx context.Context, params auth_handler.SendCustomerLoginOTPParams) error {
	h := auth_handler.NewAuthService(s.query)
	return h.SendCustomerLoginOTP(ctx, params)
}

// encore:api public method=POST path=/customer-login/validate-otp
func (s *Service) ValidateCustomerLoginOTP(ctx context.Context, params auth_handler.ValidateCustomerLoginOTPParams) (*auth_handler.LoginResponse, error) {
	h := auth_handler.NewAuthService(s.query)
	return h.ValidateCustomerLoginOTP(ctx, params)
}

// encore:api public method=POST path=/refresh
func (s *Service) RefreshTokens(ctx context.Context, p auth_handler.RefreshTokensParams) (*auth_handler.LoginResponse, error) {
	h := auth_handler.NewAuthService(s.query)
	return h.RefreshTokens(ctx, p)
}
