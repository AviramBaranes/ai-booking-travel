package reservation

import (
	"context"
	"encoding/json"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

// CreateReservationRequest defines the parameters required to create a reservation.
type CreateReservationRequest struct {
	UserID                 int32              `json:"userId" validate:"required"`
	BrokerReservationID    string             `json:"brokerReservationId" validate:"required,notblank"`
	Broker                 string             `json:"broker" validate:"required,oneof=flex hertz"`
	SupplierCode           string             `json:"supplierCode" validate:"required,notblank"`
	CarDetails             *broker.CarDetails `json:"carDetails" validate:"required"`
	PlanInclusions         []string           `json:"planInclusions" validate:"required"`
	CountryCode            string             `json:"countryCode" validate:"required,notblank"`
	CurrencyCode           string             `json:"currencyCode" validate:"required,notblank"`
	CurrencyRate           float64            `json:"currencyRate" validate:"required,gt=0"`
	PurchasePrice          float64            `json:"purchasePrice" validate:"required,gte=0"`
	PriceBeforeDiscount    float64            `json:"priceBeforeDiscount" validate:"required,gte=0"`
	PriceAfterDiscount     float64            `json:"priceAfterDiscount" validate:"required,gte=0"`
	DiscountPercentage     int                `json:"discountPercentage" validate:"gte=0,lte=100"`
	ErpPrice               float64            `json:"erpPrice" validate:"gte=0"`
	TotalPrice             float64            `json:"totalPrice" validate:"required,gte=0"`
	PickupDate             string             `json:"pickupDate" validate:"required,datetime=2006-01-02"`
	ReturnDate             string             `json:"returnDate" validate:"required,datetime=2006-01-02"`
	RentalDays             int                `json:"rentalDays" validate:"required,gte=1"`
	DriverTitle            string             `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Ms'"`
	DriverFirstName        string             `json:"driverFirstName" validate:"required,notblank"`
	DriverLastName         string             `json:"driverLastName" validate:"required,notblank"`
	DriverAge              int                `json:"driverAge" validate:"required,gte=18"`
	PickupBrokerLocationID string             `json:"pickupBrokerLocationId" validate:"required,notblank"`
	ReturnBrokerLocationID string             `json:"returnBrokerLocationId" validate:"required,notblank"`
}

// Validate validates the fields of CreateReservationRequest.
func (p CreateReservationRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// CreateReservationResponse is the response returned after creating a reservation.
type CreateReservationResponse struct {
	ID int64 `json:"id"`
}

// encore:api private method=POST path=/reservations
func (s *Service) CreateReservation(ctx context.Context, params CreateReservationRequest) (*CreateReservationResponse, error) {
	carDetailsJSON, err := json.Marshal(params.CarDetails)
	if err != nil {
		rlog.Error("failed to marshal reservation car details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	id, err := s.query.InsertReservation(ctx, db.InsertReservationParams{
		UserID:                 int32(params.UserID),
		BrokerReservationID:    params.BrokerReservationID,
		Status:                 db.ReservationStatusBooked,
		Broker:                 db.Broker(params.Broker),
		SupplierCode:           params.SupplierCode,
		CarDetails:             carDetailsJSON,
		PlanInclusions:         params.PlanInclusions,
		CountryCode:            params.CountryCode,
		CurrencyCode:           params.CurrencyCode,
		CurrencyRate:           db.NumericFromFloat64(params.CurrencyRate),
		PurchasePrice:          db.NumericFromFloat64(params.PurchasePrice),
		PriceBeforeDiscount:    db.NumericFromFloat64(params.PriceBeforeDiscount),
		PriceAfterDiscount:     db.NumericFromFloat64(params.PriceAfterDiscount),
		DiscountPercentage:     int32(params.DiscountPercentage),
		ErpPrice:               db.NumericFromFloat64(params.ErpPrice),
		TotalPrice:             db.NumericFromFloat64(params.TotalPrice),
		PickupDate:             db.DateFromString(params.PickupDate),
		ReturnDate:             db.DateFromString(params.ReturnDate),
		RentalDays:             int32(params.RentalDays),
		DriverTitle:            params.DriverTitle,
		DriverFirstName:        params.DriverFirstName,
		DriverLastName:         params.DriverLastName,
		DriverAge:              int32(params.DriverAge),
		PickupBrokerLocationID: params.PickupBrokerLocationID,
		ReturnBrokerLocationID: params.ReturnBrokerLocationID,
	})
	if err != nil {
		rlog.Error("failed to insert reservation", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &CreateReservationResponse{ID: id}, nil
}
