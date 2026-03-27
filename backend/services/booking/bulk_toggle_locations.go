package booking

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.dev/rlog"
)

type BulkToggleLocationsRequest struct {
	IDs     []int64 `json:"ids" validate:"required,min=1,dive,min=1"`
	Enabled bool    `json:"enabled"`
}

func (p BulkToggleLocationsRequest) Validate() error {
	return validation.ValidateStruct(p)
}

//encore:api auth method=PATCH path=/location-bulk-toggle tag:admin
func (s *Service) BulkToggleLocations(ctx context.Context, p *BulkToggleLocationsRequest) error {
	for _, id := range p.IDs {
		var err error
		if p.Enabled {
			err = s.query.EnableLocationBrokerCode(ctx, id)
		} else {
			err = s.query.DisableLocationBrokerCode(ctx, id)
		}
		if err != nil {
			rlog.Error("failed to toggle location broker code", "error", err, "id", id, "enabled", p.Enabled)
			return api_errors.ErrInternalError
		}
	}
	return nil
}
