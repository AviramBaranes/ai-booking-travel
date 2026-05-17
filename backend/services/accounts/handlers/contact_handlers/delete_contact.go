package contact_handlers

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

func (s *ContactService) DeleteContact(ctx context.Context, id int64) error {
	err := s.query.DeleteContact(ctx, id)
	if err != nil {
		rlog.Error("failed to delete contact", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
