package user

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type GetUserMarkupGrossResponse struct {
	GrossMarkup float64 `json:"gross_markup"`
}

// GetUserMarkupGross retrieves the gross markup for a user by their ID.
// Returns 0 for users without an associated office (e.g. admins, accountants).
func (s *UserService) GetUserMarkupGross(ctx context.Context, userID int64) (*GetUserMarkupGrossResponse, error) {
	gross, err := s.query.GetUserGrossMarkup(ctx, userID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &GetUserMarkupGrossResponse{GrossMarkup: 0}, nil
		}
		rlog.Error("failed to get user gross markup", "error", err, "userID", userID)
		return nil, api_errors.ErrInternalError
	}

	return &GetUserMarkupGrossResponse{
		GrossMarkup: dbadapters.NumericToFloat64(gross),
	}, nil
}
