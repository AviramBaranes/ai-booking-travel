package user

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// GetUserEmailParams are the params for retrieving a user's email.
type GetUserEmailParams struct {
	UserID int64
}

// GetUserEmailResponse is the response for retrieving a user's email.
type GetUserEmailResponse struct {
	Email string
}

// GetUserEmail retrieves the email address of a user by ID.
func (s *UserService) GetUserEmail(ctx context.Context, params GetUserEmailParams) (*GetUserEmailResponse, error) {
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
