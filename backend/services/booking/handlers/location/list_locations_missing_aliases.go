package location

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev"
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
	var search *string
	meta := encore.Meta()
	if meta.Environment.Type == encore.EnvDevelopment && meta.Environment.Cloud == encore.CloudLocal {
		str := "ams"
		search = &str
	}

	rlog.Info("listing locations missing aliases", "search", search)

	locations, err := s.query.ListLocationsWithoutAliases(ctx, search)
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
