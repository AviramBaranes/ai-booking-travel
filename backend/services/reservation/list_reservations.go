package reservation

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type ListReservationsRequest struct {
	SortBy     string `query:"sortBy" validate:"required,oneof=created_at pickup_date" encore:"optional"`
	Name       string `query:"name" encore:"optional"`
	BookingID  string `query:"bookingId" encore:"optional"`
	Status     string `query:"status" encore:"optional"`
	PickupDate string `query:"pickupDate" validate:"omitempty,datetime=2006-01-02" encore:"optional"`
	Page       int32  `query:"page" validate:"required,gte=1"`
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
	PickupLocationName  string `json:"pickupLocationName"`
	DriverTitle         string `json:"driverTitle"`
	DriverFirstName     string `json:"driverFirstName"`
	DriverLastName      string `json:"driverLastName"`
	Status              string `json:"status"`
	TotalPrice          int32  `json:"totalPrice"`
}

type ListReservationsResponse struct {
	Reservations []ReservationSummary `json:"reservations"`
	Total        int64                `json:"total"`
}

const listReservationsLimit int64 = 10

// encore:api auth method=GET path=/reservations
func (s Service) ListReservations(ctx context.Context, params ListReservationsRequest) (*ListReservationsResponse, error) {
	rows, err := s.listReservationsByUser(ctx, params)
	if err != nil {
		return nil, err
	}

	reservations := mapRowsToSummaries(rows)

	total, err := s.countReservationsByUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ListReservationsResponse{Reservations: reservations, Total: total}, nil
}

// listReservationsByUser returns a paginated list of reservations for a given user, ordered by creation date descending.
func (s Service) listReservationsByUser(ctx context.Context, p ListReservationsRequest) ([]db.ListReservationsByUserRow, error) {
	authData := accounts.GetAuthData()
	offset := int64(p.Page-1) * listReservationsLimit

	rows, err := s.query.ListReservationsByUser(ctx, db.ListReservationsByUserParams{
		UserID:     authData.UserID,
		Status:     nullStatusFromString(p.Status),
		Name:       nilIfEmpty(p.Name),
		BookingID:  nilIfEmpty(p.BookingID),
		PickupDate: db.DateFromString(p.PickupDate),
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

// countReservationsByUser returns the total number of reservations for a given user, optionally filtered by various criteria.
func (s Service) countReservationsByUser(ctx context.Context, p ListReservationsRequest) (int64, error) {
	authData := accounts.GetAuthData()

	count, err := s.query.CountReservationsByUser(ctx, db.CountReservationsByUserParams{
		UserID:     authData.UserID,
		Status:     nullStatusFromString(p.Status),
		Name:       nilIfEmpty(p.Name),
		BookingID:  nilIfEmpty(p.BookingID),
		PickupDate: db.DateFromString(p.PickupDate),
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
			PickupLocationName:  r.PickupLocationName,
			DriverTitle:         r.DriverTitle,
			DriverFirstName:     r.DriverFirstName,
			DriverLastName:      r.DriverLastName,
			Status:              string(r.Status),
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
