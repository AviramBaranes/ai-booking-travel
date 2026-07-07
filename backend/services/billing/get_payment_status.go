package billing

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts"
	"encore.app/services/accounts/handlers/auth"
	"encore.app/services/billing/db"
	"encore.dev/rlog"
)

type GetPaymentStatusResponse struct {
	PaymentStatus string              `json:"paymentStatus"`
	ReservationID *int64              `json:"reservationId,omitempty" encore:"optional"`
	Login         *auth.LoginResponse `json:"login,omitempty" encore:"optional"`
}

// encore:api public path=/billing/customer-payment-status/:token method=GET
func (s *Service) GetPaymentStatus(ctx context.Context, token string) (*GetPaymentStatusResponse, error) {
	pp, err := s.query.GetPendingPaymentByToken(ctx, dbadapters.StringToUuid(token))
	if err != nil {
		if errors.Is(err, dbadapters.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get pending customer payment by token", "error", err, "token", token)
		return nil, api_errors.ErrInternalError
	}

	if pp.Status == db.PaymentStatusCompleted || pp.Status == db.PaymentStatusFailed {
		customerLogin, err := accounts.GetCustomerToken(ctx, auth.GetCustomerTokenParams{
			UserID: pp.UserID,
		})
		if err != nil {
			rlog.Error("failed to get customer token after payment completed or failed", "error", err, "userID", pp.UserID)
			return nil, api_errors.ErrInternalError
		}
		return &GetPaymentStatusResponse{
			PaymentStatus: string(pp.Status),
			ReservationID: pp.ReservationID,
			Login:         customerLogin,
		}, nil
	}

	return &GetPaymentStatusResponse{
		PaymentStatus: string(pp.Status),
	}, nil
}
