package reservation

import (
	"context"
	"encoding/json"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/auth"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type GetReservationResponse struct {
	ID                     int64             `json:"id"`
	UserID                 int32             `json:"userId"`
	BrokerReservationID    string            `json:"brokerReservationId"`
	Status                 string            `json:"status"`
	Broker                 string            `json:"broker"`
	SupplierCode           string            `json:"supplierCode"`
	CarDetails             broker.CarDetails `json:"carDetails"`
	PlanInclusions         []string          `json:"planInclusions"`
	CountryCode            string            `json:"countryCode"`
	CurrencyCode           string            `json:"currencyCode"`
	CurrencyRate           float64           `json:"currencyRate"`
	PriceAfterDiscount     float64           `json:"priceAfterDiscount"`
	DiscountPercentage     int32             `json:"discountPercentage"`
	ErpPrice               float64           `json:"erpPrice"`
	TotalPrice             float64           `json:"totalPrice"`
	PickupDate             string            `json:"pickupDate"`
	ReturnDate             string            `json:"returnDate"`
	RentalDays             int32             `json:"rentalDays"`
	DriverTitle            string            `json:"driverTitle"`
	DriverFirstName        string            `json:"driverFirstName"`
	DriverLastName         string            `json:"driverLastName"`
	DriverAge              int32             `json:"driverAge"`
	PickupBrokerLocationID string            `json:"pickupBrokerLocationId"`
	ReturnBrokerLocationID string            `json:"returnBrokerLocationId"`
	CreatedAt              string            `json:"createdAt"`
}

// encore:api auth method=GET path=/reservations/:id
func (s *Service) GetReservation(ctx context.Context, id int64) (*GetReservationResponse, error) {
	row, err := s.query.GetReservationByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get reservation", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	authData := auth.GetAuthData()
	if authData.UserID != row.UserID {
		return nil, api_errors.ErrNotFound
	}

	var carDetails broker.CarDetails
	if err := json.Unmarshal(row.CarDetails, &carDetails); err != nil {
		rlog.Error("failed to unmarshal car details", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &GetReservationResponse{
		ID:                     row.ID,
		UserID:                 row.UserID,
		BrokerReservationID:    row.BrokerReservationID,
		Status:                 string(row.Status),
		Broker:                 string(row.Broker),
		SupplierCode:           row.SupplierCode,
		CarDetails:             carDetails,
		PlanInclusions:         row.PlanInclusions,
		CountryCode:            row.CountryCode,
		CurrencyCode:           row.CurrencyCode,
		CurrencyRate:           db.NumericToFloat64(row.CurrencyRate),
		PriceAfterDiscount:     db.NumericToFloat64(row.PriceAfterDiscount),
		DiscountPercentage:     row.DiscountPercentage,
		ErpPrice:               db.NumericToFloat64(row.ErpPrice),
		TotalPrice:             db.NumericToFloat64(row.TotalPrice),
		PickupDate:             db.DateToString(row.PickupDate),
		ReturnDate:             db.DateToString(row.ReturnDate),
		RentalDays:             row.RentalDays,
		DriverTitle:            row.DriverTitle,
		DriverFirstName:        row.DriverFirstName,
		DriverLastName:         row.DriverLastName,
		DriverAge:              row.DriverAge,
		PickupBrokerLocationID: row.PickupBrokerLocationID,
		ReturnBrokerLocationID: row.ReturnBrokerLocationID,
		CreatedAt:              db.TimestamptzToString(row.CreatedAt),
	}, nil
}
