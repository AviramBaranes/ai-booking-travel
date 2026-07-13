package accounts

import (
	"context"

	auth "encore.app/services/accounts/handlers/auth"
)

// encore:api public path=/login method=POST
func (s *Service) Login(ctx context.Context, p auth.LoginParams) (*auth.LoginResponse, error) {
	h := auth.NewAuthService(s.query)
	return h.Login(ctx, p)
}

// encore:api auth method=POST path=/login/as-user tag:admin
func (s *Service) LoginAsUser(ctx context.Context, params auth.LoginAsUserParams) (*auth.LoginResponse, error) {
	authData := GetAuthData()
	h := auth.NewAuthService(s.query)
	return h.LoginAsUser(ctx, params, authData.UserID)
}

// encore:api public method=POST path=/login/back-to-admin tag:agent_customer
func (s *Service) LoginBackToAdmin(ctx context.Context) (*auth.LoginResponse, error) {
	authData := GetAuthData()
	h := auth.NewAuthService(s.query)
	return h.LoginBackToAdmin(ctx, authData.AdminRefID)
}

// encore:api public method=POST path=/customer-login/send-otp
func (s *Service) SendCustomerLoginOTP(ctx context.Context, params auth.SendCustomerLoginOTPParams) error {
	h := auth.NewAuthService(s.query)
	return h.SendCustomerLoginOTP(ctx, params)
}

// encore:api public method=POST path=/customer-login/validate-otp
func (s *Service) ValidateCustomerLoginOTP(ctx context.Context, params auth.ValidateCustomerLoginOTPParams) (*auth.LoginResponse, error) {
	h := auth.NewAuthService(s.query)
	return h.ValidateCustomerLoginOTP(ctx, params)
}

// encore:api public method=POST path=/refresh
func (s *Service) RefreshTokens(ctx context.Context, p auth.RefreshTokensParams) (*auth.LoginResponse, error) {
	h := auth.NewAuthService(s.query)
	return h.RefreshTokens(ctx, p)
}

// encore:api private
func (s *Service) GetCustomerToken(ctx context.Context, p auth.GetCustomerTokenParams) (*auth.LoginResponse, error) {
	h := auth.NewAuthService(s.query)
	return h.GetCustomerToken(ctx, p)
}

// encore:api public method=POST path=/auth/session
func (s *Service) ResolveSession(ctx context.Context, p auth.ResolveSessionParams) (*auth.SessionUser, error) {
	h := auth.NewAuthService(s.query)
	return h.ResolveSession(ctx, p)
}

// encore:api public method=POST path=/logout
func (s *Service) Logout(ctx context.Context) (*auth.LogoutResponse, error) {
	h := auth.NewAuthService(s.query)
	return h.Logout(ctx)
}
