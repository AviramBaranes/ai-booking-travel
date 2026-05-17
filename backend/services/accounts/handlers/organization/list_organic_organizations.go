package organization

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type OrganicOrganization struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ListOrganicOrganizationResponse struct {
	Organizations []OrganicOrganization `json:"organizations"`
}

// ListOrganicOrganizations lists all organic organizations for accountant use.
func (s *OrganizationService) ListOrganicOrganizations(ctx context.Context) (*ListOrganicOrganizationResponse, error) {
	rows, err := s.query.ListOrganicOrganizations(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &ListOrganicOrganizationResponse{Organizations: []OrganicOrganization{}}, nil
		}
		rlog.Error("failed to list organic organizations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	orgs := make([]OrganicOrganization, 0, len(rows))
	for _, r := range rows {
		orgs = append(orgs, OrganicOrganization{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &ListOrganicOrganizationResponse{Organizations: orgs}, nil
}
