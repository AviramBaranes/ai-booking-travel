package user

import (
	"context"

	dbadapters "encore.app/internal/db_adapters"
	"encore.dev/rlog"
)

// GetUserCreditByIDParams names the agent whose credit to get.
type GetUserCreditByIDParams struct {
	UserID int64 `json:"user_id"`
}

// GetUserCreditParams are the params for getting a user's credit.
// OfficeID, OrganizationID and IsOrganic identify the billing entity the credit belongs to:
// the organization when it is organic, the office otherwise.
type GetUserCreditResponse struct {
	Obligo         float64 `json:"obligo"`
	BalanceDue     float64 `json:"balance_due"`
	OfficeID       int64   `json:"office_id"`
	OrganizationID int64   `json:"organization_id"`
	IsOrganic      bool    `json:"is_organic"`
}

func (s *UserService) GetUserCredit(ctx context.Context, userID int64) (*GetUserCreditResponse, error) {
	res, err := s.query.GetUserCredit(ctx, userID)
	if err != nil {
		rlog.Error("failed to get user credit", "error", err, "userID", userID)
		return nil, err
	}

	return &GetUserCreditResponse{
		Obligo:         float64(res.Obligo),
		BalanceDue:     dbadapters.NumericToFloat64(res.BalanceDue),
		OfficeID:       res.OfficeID,
		OrganizationID: res.OrganizationID,
		IsOrganic:      res.IsOrganic,
	}, nil
}
