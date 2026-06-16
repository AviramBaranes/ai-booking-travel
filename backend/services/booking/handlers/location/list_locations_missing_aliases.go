package location

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type ListLocationsMissingAliasesResponse struct {
	Locations []MissingAliasLocation `json:"locations"`
}

type MissingAliasLocation struct {
	ID   int64   `json:"id"`
	Iata *string `json:"iata"`
	Name string  `json:"name"`
}

func (s *LocationService) ListLocationsMissingAliases(ctx context.Context) (*ListLocationsMissingAliasesResponse, error) {
	locations, err := s.query.ListLocationsWithoutAliases(ctx)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to list locations without aliases", "error", err)
		return nil, api_errors.ErrInternalError
	}

	mLocs := make([]MissingAliasLocation, 0, len(locations))
	for _, l := range locations {
		mLocs = append(mLocs, MissingAliasLocation{
			ID:   l.ID,
			Name: l.Name,
			Iata: l.Iata,
		})
	}

	return &ListLocationsMissingAliasesResponse{
		Locations: mLocs,
	}, nil
}
