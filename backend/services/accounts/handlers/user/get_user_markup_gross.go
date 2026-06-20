package user

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.dev/rlog"
)

type GetUserMarkupGrossResponse struct {
	GrossMarkup float64 `json:"gross_markup"`
}

// GetUserMarkupGross retrieves the gross markup for a user by their ID.
func (s *UserService) GetUserMarkupGross(ctx context.Context, userID int64) (*GetUserMarkupGrossResponse, error) {
	gross, err := s.query.GetUserGrossMarkup(ctx, userID)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotFound) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get user gross markup", "error", err, "userID", userID)
		return nil, api_errors.ErrInternalError
	}

	return &GetUserMarkupGrossResponse{
		GrossMarkup: dbadapters.NumericToFloat64(gross),
	}, nil
}
