package accounts

import (
	"context"

	auth "encore.app/services/accounts/handlers/auth"
)

// encore:api public path=/password-reset method=POST
func (s *Service) RequestPasswordReset(ctx context.Context, params auth.SendPasswordResetTokenParams) error {
	h := auth.NewAuthService(s.query)
	return h.SendPasswordResetToken(ctx, params)
}
