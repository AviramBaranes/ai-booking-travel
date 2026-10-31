package location

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type MergeSuggestionGroup struct {
	CountryCode    string          `json:"country_code"`
	NormalizedName string          `json:"normalized_name"`
	SharedBroker   bool            `json:"shared_broker"`
	Locations      []MergeLocation `json:"locations"`
}

type ListMergeSuggestionsResponse struct {
	Groups []MergeSuggestionGroup `json:"groups"`
}

// ListMergeSuggestions groups locations whose names differ only in punctuation, spacing or case.
func (s *LocationService) ListMergeSuggestions(ctx context.Context) (*ListMergeSuggestionsResponse, error) {
	rows, err := s.query.ListLocationMergeSuggestions(ctx)
	if err != nil {
		rlog.Error("failed to list location merge suggestions", "error", err)
		return nil, api_errors.ErrInternalError
	}

	var ids []int64
	for _, row := range rows {
		ids = append(ids, row.LocationIds...)
	}

	locations, err := loadMergeLocations(ctx, s.query, ids)
	if err != nil {
		rlog.Error("failed to load merge suggestion locations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	byID := make(map[int64]MergeLocation, len(locations))
	for _, l := range locations {
		byID[l.ID] = l
	}

	groups := make([]MergeSuggestionGroup, 0, len(rows))
	for _, row := range rows {
		group := MergeSuggestionGroup{
			CountryCode:    row.CountryCode,
			NormalizedName: row.NormalizedName,
			Locations:      make([]MergeLocation, 0, len(row.LocationIds)),
		}
		for _, id := range row.LocationIds {
			l, ok := byID[id]
			if !ok {
				continue
			}
			for _, other := range group.Locations {
				if mergeLocationsConflict(other, l) {
					group.SharedBroker = true
				}
			}
			group.Locations = append(group.Locations, l)
		}
		groups = append(groups, group)
	}

	return &ListMergeSuggestionsResponse{Groups: groups}, nil
}

// mergeLocationsConflict reports whether two locations hold stations that cannot end up in one
// location, which blocks merging them.
func mergeLocationsConflict(a, b MergeLocation) bool {
	for _, ca := range a.BrokerCodes {
		for _, cb := range b.BrokerCodes {
			if brokerCodesConflict(db.Broker(ca.Broker), ca.Brands, db.Broker(cb.Broker), cb.Brands) {
				return true
			}
		}
	}
	return false
}
