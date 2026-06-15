package location

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type InsertLocationAliasesParams struct {
	LocationID int64   `json:"locationId"`
	Lang       string  `json:"lang"`
	Aliases    []Alias `json:"aliases"`
}

type Alias struct {
	Value     string `json:"value"`
	AliasType string `json:"type"`
}

func (s *LocationService) InsertLocationAliases(ctx context.Context, p InsertLocationAliasesParams) error {
	var dbParams db.InsertManyLocationAliasesParams
	for _, alias := range p.Aliases {
		dbParams.LocationIds = append(dbParams.LocationIds, p.LocationID)
		dbParams.Aliases = append(dbParams.Aliases, alias.Value)
	}

	err := s.query.InsertManyLocationAliases(ctx, dbParams)
	if err != nil {
		rlog.Error("failed to insert location aliases", "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
