package price_offer

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
)

// SearchAvailabilityFn is the function type for searching availability.
// It returns the snapshot ID, or 0 if no availability was found.
type SearchAvailabilityFn func(ctx context.Context, pickupLocID, dropoffLocID int64, pickupDate, dropoffDate, pickupTime, dropoffTime string, driverAge int) (snapshotID int64, err error)

// PriceOfferService handles price offer operations.
type PriceOfferService struct {
	query              db.Querier
	searchAvailability SearchAvailabilityFn
}

// NewPriceOfferService creates a new PriceOfferService.
// searchFn may be nil for operations that do not require availability search.
func NewPriceOfferService(query db.Querier, searchFn SearchAvailabilityFn) *PriceOfferService {
	return &PriceOfferService{query: query, searchAvailability: searchFn}
}

var (
	// ErrOfferRenewalTooSoon is returned when a renewal is attempted before the 15-minute cooldown.
	ErrOfferRenewalTooSoon = api_errors.NewErrorWithDetail(errs.PermissionDenied, "offer renewal only allowed after 15 minutes from last renewal", api_errors.ErrorDetails{
		Code: api_errors.CodeOfferRenewalTooSoon,
	})
	errSnapshotNotFound = api_errors.NewErrorWithDetail(errs.NotFound, "Snapshot not found", api_errors.ErrorDetails{
		Code: api_errors.CodeSnapshotNotFound,
	})
	errPlanNotFound = api_errors.NewErrorWithDetail(errs.NotFound, "Plan not found", api_errors.ErrorDetails{
		Code: api_errors.CodePlanNotFound,
	})
)

func isNotFound(err error) bool {
	return errors.Is(err, db.ErrNoRows)
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nullOfferStatusFromString(s string) db.NullOfferStatus {
	if s == "" {
		return db.NullOfferStatus{}
	}
	return db.NullOfferStatus{OfferStatus: db.OfferStatus(s), Valid: true}
}
