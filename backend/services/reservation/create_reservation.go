package reservation

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/internal/validation"
	"encore.app/services/notifications"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

// CreateReservationParams defines the parameters required to create a reservation.
type CreateReservationParams struct {
	UserID                int64              `json:"userId" validate:"required"`
	OfficeID              *int64             `json:"officeId,omitempty" validate:"omitempty"`
	OrganizationID        *int64             `json:"organizationId,omitempty" validate:"omitempty"`
	IsOrganizationOrganic *bool              `json:"isOrganizationOrganic,omitempty" validate:"omitempty"`
	AdminRefID            *int64             `json:"adminRefId,omitempty" validate:"omitempty"`
	BrokerReservationID   string             `json:"brokerReservationId" validate:"required,notblank"`
	Broker                string             `json:"broker" validate:"required,oneof=flex hertz"`
	SupplierCode          string             `json:"supplierCode" validate:"required,notblank"`
	CarDetails            *broker.CarDetails `json:"carDetails" validate:"required"`
	PlanInclusions        []string           `json:"planInclusions" validate:"required"`
	CountryCode           string             `json:"countryCode" validate:"required,notblank"`
	CurrencyCode          string             `json:"currencyCode" validate:"required,notblank"`
	CurrencyRate          float64            `json:"currencyRate" validate:"required,gt=0"`
	PurchasePrice         float64            `json:"purchasePrice" validate:"required,gt=0"`
	MarkupPercentage      float64            `json:"markupPercentage" validate:"required,gt=0"`
	DiscountPercentage    int                `json:"discountPercentage" validate:"gte=0,lte=100"`
	BrokerErpPrice        float64            `json:"brokerErpPrice" validate:"gte=0"`
	BtErpPrice            int                `json:"btErpPrice" validate:"gte=0"`
	PickupDate            string             `json:"pickupDate" validate:"required,datetime=2006-01-02"`
	DropoffDate           string             `json:"dropoffDate" validate:"required,datetime=2006-01-02"`
	PickupTime            string             `json:"pickupTime" validate:"required,notblank"`
	DropoffTime           string             `json:"dropoffTime" validate:"required,notblank"`
	RentalDays            int                `json:"rentalDays" validate:"required,gte=1"`
	DriverTitle           string             `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Ms' 'Mrs' 'Miss' 'Dr'"`
	DriverFirstName       string             `json:"driverFirstName" validate:"required,notblank"`
	DriverLastName        string             `json:"driverLastName" validate:"required,notblank"`
	DriverAge             int                `json:"driverAge" validate:"required,gte=18"`
	PickupLocationName    string             `json:"pickupBrokerLocationId" validate:"required,notblank"`
	DropoffLocationName   string             `json:"dropoffBrokerLocationId" validate:"required,notblank"`
}

// Validate validates the fields of CreateReservationParams.
func (p CreateReservationParams) Validate() error {
	return validation.ValidateStruct(p)
}

// CreateReservationResponse is the response returned after creating a reservation.
type CreateReservationResponse struct {
	ID int64 `json:"id"`
}

// encore:api private
func (s *Service) CreateReservation(ctx context.Context, p CreateReservationParams) (*CreateReservationResponse, error) {
	carDetailsJSON, err := json.Marshal(p.CarDetails)
	if err != nil {
		rlog.Error("failed to marshal reservation car details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	totalPrice := pricing.CalculateTotalPrice(p.PurchasePrice, p.MarkupPercentage, p.BrokerErpPrice, p.BtErpPrice, p.DiscountPercentage)

	id, err := s.query.InsertReservation(ctx, db.InsertReservationParams{
		UserID:                p.UserID,
		OfficeID:              p.OfficeID,
		OrganizationID:        p.OrganizationID,
		IsOrganizationOrganic: p.IsOrganizationOrganic,
		AdminRefID:            p.AdminRefID,
		BrokerReservationID:   p.BrokerReservationID,
		Broker:                db.Broker(p.Broker),
		SupplierCode:          p.SupplierCode,
		CarDetails:            carDetailsJSON,
		PlanInclusions:        p.PlanInclusions,
		CountryCode:           p.CountryCode,
		CurrencyCode:          p.CurrencyCode,
		CurrencyRate:          dbadapters.NumericFromFloat64(p.CurrencyRate),
		PurchasePrice:         dbadapters.NumericFromFloat64(p.PurchasePrice),
		MarkupPercentage:      dbadapters.NumericFromFloat64(p.MarkupPercentage),
		DiscountPercentage:    int32(p.DiscountPercentage),
		BrokerErpPrice:        dbadapters.NumericFromFloat64(p.BrokerErpPrice),
		BtErpPrice:            int32(p.BtErpPrice),
		VatPercentage:         dbadapters.NumericFromFloat64(cfg.VAT()),
		TotalPrice:            int32(totalPrice),
		PickupDate:            dbadapters.DateFromString(p.PickupDate),
		DropoffDate:           dbadapters.DateFromString(p.DropoffDate),
		PickupTime:            p.PickupTime,
		DropoffTime:           p.DropoffTime,
		RentalDays:            int32(p.RentalDays),
		DriverTitle:           p.DriverTitle,
		DriverFirstName:       p.DriverFirstName,
		DriverLastName:        p.DriverLastName,
		DriverAge:             int32(p.DriverAge),
		PickupLocationName:    p.PickupLocationName,
		DropoffLocationName:   p.DropoffLocationName,
	})
	if err != nil {
		rlog.Error("failed to insert reservation", "error", err)
		return nil, api_errors.ErrInternalError
	}

	if _, err := notifications.PublishEmailEvent(ctx, notifications.EmailEventTypeNewOrder, notifications.NewOrderEmailPayload{
		UserID:             p.UserID,
		BookingReferenceID: p.BrokerReservationID,
		DriverFullName:     fmt.Sprintf("%s %s %s", p.DriverTitle, p.DriverFirstName, p.DriverLastName),
	}); err != nil {
		rlog.Error("failed to publish new order email event", "error", err, "brokerReservationId", p.BrokerReservationID)
	}

	return &CreateReservationResponse{ID: id}, nil
}
