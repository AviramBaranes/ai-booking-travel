package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/internal/validation"
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type SelectedAddon struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type PayAtPickup struct {
	Deposit         int             `json:"deposit"`
	DepositCurrency string          `json:"depositCurrency"`
	Fees            broker.Fees     `json:"fees"`
	SelectedAddons  []SelectedAddon `json:"selectedAddons" encore:"optional"`
}

// CreateReservationParams defines the parameters required to create a reservation.
type CreateReservationParams struct {
	UserID                int64                           `json:"userId" validate:"required"`
	OfficeID              *int64                          `json:"officeId,omitempty" validate:"omitempty"`
	OrganizationID        *int64                          `json:"organizationId,omitempty" validate:"omitempty"`
	IsOrganizationOrganic *bool                           `json:"isOrganizationOrganic,omitempty" validate:"omitempty"`
	AdminRefID            *int64                          `json:"adminRefId,omitempty" validate:"omitempty"`
	BrokerReservationID   string                          `json:"brokerReservationId" validate:"required,notblank"`
	Broker                string                          `json:"broker" validate:"required,oneof=flex hertz"`
	SupplierCode          string                          `json:"supplierCode" validate:"required,notblank"`
	CarDetails            *broker.CarDetails              `json:"carDetails" validate:"required"`
	PlanInclusions        []string                        `json:"planInclusions" validate:"required"`
	CountryCode           string                          `json:"countryCode" validate:"required,notblank"`
	CurrencyCode          string                          `json:"currencyCode" validate:"required,notblank"`
	CurrencyRate          float64                         `json:"currencyRate" validate:"required,gt=0"`
	PurchasePrice         float64                         `json:"purchasePrice" validate:"required,gt=0"`
	MarkupPercentage      float64                         `json:"markupPercentage" validate:"required,gt=0"`
	DiscountPercentage    float64                         `json:"discountPercentage" validate:"gte=0,lte=100"`
	CouponName            string                          `json:"couponName" encore:"optional" validate:"omitempty"` //only customers can book with a coupon, so this is empty for agent reservations
	BrokerErpPrice        float64                         `json:"brokerErpPrice" validate:"gte=0"`
	BtErpPrice            float64                         `json:"btErpPrice" validate:"gte=0"`
	PickupDate            string                          `json:"pickupDate" validate:"required,datetime=2006-01-02"`
	DropoffDate           string                          `json:"dropoffDate" validate:"required,datetime=2006-01-02"`
	PickupTime            string                          `json:"pickupTime" validate:"required,notblank"`
	DropoffTime           string                          `json:"dropoffTime" validate:"required,notblank"`
	RentalDays            int                             `json:"rentalDays" validate:"required,gte=1"`
	DriverTitle           string                          `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Ms' 'Mrs' 'Miss' 'Dr'"`
	DriverFirstName       string                          `json:"driverFirstName" validate:"required,notblank"`
	DriverLastName        string                          `json:"driverLastName" validate:"required,notblank"`
	DriverAge             int                             `json:"driverAge" validate:"required,gte=18"`
	PickupLocationName    string                          `json:"pickupBrokerLocationId" validate:"required,notblank"`
	DropoffLocationName   string                          `json:"dropoffBrokerLocationId" validate:"required,notblank"`
	FlightNumber          *string                         `json:"flightNumber,omitempty" validate:"omitempty"`
	PayAtPickup           PayAtPickup                     `json:"payAtPickup"`
	Excess                int                             `json:"excess" validate:"gte=0"`
	ExcessCurrency        string                          `json:"excessCurrency" validate:"omitempty"` //not every plan carries an excess, so an empty currency is valid
	PickupLocationCode    string                          `json:"pickupLocationCode" validate:"required,notblank"`
	DropoffLocationCode   string                          `json:"dropoffLocationCode" validate:"required,notblank"`
	SupplierTerms         []broker.TermsAndConditionsItem `json:"supplierTerms" encore:"optional" validate:"omitempty"`
	PickupDetails         broker.StationInfo              `json:"pickupDetails" encore:"optional" validate:"omitempty"`
	DropoffDetails        broker.StationInfo              `json:"dropoffDetails" encore:"optional" validate:"omitempty"`
}

// Validate validates the fields of CreateReservationParams.
func (p CreateReservationParams) Validate() error {
	return validation.ValidateStruct(p)
}

// CreateReservationResponse is the response returned after creating a reservation.
type CreateReservationResponse struct {
	ID int64 `json:"id"`
}

func (s *ActionService) CreateReservation(ctx context.Context, p CreateReservationParams) (*CreateReservationResponse, error) {
	carDetailsJSON, err := json.Marshal(p.CarDetails)
	if err != nil {
		rlog.Error("failed to marshal reservation car details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	payAtPickupJSON, err := json.Marshal(p.PayAtPickup)
	if err != nil {
		rlog.Error("failed to marshal reservation pay at pickup details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	pickupDetailsJSON, err := json.Marshal(p.PickupDetails)
	if err != nil {
		rlog.Error("failed to marshal reservation pickup details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	dropoffDetailsJSON, err := json.Marshal(p.DropoffDetails)
	if err != nil {
		rlog.Error("failed to marshal reservation dropoff details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	supplierTermsID := s.upsertSupplierTerms(ctx, p)

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
		DiscountPercentage:    dbadapters.NumericFromFloat64(p.DiscountPercentage),
		CouponName:            p.CouponName,
		BrokerErpPrice:        dbadapters.NumericFromFloat64(p.BrokerErpPrice),
		BtErpPrice:            dbadapters.NumericFromFloat64(p.BtErpPrice),
		VatPercentage:         dbadapters.NumericFromFloat64(s.cfg.VAT),
		TotalPrice:            dbadapters.NumericFromFloat64(totalPrice),
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
		FlightNumber:          p.FlightNumber,
		PayAtPickup:           payAtPickupJSON,
		Excess:                int32(p.Excess),
		ExcessCurrency:        p.ExcessCurrency,
		PickupLocationCode:    p.PickupLocationCode,
		DropoffLocationCode:   p.DropoffLocationCode,
		SupplierTermsID:       supplierTermsID,
		PickupDetails:         pickupDetailsJSON,
		DropoffDetails:        dropoffDetailsJSON,
	})
	if err != nil {
		rlog.Error("failed to insert reservation", "error", err)
		return nil, api_errors.ErrInternalError
	}

	if _, err := emailPublisher.Publish(ctx, emailevents.EmailEventTypeNewOrder, emailevents.NewOrderEmailPayload{
		UserID:             p.UserID,
		BookingReferenceID: p.BrokerReservationID,
		ReservationPDFData: buildReservationPDFData(id, p, totalPrice),
		DriverFullName:     fmt.Sprintf("%s %s %s", p.DriverTitle, p.DriverFirstName, p.DriverLastName),
	}); err != nil {
		rlog.Error("failed to publish new order email event", "error", err, "brokerReservationId", p.BrokerReservationID)
	}

	return &CreateReservationResponse{ID: id}, nil
}

// upsertSupplierTerms stores the supplier terms of the reservation and returns the row they were stored in.
// Terms are shared by every reservation of the same supplier at the same pickup location, so they live in
// their own table rather than being duplicated onto each reservation.
//
// The terms are informational and are resolved after the car has already been booked at the broker, so a
// failure here must never abort the reservation - the reservation is simply left without a terms reference.
func (s *ActionService) upsertSupplierTerms(ctx context.Context, p CreateReservationParams) *int64 {
	if len(p.SupplierTerms) == 0 {
		return nil
	}

	termsJSON, err := json.Marshal(p.SupplierTerms)
	if err != nil {
		rlog.Error("failed to marshal supplier terms", "brokerReservationId", p.BrokerReservationID, "error", err)
		return nil
	}

	id, err := s.query.UpsertSupplierTerms(ctx, db.UpsertSupplierTermsParams{
		Broker:           db.Broker(p.Broker),
		SupplierCode:     p.SupplierCode,
		PickupLocationID: p.PickupLocationCode,
		Terms:            termsJSON,
	})
	if err != nil {
		rlog.Error("failed to upsert supplier terms", "brokerReservationId", p.BrokerReservationID, "error", err)
		return nil
	}

	return &id
}

func buildReservationPDFData(id int64, p CreateReservationParams, totalPrice float64) emailevents.ReservationPDFData {
	carPriceWithMarkup := pricing.ApplyMarkup(p.PurchasePrice, p.MarkupPercentage)
	erpFullPrice := pricing.ApplyMarkup(p.BrokerErpPrice, p.MarkupPercentage) + p.BtErpPrice

	var discountAmount float64
	if p.DiscountPercentage > 0 {
		discountAmount = (erpFullPrice + carPriceWithMarkup) - totalPrice
	}

	addons := make([]emailevents.SelectedAddon, len(p.PayAtPickup.SelectedAddons))
	for i, addon := range p.PayAtPickup.SelectedAddons {
		addons[i] = emailevents.SelectedAddon{
			ID:       addon.ID,
			Name:     addon.Name,
			Price:    addon.Price,
			Quantity: addon.Quantity,
		}
	}

	return emailevents.ReservationPDFData{
		BrokerReservationID: p.BrokerReservationID,
		CarDetails:          *p.CarDetails,
		PlanInclusions:      p.PlanInclusions,
		CurrencyCode:        getCurrencyCode(p.CurrencyCode),
		CarFullPrice:        int(math.Round(carPriceWithMarkup)),
		DiscountAmount:      int(math.Round(discountAmount)),
		ErpPrice:            int(math.Round(erpFullPrice)),
		TotalPrice:          int(math.Round(totalPrice)),
		PayAtPickup: emailevents.PayAtPickup{
			Fees: broker.Fees{
				DropCharge:             p.PayAtPickup.Fees.DropCharge,
				DropChargeCurrency:     getCurrencyCode(p.PayAtPickup.Fees.DropChargeCurrency),
				YoungDriverFee:         p.PayAtPickup.Fees.YoungDriverFee,
				YoungDriverFeeCurrency: getCurrencyCode(p.PayAtPickup.Fees.YoungDriverFeeCurrency),
			},
			SelectedAddons: addons,
		},
		FlightNumber:        p.FlightNumber,
		PickupLocationName:  p.PickupLocationName,
		DropoffLocationName: p.DropoffLocationName,
		PickupDate:          formatDate(p.PickupDate),
		DropoffDate:         formatDate(p.DropoffDate),
		PickupTime:          p.PickupTime,
		DropoffTime:         p.DropoffTime,
		RentalDays:          int32(p.RentalDays),
		DriverTitle:         p.DriverTitle,
		DriverFirstName:     p.DriverFirstName,
		DriverLastName:      p.DriverLastName,
		DriverAge:           int32(p.DriverAge),
		Voucher:             nil,
		VoucheredAt:         nil,
		CreatedAt:           time.Now().Format("02/01/2006 15:04"),
	}
}

func formatDate(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("02/01/2006")

}

func getCurrencyCode(currencyCode string) string {
	switch currencyCode {
	case "USD":
		return "$"
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	default:
		return currencyCode
	}
}
