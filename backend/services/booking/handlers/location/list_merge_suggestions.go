package location

import (
	"context"

	"encore.app/internal/api_errors"
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
		seen := make(map[string]bool)
		for _, id := range row.LocationIds {
			l, ok := byID[id]
			if !ok {
				continue
			}
			// A location has at most one code per broker, so a repeat comes from another location.
			for _, b := range l.BrokerCodes {
				if seen[b.Broker] {
					group.SharedBroker = true
				}
				seen[b.Broker] = true
			}
			group.Locations = append(group.Locations, l)
		}
		groups = append(groups, group)
	}

	return &ListMergeSuggestionsResponse{Groups: groups}, nil
}
