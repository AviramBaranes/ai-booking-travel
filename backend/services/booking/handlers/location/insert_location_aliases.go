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
		dbParams.Names = append(dbParams.Names, alias.Value)
		dbParams.Types = append(dbParams.Types, string(alias.AliasType))
		aliasType := db.LocationAliasType(alias.AliasType)
		if aliasType == db.LocationAliasTypeTranslation {
			dbParams.LanguageCodes = append(dbParams.LanguageCodes, p.Lang)
		} else {
			dbParams.LanguageCodes = append(dbParams.LanguageCodes, "")
		}
	}

	err := s.query.InsertManyLocationAliases(ctx, dbParams)
	if err != nil {
		rlog.Error("failed to insert location aliases", "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
