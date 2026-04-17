package reservation

import (
	"context"
	"encoding/json"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/pricing"
	"encore.app/internal/validation"
	"encore.app/services/reservation/db"
	"encore.dev/config"
	"encore.dev/rlog"
)

type ReservationCfg struct {
	VAT config.Float64
}

var cfg = config.Load[*ReservationCfg]()

// CreateReservationRequest defines the parameters required to create a reservation.
type CreateReservationRequest struct {
	UserID              int32              `json:"userId" validate:"required"`
	BrokerReservationID string             `json:"brokerReservationId" validate:"required,notblank"`
	Broker              string             `json:"broker" validate:"required,oneof=flex hertz"`
	SupplierCode        string             `json:"supplierCode" validate:"required,notblank"`
	CarDetails          *broker.CarDetails `json:"carDetails" validate:"required"`
	PlanInclusions      []string           `json:"planInclusions" validate:"required"`
	CountryCode         string             `json:"countryCode" validate:"required,notblank"`
	CurrencyCode        string             `json:"currencyCode" validate:"required,notblank"`
	CurrencyRate        float64            `json:"currencyRate" validate:"required,gt=0"`
	PurchasePrice       float64            `json:"purchasePrice" validate:"required,gt=0"`
	MarkupPercentage    float64            `json:"markupPercentage" validate:"required,gt=0"`
	DiscountPercentage  int                `json:"discountPercentage" validate:"gte=0,lte=100"`
	BrokerErpPrice      float64            `json:"brokerErpPrice" validate:"gte=0"`
	BtErpPrice          int                `json:"btErpPrice" validate:"gte=0"`
	PickupDate          string             `json:"pickupDate" validate:"required,datetime=2006-01-02"`
	ReturnDate          string             `json:"returnDate" validate:"required,datetime=2006-01-02"`
	PickupTime          string             `json:"pickupTime" validate:"required,notblank"`
	DropoffTime         string             `json:"dropoffTime" validate:"required,notblank"`
	RentalDays          int                `json:"rentalDays" validate:"required,gte=1"`
	DriverTitle         string             `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Ms'"`
	DriverFirstName     string             `json:"driverFirstName" validate:"required,notblank"`
	DriverLastName      string             `json:"driverLastName" validate:"required,notblank"`
	DriverAge           int                `json:"driverAge" validate:"required,gte=18"`
	PickupLocationName  string             `json:"pickupBrokerLocationId" validate:"required,notblank"`
	DropoffLocationName string             `json:"dropoffBrokerLocationId" validate:"required,notblank"`
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
func (s *Service) CreateReservation(ctx context.Context, p CreateReservationRequest) (*CreateReservationResponse, error) {
	carDetailsJSON, err := json.Marshal(p.CarDetails)
	if err != nil {
		rlog.Error("failed to marshal reservation car details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	totalPrice := pricing.CalculateTotalPrice(p.PurchasePrice, p.MarkupPercentage, p.BrokerErpPrice, p.BtErpPrice, p.DiscountPercentage)

	id, err := s.query.InsertReservation(ctx, db.InsertReservationParams{
		UserID:              int32(p.UserID),
		BrokerReservationID: p.BrokerReservationID,
		Status:              db.ReservationStatusBooked,
		Broker:              db.Broker(p.Broker),
		SupplierCode:        p.SupplierCode,
		CarDetails:          carDetailsJSON,
		PlanInclusions:      p.PlanInclusions,
		CountryCode:         p.CountryCode,
		CurrencyCode:        p.CurrencyCode,
		CurrencyRate:        db.NumericFromFloat64(p.CurrencyRate),
		PurchasePrice:       db.NumericFromFloat64(p.PurchasePrice),
		MarkupPercentage:    db.NumericFromFloat64(p.MarkupPercentage),
		DiscountPercentage:  int32(p.DiscountPercentage),
		BrokerErpPrice:      db.NumericFromFloat64(p.BrokerErpPrice),
		BtErpPrice:          int32(p.BtErpPrice),
		VatPercentage:       db.NumericFromFloat64(cfg.VAT()),
		TotalPrice:          int32(totalPrice),
		PickupDate:          db.DateFromString(p.PickupDate),
		ReturnDate:          db.DateFromString(p.ReturnDate),
		PickupTime:          p.PickupTime,
		DropoffTime:         p.DropoffTime,
		RentalDays:          int32(p.RentalDays),
		DriverTitle:         p.DriverTitle,
		DriverFirstName:     p.DriverFirstName,
		DriverLastName:      p.DriverLastName,
		DriverAge:           int32(p.DriverAge),
		PickupLocationName:  p.PickupLocationName,
		DropoffLocationName: p.DropoffLocationName,
	})
	if err != nil {
		rlog.Error("failed to insert reservation", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &CreateReservationResponse{ID: id}, nil
}
