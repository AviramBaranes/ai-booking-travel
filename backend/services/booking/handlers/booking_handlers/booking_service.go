package booking_handlers

import (
	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
)

var (
	errSnapshotNotFound = api_errors.NewErrorWithDetail(
		errs.NotFound,
		"Snapshot not found",
		api_errors.ErrorDetails{Code: api_errors.CodeSnapshotNotFound},
	)
	errPlanNotFound = api_errors.NewErrorWithDetail(
		errs.NotFound,
		"Plan not found",
		api_errors.ErrorDetails{Code: api_errors.CodePlanNotFound},
	)
	ErrBookingFailed = api_errors.NewErrorWithDetail(
		errs.Unknown,
		"Booking failed",
		api_errors.ErrorDetails{Code: api_errors.CodeBookingFailed},
	)
	ErrReservationCreationFailed = api_errors.NewErrorWithDetail(
		errs.Unknown,
		"Reservation creation failed",
		api_errors.ErrorDetails{Code: api_errors.CodeReservationCreationFailed},
	)

	errFlightNumberRequired = api_errors.NewErrorWithDetail(
		errs.InvalidArgument,
		"Flight number is required for this office",
		api_errors.ErrorDetails{Code: api_errors.CodeFlightNumberRequired},
	)

	errOldOffer = api_errors.NewErrorWithDetail(
		errs.PermissionDenied,
		"offer is too old, renew it first",
		api_errors.ErrorDetails{Code: api_errors.CodeOldOffer},
	)
)

type BookingService struct {
	query db.Querier
}

func NewBookingService(query db.Querier) *BookingService {
	return &BookingService{query: query}
}

// getBroker returns the broker implementation based on the broker stored in the db.
func getBroker(dbBroker db.Broker) (broker.Booker, error) {
	switch dbBroker {
	case db.BrokerHertz:
		return broker.NewHertz(), nil
	case db.BrokerFlex:
		return broker.NewFlex(), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}

// getCanceler returns the appropriate canceler implementation based on the broker.
func getCanceler(b db.Broker) (broker.Canceler, error) {
	switch b {
	case db.BrokerHertz:
		return broker.NewHertz(), nil
	case db.BrokerFlex:
		return broker.NewFlex(), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}
