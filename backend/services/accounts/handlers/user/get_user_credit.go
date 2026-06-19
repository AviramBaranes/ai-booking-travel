package user

import (
	"context"

	dbadapters "encore.app/internal/db_adapters"
	"encore.dev/rlog"
)

// GetUserCreditParams are the params for getting a user's credit.
type GetUserCreditResponse struct {
	Obligo     float64 `json:"obligo"`
	BalanceDue float64 `json:"balance_due"`
}

func (s *UserService) GetUserCredit(ctx context.Context, userID int64) (*GetUserCreditResponse, error) {
	res, err := s.query.GetUserCredit(ctx, userID)
	if err != nil {
		rlog.Error("failed to get user credit", "error", err, "userID", userID)
		return nil, err
	}

	return &GetUserCreditResponse{
		Obligo:     float64(res.Obligo),
		BalanceDue: dbadapters.NumericToFloat64(res.BalanceDue),
	}, nil
}
