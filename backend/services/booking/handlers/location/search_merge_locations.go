package location

import (
	"context"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.dev/rlog"
)

type SearchMergeLocationsParams struct {
	Search string `query:"search" validate:"required,min=2"`
}

func (p SearchMergeLocationsParams) Validate() error {
	return validation.ValidateStruct(p)
}

type SearchMergeLocationsResponse struct {
	Locations []MergeLocation `json:"locations"`
}

// SearchMergeLocations finds locations by id, name, city or IATA, including ones whose codes are all disabled.
func (s *LocationService) SearchMergeLocations(ctx context.Context, p SearchMergeLocationsParams) (*SearchMergeLocationsResponse, error) {
	rows, err := s.query.SearchLocationsForMerge(ctx, strings.TrimSpace(p.Search))
	if err != nil {
		rlog.Error("failed to search locations for merge", "error", err, "search", p.Search)
		return nil, api_errors.ErrInternalError
	}

	locations, err := withMergeDetails(ctx, s.query, rows)
	if err != nil {
		rlog.Error("failed to load merge location details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &SearchMergeLocationsResponse{Locations: locations}, nil
}
