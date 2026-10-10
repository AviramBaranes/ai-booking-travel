package location

import (
	"context"
	"errors"
	"strings"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	errLocationsShareBroker = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "both locations have a code of the same broker", api_errors.ErrorDetails{
		Code: api_errors.CodeLocationsShareBroker,
	})
	errLocationMergeConflict = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "another location already has this name or iata", api_errors.ErrorDetails{
		Code: api_errors.CodeLocationMergeConflict,
	})
)

// MergeLocation is a location with everything a merge carries over: its broker codes and aliases.
type MergeLocation struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Country     string                `json:"country"`
	CountryCode string                `json:"country_code"`
	City        *string               `json:"city"`
	Iata        *string               `json:"iata"`
	IsAirport   bool                  `json:"is_airport"`
	BrokerCodes []MergeLocationBroker `json:"broker_codes"`
	Aliases     []string              `json:"aliases"`
}

type MergeLocationBroker struct {
	ID               int64  `json:"id"`
	Broker           string `json:"broker"`
	BrokerLocationID string `json:"broker_location_id"`
	Enabled          bool   `json:"enabled"`
}

type MergeLocationsParams struct {
	KeepID      int64  `json:"keep_id" validate:"required,nefield=RemoveID"`
	RemoveID    int64  `json:"remove_id" validate:"required"`
	Name        string `json:"name" validate:"required,notblank"`
	Country     string `json:"country" validate:"required,notblank"`
	CountryCode string `json:"country_code" validate:"required,notblank"`
	City        string `json:"city"`
	Iata        string `json:"iata" validate:"omitempty,len=3"`
	IsAirport   bool   `json:"is_airport"`
}

func (p MergeLocationsParams) Validate() error {
	return validation.ValidateStruct(p)
}

type MergeLocationsResponse struct {
	Location MergeLocation `json:"location"`
}

// MergeLocations folds the removed location into the kept one: broker codes, aliases and price
// offers move over, the removed location is deleted, and the kept one takes the chosen fields.
// Both original names stay searchable as aliases.
func (s *LocationService) MergeLocations(ctx context.Context, inTx db.TxRunner, p MergeLocationsParams) (*MergeLocationsResponse, error) {
	fields := db.UpdateLocationFieldsParams{
		ID:          p.KeepID,
		Name:        strings.TrimSpace(p.Name),
		Country:     strings.TrimSpace(p.Country),
		CountryCode: strings.ToUpper(strings.TrimSpace(p.CountryCode)),
		City:        nilIfBlank(p.City),
		Iata:        nilIfBlank(strings.ToUpper(p.Iata)),
		IsAirport:   p.IsAirport,
	}

	err := inTx(ctx, func(q db.Querier) error {
		ids := []int64{p.KeepID, p.RemoveID}
		locs, err := q.LockLocationsForMerge(ctx, ids)
		if err != nil {
			rlog.Error("failed to lock locations for merge", "error", err, "ids", ids)
			return api_errors.ErrInternalError
		}
		if len(locs) != 2 {
			return api_errors.ErrNotFound
		}

		codes, err := q.GetAllLocationBrokerCodesByLocationIDs(ctx, ids)
		if err != nil {
			rlog.Error("failed to load broker codes for merge", "error", err, "ids", ids)
			return api_errors.ErrInternalError
		}
		brokersOf := make(map[int64]map[db.Broker]bool)
		for _, c := range codes {
			if brokersOf[c.LocationID] == nil {
				brokersOf[c.LocationID] = make(map[db.Broker]bool)
			}
			brokersOf[c.LocationID][c.Broker] = true
		}
		for b := range brokersOf[p.RemoveID] {
			if brokersOf[p.KeepID][b] {
				return errLocationsShareBroker
			}
		}

		move := func(what string, fn func() error) error {
			if err := fn(); err != nil {
				rlog.Error("failed to move "+what+" for merge", "error", err, "keep_id", p.KeepID, "remove_id", p.RemoveID)
				return api_errors.ErrInternalError
			}
			return nil
		}
		if err := move("broker codes", func() error {
			return q.MoveLocationBrokerCodes(ctx, db.MoveLocationBrokerCodesParams{ToID: p.KeepID, FromID: p.RemoveID})
		}); err != nil {
			return err
		}
		if err := move("aliases", func() error {
			return q.MoveLocationAliases(ctx, db.MoveLocationAliasesParams{ToID: p.KeepID, FromID: p.RemoveID})
		}); err != nil {
			return err
		}
		if err := move("price offers", func() error {
			return q.MovePriceOfferLocations(ctx, db.MovePriceOfferLocationsParams{ToID: p.KeepID, FromID: p.RemoveID})
		}); err != nil {
			return err
		}

		// The removed location goes first, so the kept one can take its name or IATA without
		// tripping the unique indexes.
		if err := q.DeleteLocationByID(ctx, p.RemoveID); err != nil {
			rlog.Error("failed to delete merged location", "error", err, "remove_id", p.RemoveID)
			return api_errors.ErrInternalError
		}

		if _, err := q.UpdateLocationFields(ctx, fields); err != nil {
			if db.IsUniqueViolation(err) {
				return errLocationMergeConflict
			}
			rlog.Error("failed to update merged location", "error", err, "keep_id", p.KeepID)
			return api_errors.ErrInternalError
		}

		var aliases db.InsertManyLocationAliasesParams
		for _, l := range locs {
			if !strings.EqualFold(strings.TrimSpace(l.Name), fields.Name) {
				aliases.LocationIds = append(aliases.LocationIds, p.KeepID)
				aliases.Aliases = append(aliases.Aliases, l.Name)
			}
		}
		if len(aliases.Aliases) > 0 {
			if err := q.InsertManyLocationAliases(ctx, aliases); err != nil {
				rlog.Error("failed to keep merged names as aliases", "error", err, "keep_id", p.KeepID)
				return api_errors.ErrInternalError
			}
		}

		return nil
	})
	if err != nil {
		var apiErr *errs.Error
		if errors.As(err, &apiErr) {
			return nil, err
		}
		rlog.Error("failed to merge locations", "error", err, "keep_id", p.KeepID, "remove_id", p.RemoveID)
		return nil, api_errors.ErrInternalError
	}

	merged, err := loadMergeLocations(ctx, s.query, []int64{p.KeepID})
	if err != nil || len(merged) != 1 {
		rlog.Error("failed to load merged location", "error", err, "keep_id", p.KeepID)
		return nil, api_errors.ErrInternalError
	}

	return &MergeLocationsResponse{Location: merged[0]}, nil
}

func loadMergeLocations(ctx context.Context, q db.Querier, ids []int64) ([]MergeLocation, error) {
	if len(ids) == 0 {
		return []MergeLocation{}, nil
	}
	rows, err := q.GetLocationsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return withMergeDetails(ctx, q, rows)
}

func withMergeDetails(ctx context.Context, q db.Querier, rows []db.Location) ([]MergeLocation, error) {
	locations := make([]MergeLocation, 0, len(rows))
	if len(rows) == 0 {
		return locations, nil
	}

	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}

	codes, err := q.GetAllLocationBrokerCodesByLocationIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	aliases, err := q.ListLocationAliasesByLocationIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	codesOf := make(map[int64][]MergeLocationBroker)
	for _, c := range codes {
		codesOf[c.LocationID] = append(codesOf[c.LocationID], MergeLocationBroker{
			ID:               c.ID,
			Broker:           string(c.Broker),
			BrokerLocationID: c.BrokerLocationID,
			Enabled:          c.Enabled,
		})
	}
	aliasesOf := make(map[int64][]string)
	for _, a := range aliases {
		aliasesOf[a.LocationID] = append(aliasesOf[a.LocationID], a.Alias)
	}

	for _, r := range rows {
		l := MergeLocation{
			ID:          r.ID,
			Name:        r.Name,
			Country:     r.Country,
			CountryCode: r.CountryCode,
			City:        r.City,
			Iata:        r.Iata,
			IsAirport:   r.IsAirport,
			BrokerCodes: codesOf[r.ID],
			Aliases:     aliasesOf[r.ID],
		}
		if l.BrokerCodes == nil {
			l.BrokerCodes = []MergeLocationBroker{}
		}
		if l.Aliases == nil {
			l.Aliases = []string{}
		}
		locations = append(locations, l)
	}

	return locations, nil
}

func nilIfBlank(s string) *string {
	return nilIfEmpty(strings.TrimSpace(s))
}
