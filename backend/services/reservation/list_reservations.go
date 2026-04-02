package reservation

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type ListReservationsRequest struct {
	Page int32 `query:"page" validate:"required,gte=1"`
}

func (p ListReservationsRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type ReservationSummary struct {
	ID                  int64  `json:"id"`
	BrokerReservationID string `json:"brokerReservationId"`
	CreatedAt           string `json:"createdAt"`
	CountryCode         string `json:"countryCode"`
	PickupDate          string `json:"pickupDate"`
	DriverTitle         string `json:"driverTitle"`
	DriverFirstName     string `json:"driverFirstName"`
	DriverLastName      string `json:"driverLastName"`
	Status              string `json:"status"`
	TotalPrice          int32  `json:"totalPrice"`
}

type ListReservationsResponse struct {
	Reservations []ReservationSummary `json:"reservations"`
}

const listReservationsLimit int32 = 10

// encore:api auth method=GET path=/reservations
func (s Service) ListReservations(ctx context.Context, params ListReservationsRequest) (*ListReservationsResponse, error) {
	rows, err := s.listReservationsByUser(ctx, params.Page)
	if err != nil {
		return nil, err
	}

	reservations := mapRowsToSummaries(rows)

	return &ListReservationsResponse{Reservations: reservations}, nil
}

// listReservationsByUser returns a paginated list of reservations for a given user, ordered by creation date descending.
func (s Service) listReservationsByUser(ctx context.Context, page int32) ([]db.ListReservationsByUserRow, error) {
	authData := accounts.GetAuthData()
	offset := (page - 1) * listReservationsLimit

	rows, err := s.query.ListReservationsByUser(ctx, db.ListReservationsByUserParams{
		UserID:      authData.UserID,
		QueryLimit:  listReservationsLimit,
		QueryOffset: offset,
	})
	if err != nil {
		rlog.Error("failed to list reservations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return rows, nil
}

// mapRowsToSummaries maps database rows to reservation summaries.
func mapRowsToSummaries(rows []db.ListReservationsByUserRow) []ReservationSummary {
	summaries := make([]ReservationSummary, len(rows))
	for i, r := range rows {
		summaries[i] = ReservationSummary{
			ID:                  r.ID,
			BrokerReservationID: r.BrokerReservationID,
			CreatedAt:           db.TimestamptzToString(r.CreatedAt),
			CountryCode:         r.CountryCode,
			PickupDate:          db.DateToString(r.PickupDate),
			DriverTitle:         r.DriverTitle,
			DriverFirstName:     r.DriverFirstName,
			DriverLastName:      r.DriverLastName,
			Status:              string(r.Status),
			TotalPrice:          r.TotalPrice,
		}
	}
	return summaries
}
