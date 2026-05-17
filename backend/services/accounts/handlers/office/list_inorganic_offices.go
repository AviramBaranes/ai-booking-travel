package office

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type InorganicOffice struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ListInorganicOfficeResponse struct {
	Offices []InorganicOffice `json:"offices"`
}

// ListInorganicOffices lists all inorganic offices for accountant use.
func (s *OfficeService) ListInorganicOffices(ctx context.Context) (*ListInorganicOfficeResponse, error) {
	rows, err := s.query.ListInorganicOffices(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &ListInorganicOfficeResponse{Offices: []InorganicOffice{}}, nil
		}
		rlog.Error("failed to list inorganic offices", "error", err)
		return nil, api_errors.ErrInternalError
	}

	offices := make([]InorganicOffice, 0, len(rows))
	for _, r := range rows {
		offices = append(offices, InorganicOffice{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &ListInorganicOfficeResponse{Offices: offices}, nil
}
