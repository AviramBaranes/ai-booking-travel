package accounts

import (
	"context"

	auth "encore.app/services/accounts/handlers/auth"
)

// encore:api public path=/send-password-reset method=POST
func (s *Service) RequestPasswordReset(ctx context.Context, params auth.SendPasswordResetTokenParams) error {
	h := auth.NewAuthService(s.query)
	return h.SendPasswordResetToken(ctx, params)
}

// encore:api public path=/reset-password method=POST
func (s *Service) ResetPassword(ctx context.Context, params auth.ResetPasswordParams) error {
	h := auth.NewAuthService(s.query)
	return h.ResetPassword(ctx, params)
}
