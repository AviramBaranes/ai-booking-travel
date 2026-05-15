package accounts

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type GetUserEmailParams struct {
	UserID int64
}

type GetUserEmailResponse struct {
	Email string
}

// encore:api private
func (s *Service) GetUserEmail(ctx context.Context, params GetUserEmailParams) (*GetUserEmailResponse, error) {
	user, err := s.query.GetUserById(ctx, params.UserID)
	if err != nil {
		if errs.Code(err) == errs.NotFound {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get user email", "error", err, "id", params.UserID)
		return nil, api_errors.ErrInternalError
	}
	return &GetUserEmailResponse{Email: user.Email}, nil
}
