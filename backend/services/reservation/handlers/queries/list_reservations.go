package queries

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type ListReservationsParams struct {
	SortBy     string `query:"sortBy" validate:"required,oneof=created_at pickup_date" encore:"optional"`
	Name       string `query:"name" encore:"optional"`
	BookingID  string `query:"bookingId" encore:"optional"`
	Status     string `query:"status" encore:"optional"`
	PickupDate string `query:"pickupDate" validate:"omitempty,datetime=2006-01-02" encore:"optional"`
	Page       int32  `query:"page" validate:"required,gte=1"`
}

func (p ListReservationsParams) Validate() error {
	return validation.ValidateStruct(p)
}

type ReservationSummary struct {
	ID                  int64  `json:"id"`
	BrokerReservationID string `json:"brokerReservationId"`
	CreatedAt           string `json:"createdAt"`
	CountryCode         string `json:"countryCode"`
	PickupDate          string `json:"pickupDate"`
	PickupLocationName  string `json:"pickupLocationName"`
	DriverTitle         string `json:"driverTitle"`
	DriverFirstName     string `json:"driverFirstName"`
	DriverLastName      string `json:"driverLastName"`
	ReservationStatus   string `json:"reservationStatus"`
	TotalPrice          int32  `json:"totalPrice"`
}

type ListReservationsResponse struct {
	Reservations []ReservationSummary `json:"reservations"`
	Total        int64                `json:"total"`
}

const listReservationsLimit int64 = 8

func (s *QueryService) ListReservations(ctx context.Context, p ListReservationsParams) (*ListReservationsResponse, error) {
	rows, err := s.listReservationsByUser(ctx, p)
	if err != nil {
		return nil, err
	}

	reservations := mapRowsToSummaries(rows)

	total, err := s.countReservationsByUser(ctx, p)
	if err != nil {
		return nil, err
	}

	return &ListReservationsResponse{Reservations: reservations, Total: total}, nil
}

func (s *QueryService) listReservationsByUser(ctx context.Context, p ListReservationsParams) ([]db.ListReservationsByUserRow, error) {
	authData := accounts.GetAuthData()
	offset := int64(p.Page-1) * listReservationsLimit

	rows, err := s.query.ListReservationsByUser(ctx, db.ListReservationsByUserParams{
		UserID:     authData.UserID,
		Status:     nullStatusFromString(p.Status),
		Name:       nilIfEmpty(p.Name),
		BookingID:  nilIfEmpty(p.BookingID),
		PickupDate: dbadapters.DateFromString(p.PickupDate),
		SortBy:     p.SortBy,
		PageSize:   listReservationsLimit,
		PageOffset: offset,
	})

	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return []db.ListReservationsByUserRow{}, nil
		}
		rlog.Error("failed to list reservations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return rows, nil
}

func (s *QueryService) countReservationsByUser(ctx context.Context, p ListReservationsParams) (int64, error) {
	authData := accounts.GetAuthData()

	count, err := s.query.CountReservationsByUser(ctx, db.CountReservationsByUserParams{
		UserID:     authData.UserID,
		Status:     nullStatusFromString(p.Status),
		Name:       nilIfEmpty(p.Name),
		BookingID:  nilIfEmpty(p.BookingID),
		PickupDate: dbadapters.DateFromString(p.PickupDate),
	})

	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return 0, nil
		}
		rlog.Error("failed to count reservations", "error", err)
		return 0, api_errors.ErrInternalError
	}

	return count, nil
}

func mapRowsToSummaries(rows []db.ListReservationsByUserRow) []ReservationSummary {
	summaries := make([]ReservationSummary, len(rows))
	for i, r := range rows {
		summaries[i] = ReservationSummary{
			ID:                  r.ID,
			BrokerReservationID: r.BrokerReservationID,
			CreatedAt:           dbadapters.TimestamptzToString(r.CreatedAt),
			CountryCode:         r.CountryCode,
			PickupDate:          dbadapters.DateToString(r.PickupDate),
			PickupLocationName:  r.PickupLocationName,
			DriverTitle:         r.DriverTitle,
			DriverFirstName:     r.DriverFirstName,
			DriverLastName:      r.DriverLastName,
			ReservationStatus:   string(r.ReservationStatus),
			TotalPrice:          r.TotalPrice,
		}
	}
	return summaries
}

func nullStatusFromString(s string) db.NullReservationStatus {
	if s == "" {
		return db.NullReservationStatus{}
	}
	return db.NullReservationStatus{ReservationStatus: db.ReservationStatus(s), Valid: true}
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
