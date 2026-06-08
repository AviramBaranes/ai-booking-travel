package queries

import (
	"context"
	"encoding/json"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type SelectedAddon struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type GetFullReservationResponse struct {
	// Identity
	ID                    int64  `json:"id"`
	BrokerReservationID   string `json:"brokerReservationId"`
	UserID                int64  `json:"userId"`
	OfficeID              *int64 `json:"officeId,omitempty"`
	OrganizationID        *int64 `json:"organizationId,omitempty"`
	IsOrganizationOrganic *bool  `json:"isOrganizationOrganic,omitempty"`
	AdminRefID            *int64 `json:"adminRefId,omitempty"`

	// Status
	ReservationStatus string `json:"reservationStatus"`
	PaymentStatus     string `json:"paymentStatus"`

	// Broker / supplier
	Broker       string `json:"broker"`
	SupplierCode string `json:"supplierCode"`

	// Car details (flattened from JSONB)
	Model        string `json:"model"`
	CarGroup     string `json:"carGroup"`
	ImageURL     string `json:"imageUrl"`
	SupplierName string `json:"supplierName"`
	CarType      string `json:"carType"`
	Acriss       string `json:"acriss"`
	HasAC        bool   `json:"hasAC"`
	IsAutoGear   bool   `json:"isAutoGear"`
	IsElectric   bool   `json:"isElectric"`
	Seats        int    `json:"seats"`
	Bags         int    `json:"bags"`
	Doors        int    `json:"doors"`

	// Plan
	PlanInclusions []string `json:"planInclusions"`

	// Pay-at-pickup (flattened from JSONB)
	Fees           broker.Fees     `json:"fees"`
	SelectedAddons []SelectedAddon `json:"selectedAddons"`

	// Pricing
	CurrencyCode                   string  `json:"currencyCode"`
	CurrencyRate                   float64 `json:"currencyRate"`
	VatPercentage                  float64 `json:"vatPercentage"`
	PurchasePrice                  float64 `json:"purchasePrice"`
	PurchasePriceInILS             float64 `json:"purchasePriceInILS"`
	MarkupPercentage               float64 `json:"markupPercentage"`
	CarSellPriceWithBrokerERP      float64 `json:"carSellPriceWithBrokerERP"`
	CarSellPriceWithBrokerERPInILS float64 `json:"carSellPriceWithBrokerERPInILS"`
	BrokerErpPrice                 float64 `json:"brokerErpPrice"`
	BtErpPrice                     float64 `json:"btErpPrice"`
	BtErpPriceInILS                float64 `json:"btErpPriceInILS"`
	DiscountPercentage             float64 `json:"discountPercentage"`
	TotalPrice                     float64 `json:"totalPrice"`
	TotalPriceInILS                float64 `json:"totalPriceInILS"`
	Profit                         float64 `json:"profit"`
	ProfitInILS                    float64 `json:"profitInILS"`
	ProfitPercentage               float64 `json:"profitPercentage"`

	// Trip details
	FlightNumber        *string `json:"flightNumber,omitempty"`
	CountryCode         string  `json:"countryCode"`
	PickupDate          string  `json:"pickupDate"`
	DropoffDate         string  `json:"dropoffDate"`
	PickupTime          string  `json:"pickupTime"`
	DropoffTime         string  `json:"dropoffTime"`
	RentalDays          int32   `json:"rentalDays"`
	PickupLocationName  string  `json:"pickupLocationName"`
	DropoffLocationName string  `json:"dropoffLocationName"`

	// Driver
	DriverTitle     string `json:"driverTitle"`
	DriverFirstName string `json:"driverFirstName"`
	DriverLastName  string `json:"driverLastName"`
	DriverAge       int32  `json:"driverAge"`

	// Voucher
	VoucherNumber *string `json:"voucherNumber,omitempty"`
	VoucheredAt   string  `json:"voucheredAt,omitempty"`

	// Timestamps
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (s *QueryService) GetFullReservation(ctx context.Context, reservationID int64) (*GetFullReservationResponse, error) {
	reservation, err := s.query.GetReservationByID(ctx, reservationID)
	if err != nil {
		rlog.Error("failed to get full reservation", "error", err, "reservationID", reservationID)
		return nil, api_errors.ErrInternalError
	}

	return buildReservationResponse(reservation), nil
}

func buildReservationResponse(reservation db.Reservation) *GetFullReservationResponse {
	var carDetails broker.CarDetails
	_ = json.Unmarshal(reservation.CarDetails, &carDetails)

	var payAtPickup struct {
		Fees           broker.Fees     `json:"fees"`
		SelectedAddons []SelectedAddon `json:"selectedAddons"`
	}
	_ = json.Unmarshal(reservation.PayAtPickup, &payAtPickup)

	selectedAddons := payAtPickup.SelectedAddons
	if selectedAddons == nil {
		selectedAddons = []SelectedAddon{}
	}

	currencyRate := dbadapters.NumericToFloat64(reservation.CurrencyRate)
	purchasePrice := dbadapters.NumericToFloat64(reservation.PurchasePrice)
	markupPercentage := dbadapters.NumericToFloat64(reservation.MarkupPercentage)
	brokerErpPrice := dbadapters.NumericToFloat64(reservation.BrokerErpPrice)
	btErpPrice := dbadapters.NumericToFloat64(reservation.BtErpPrice)
	totalPrice := dbadapters.NumericToFloat64(reservation.TotalPrice)
	carSellPrice := pricing.ApplyMarkup(purchasePrice, markupPercentage) + pricing.ApplyMarkup(brokerErpPrice, markupPercentage)
	carPurchasePrice := purchasePrice + brokerErpPrice
	profit := carSellPrice - carPurchasePrice + btErpPrice
	var profitPercentage float64
	if carSellPrice != 0 {
		profitPercentage = (profit / carSellPrice) * 100
	}

	return &GetFullReservationResponse{
		ID:                             reservation.ID,
		BrokerReservationID:            reservation.BrokerReservationID,
		UserID:                         reservation.UserID,
		OfficeID:                       reservation.OfficeID,
		OrganizationID:                 reservation.OrganizationID,
		IsOrganizationOrganic:          reservation.IsOrganizationOrganic,
		AdminRefID:                     reservation.AdminRefID,
		ReservationStatus:              string(reservation.ReservationStatus),
		PaymentStatus:                  string(reservation.PaymentStatus),
		Broker:                         string(reservation.Broker),
		SupplierCode:                   reservation.SupplierCode,
		Model:                          carDetails.Model,
		CarGroup:                       carDetails.CarGroup,
		ImageURL:                       carDetails.ImageURL,
		SupplierName:                   carDetails.SupplierName,
		CarType:                        carDetails.CarType,
		Acriss:                         carDetails.Acriss,
		HasAC:                          carDetails.HasAC,
		IsAutoGear:                     carDetails.IsAutoGear,
		IsElectric:                     carDetails.IsElectric,
		Seats:                          carDetails.Seats,
		Bags:                           carDetails.Bags,
		Doors:                          carDetails.Doors,
		PlanInclusions:                 reservation.PlanInclusions,
		Fees:                           payAtPickup.Fees,
		SelectedAddons:                 selectedAddons,
		CurrencyCode:                   reservation.CurrencyCode,
		CurrencyRate:                   currencyRate,
		VatPercentage:                  dbadapters.NumericToFloat64(reservation.VatPercentage),
		PurchasePrice:                  purchasePrice,
		PurchasePriceInILS:             purchasePrice * currencyRate,
		MarkupPercentage:               markupPercentage,
		CarSellPriceWithBrokerERP:      carSellPrice,
		CarSellPriceWithBrokerERPInILS: carSellPrice * currencyRate,
		BrokerErpPrice:                 brokerErpPrice,
		BtErpPrice:                     btErpPrice,
		BtErpPriceInILS:                btErpPrice * currencyRate,
		DiscountPercentage:             dbadapters.NumericToFloat64(reservation.DiscountPercentage),
		TotalPrice:                     totalPrice,
		TotalPriceInILS:                totalPrice * currencyRate,
		Profit:                         profit,
		ProfitInILS:                    profit * currencyRate,
		ProfitPercentage:               profitPercentage,
		FlightNumber:                   reservation.FlightNumber,
		CountryCode:                    reservation.CountryCode,
		PickupDate:                     dbadapters.DateToString(reservation.PickupDate),
		DropoffDate:                    dbadapters.DateToString(reservation.DropoffDate),
		PickupTime:                     reservation.PickupTime,
		DropoffTime:                    reservation.DropoffTime,
		RentalDays:                     reservation.RentalDays,
		PickupLocationName:             reservation.PickupLocationName,
		DropoffLocationName:            reservation.DropoffLocationName,
		DriverTitle:                    reservation.DriverTitle,
		DriverFirstName:                reservation.DriverFirstName,
		DriverLastName:                 reservation.DriverLastName,
		DriverAge:                      reservation.DriverAge,
		VoucherNumber:                  reservation.VoucherNumber,
		VoucheredAt:                    dbadapters.TimestamptzToString(reservation.VoucheredAt),
		CreatedAt:                      dbadapters.TimestamptzToString(reservation.CreatedAt),
		UpdatedAt:                      dbadapters.TimestamptzToString(reservation.UpdatedAt),
	}
}
