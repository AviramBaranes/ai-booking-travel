package accounts

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type GetUserEmailRequest struct {
	UserID int32
}

type GetUserEmailResponse struct {
	Email string
}

// encore:api private
func (s *Service) GetUserEmail(ctx context.Context, id int32) (*GetUserEmailResponse, error) {
	user, err := s.query.GetUserById(ctx, id)
	if err != nil {
		if errs.Code(err) == errs.NotFound {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get user email", "error", err, "id", id)
		return nil, api_errors.ErrInternalError
	}
	return &GetUserEmailResponse{Email: user.Email}, nil
}
