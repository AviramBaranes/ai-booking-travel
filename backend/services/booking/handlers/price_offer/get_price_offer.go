package price_offer

import (
	"context"
	"encoding/json"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.app/services/reservation"
	"encore.dev/rlog"
)

// GetPriceOfferResponse is the public-facing details of a price offer (no internal pricing breakdown).
type GetPriceOfferResponse struct {
	ID                  int64                   `json:"id"`
	Status              string                  `json:"status"`
	Name                string                  `json:"name"`
	CarDetails          broker.CarDetails       `json:"carDetails"`
	PlanInclusions      []string                `json:"planInclusions"`
	IsErpIncluded       bool                    `json:"isErpIncluded"`
	CurrencyCode        string                  `json:"currencyCode"`
	TotalPrice          int32                   `json:"totalPrice"`
	PickupLocationName  string                  `json:"pickupLocationName"`
	DropoffLocationName string                  `json:"dropoffLocationName"`
	PickupDate          string                  `json:"pickupDate"`
	DropoffDate         string                  `json:"dropoffDate"`
	RentalDays          int32                   `json:"rentalDays"`
	PickupTime          string                  `json:"pickupTime"`
	DropoffTime         string                  `json:"dropoffTime"`
	DriverAge           string                  `json:"driverAge"`
	PayAtPickup         reservation.PayAtPickup `json:"payAtPickup"`
	CreatedAt           string                  `json:"createdAt"`
}

// GetAgentPriceOfferResponse is the agent-facing details of a price offer, including internal pricing.
type GetAgentPriceOfferResponse struct {
	ID                  int64                   `json:"id"`
	ReservationID       *int64                  `json:"reservationId,omitempty" encore:"optional"`
	Token               string                  `json:"token"`
	Status              string                  `json:"status"`
	Name                string                  `json:"name"`
	CarDetails          broker.CarDetails       `json:"carDetails"`
	PlanInclusions      []string                `json:"planInclusions"`
	SupplierCode        string                  `json:"supplierCode"`
	CurrencyCode        string                  `json:"currencyCode"`
	CarFullPrice        int                     `json:"priceBefDesc"`
	ErpPrice            int                     `json:"erpPrice"`
	TotalPrice          int32                   `json:"totalPrice"`
	OfferedCurrencyCode string                  `json:"offeredCurrencyCode"`
	OfferedPrice        int32                   `json:"offeredPrice"`
	PickupLocationName  string                  `json:"pickupLocationName"`
	DropoffLocationName string                  `json:"dropoffLocationName"`
	PickupLocationID    int64                   `json:"pickupLocationId"`
	DropoffLocationID   int64                   `json:"dropoffLocationId"`
	PickupDate          string                  `json:"pickupDate"`
	DropoffDate         string                  `json:"dropoffDate"`
	PickupTime          string                  `json:"pickupTime"`
	DropoffTime         string                  `json:"dropoffTime"`
	RentalDays          int32                   `json:"rentalDays"`
	DriverAge           string                  `json:"driverAge"`
	RenewedAt           string                  `json:"renewedAt"`
	PayAtPickup         reservation.PayAtPickup `json:"payAtPickup"`
	CreatedAt           string                  `json:"createdAt"`
}

// GetClientPriceOffer retrieves public price offer details by token.
func (s *PriceOfferService) GetClientPriceOffer(ctx context.Context, token string) (*GetPriceOfferResponse, error) {
	uuid := dbadapters.StringToUuid(token)
	if !uuid.Valid {
		return nil, api_errors.ErrNotFound
	}

	row, err := s.query.GetPriceOfferByToken(ctx, uuid)
	if err != nil {
		if isNotFound(err) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get price offer", "token", token, "error", err)
		return nil, api_errors.ErrInternalError
	}

	var carDetails broker.CarDetails
	if err := json.Unmarshal(row.CarDetails, &carDetails); err != nil {
		rlog.Error("failed to unmarshal car details", "token", token, "error", err)
		return nil, api_errors.ErrInternalError
	}

	isErpIncluded := (float64(row.BtErpPrice) + dbadapters.NumericToFloat64(row.BrokerErpPrice)) > 0

	return &GetPriceOfferResponse{
		ID:                  row.ID,
		Status:              string(row.Status),
		Name:                row.Name,
		CarDetails:          carDetails,
		PlanInclusions:      row.PlanInclusions,
		IsErpIncluded:       isErpIncluded,
		CurrencyCode:        row.OfferedCurrencyCode,
		TotalPrice:          row.OfferedPrice,
		PickupLocationName:  row.PickupLocation,
		DropoffLocationName: row.DropoffLocation,
		PickupDate:          dbadapters.DateToString(row.PickupDate),
		DropoffDate:         dbadapters.DateToString(row.DropoffDate),
		RentalDays:          row.RentalDays,
		PickupTime:          row.PickupTime,
		DropoffTime:         row.DropoffTime,
		DriverAge:           row.DriverAge,
		PayAtPickup:         getPayAtPickup(row.PayAtPickup),
		CreatedAt:           dbadapters.TimestamptzToString(row.CreatedAt),
	}, nil
}

// GetAgentPriceOffer retrieves price offer details for the authenticated agent, including internal pricing.
func (s *PriceOfferService) GetAgentPriceOffer(ctx context.Context, id int64) (*GetAgentPriceOfferResponse, error) {
	authData := auth.GetAuthData()

	row, err := s.query.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
		ID:      id,
		AgentID: authData.UserID,
	})
	if err != nil {
		if isNotFound(err) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get price offer", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	var carDetails broker.CarDetails
	if err := json.Unmarshal(row.CarDetails, &carDetails); err != nil {
		rlog.Error("failed to unmarshal car details", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	priceDetails := calculatePriceOfferDetails(row)

	return &GetAgentPriceOfferResponse{
		ID:                  row.ID,
		ReservationID:       row.ReservationID,
		Token:               dbadapters.UuidToString(row.Token),
		Status:              string(row.Status),
		Name:                row.Name,
		CarDetails:          carDetails,
		PlanInclusions:      row.PlanInclusions,
		SupplierCode:        row.SupplierCode,
		CurrencyCode:        row.CurrencyCode,
		CarFullPrice:        priceDetails.carFullPrice,
		ErpPrice:            priceDetails.erpPrice,
		TotalPrice:          row.TotalPrice,
		OfferedCurrencyCode: row.OfferedCurrencyCode,
		OfferedPrice:        row.OfferedPrice,
		PickupLocationName:  row.PickupLocation,
		DropoffLocationName: row.DropoffLocation,
		PickupLocationID:    row.PickupLocationID,
		DropoffLocationID:   row.DropoffLocationID,
		PickupDate:          dbadapters.DateToString(row.PickupDate),
		DropoffDate:         dbadapters.DateToString(row.DropoffDate),
		RentalDays:          row.RentalDays,
		PickupTime:          row.PickupTime,
		DropoffTime:         row.DropoffTime,
		DriverAge:           row.DriverAge,
		PayAtPickup:         getPayAtPickup(row.PayAtPickup),
		RenewedAt:           dbadapters.TimestamptzToString(row.RenewedAt),
		CreatedAt:           dbadapters.TimestamptzToString(row.CreatedAt),
	}, nil
}

// priceOfferPriceDetails holds the calculated price breakdown for an agent-facing response.
type priceOfferPriceDetails struct {
	carFullPrice int
	erpPrice     int
}

// calculatePriceOfferDetails derives the display prices from stored pricing fields.
func calculatePriceOfferDetails(offer db.GetPriceOfferByIdRow) priceOfferPriceDetails {
	pp := dbadapters.NumericToFloat64(offer.PurchasePrice)
	mp := dbadapters.NumericToFloat64(offer.MarkupPercentage)
	bErp := dbadapters.NumericToFloat64(offer.BrokerErpPrice)
	btErp := float64(offer.BtErpPrice)

	carFullPrice := pricing.RoundToInt(pricing.ApplyMarkup(pp, mp))
	erpFullPrice := pricing.RoundToInt(pricing.ApplyMarkup(bErp, mp) + btErp)

	return priceOfferPriceDetails{
		carFullPrice: carFullPrice,
		erpPrice:     erpFullPrice,
	}
}

// getPayAtPickup unmarshals the PayAtPickup JSON from the database into a struct for the response.
func getPayAtPickup(papJson []byte) reservation.PayAtPickup {
	var pap reservation.PayAtPickup
	if err := json.Unmarshal(papJson, &pap); err != nil {
		rlog.Error("failed to unmarshal pay at pickup", "error", err)
		return reservation.PayAtPickup{}
	}
	return pap
}
