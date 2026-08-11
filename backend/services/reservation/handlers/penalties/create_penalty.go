package penalties

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	ErrReservationNotCanceled = api_errors.NewErrorWithDetail(
		errs.FailedPrecondition,
		"A penalty can only be charged on a canceled reservation",
		api_errors.ErrorDetails{Code: api_errors.CodeReservationNotCanceled},
	)

	ErrPenaltyAlreadyExists = api_errors.NewErrorWithDetail(
		errs.AlreadyExists,
		"This reservation already has a penalty",
		api_errors.ErrorDetails{Code: api_errors.CodePenaltyAlreadyExists},
	)
)

// CreatePenaltyParams is the request payload for recording a penalty the supplier charged us.
type CreatePenaltyParams struct {
	ReservationID int64   `json:"reservationId" validate:"required,gte=1"`
	Type          string  `json:"type" validate:"required,oneof=cancellation no_show"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

func (p CreatePenaltyParams) Validate() error {
	return validation.ValidateStruct(p)
}

type CreatePenaltyResponse struct {
	ID            int64   `json:"id"`
	ReservationID int64   `json:"reservationId"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	CurrencyCode  string  `json:"currencyCode"`
	CurrencyRate  float64 `json:"currencyRate"`
	CreatedAt     string  `json:"createdAt"`
}

// CreatePenalty records a cancellation or no-show fee against a canceled reservation. The penalty is
// charged in the reservation's currency, at the rate of the day it is recorded rather than the
// reservation's own rate, since the debt towards the customer is created now.
func (s *PenaltiesService) CreatePenalty(ctx context.Context, p CreatePenaltyParams) (*CreatePenaltyResponse, error) {
	authData := auth.Data().(*accounts.AuthData)

	reservation, err := s.query.GetReservationByID(ctx, p.ReservationID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get reservation for penalty", "error", err, "reservation_id", p.ReservationID)
		return nil, api_errors.ErrInternalError
	}

	if reservation.ReservationStatus != db.ReservationStatusCanceled {
		return nil, ErrReservationNotCanceled
	}

	currencyRate, err := s.currencyCache.GetCurrencyRate(ctx, reservation.CurrencyCode)
	if err != nil {
		rlog.Error("failed to get currency rate", "error", err, "currency_code", reservation.CurrencyCode)
		return nil, api_errors.ErrInternalError
	}

	penalty, err := s.query.InsertReservationPenalty(ctx, db.InsertReservationPenaltyParams{
		ReservationID:   reservation.ID,
		PenaltyType:     db.PenaltyType(p.Type),
		CreatedByUserID: &authData.UserID,
		CurrencyCode:    reservation.CurrencyCode,
		CurrencyRate:    dbadapters.NumericFromFloat64(currencyRate),
		Amount:          dbadapters.NumericFromFloat64(p.Amount),
	})
	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil, ErrPenaltyAlreadyExists
		}
		rlog.Error("failed to insert reservation penalty", "error", err, "reservation_id", p.ReservationID)
		return nil, api_errors.ErrInternalError
	}

	return &CreatePenaltyResponse{
		ID:            penalty.ID,
		ReservationID: penalty.ReservationID,
		Type:          string(penalty.PenaltyType),
		Amount:        dbadapters.NumericToFloat64(penalty.Amount),
		CurrencyCode:  penalty.CurrencyCode,
		CurrencyRate:  dbadapters.NumericToFloat64(penalty.CurrencyRate),
		CreatedAt:     dbadapters.TimestamptzToString(penalty.CreatedAt),
	}, nil
}
