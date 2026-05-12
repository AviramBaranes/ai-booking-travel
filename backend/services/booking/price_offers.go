package booking

import (
	"context"
	"encoding/json"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/pricing"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// CreatePriceOfferParams defines the parameters required to create a price offer.
type CreatePriceOfferParams struct {
	SnapshotID          int64  `json:"snapshotId" validate:"required"`
	RateQualifier       string `json:"rateQualifier" validate:"required"`
	SupplierCode        string `json:"supplierCode" validate:"required"`
	IncludeERP          bool   `json:"includeERP"`
	Name                string `json:"name"`
	OfferedCurrencyCode string `json:"offeredCurrencyCode"`
	OfferedPrice        int32  `json:"offeredPrice"`
}

// CreatePriceOfferResponse represents the response returned after successfully creating a price offer, including the unique identifier and token for the created offer.
type PriceOfferResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

// CreatePriceOffer creates a new price offer based on the provided parameters, including details from the associated snapshot and plan, and returns the created offer's ID and token.
//
// encore:api auth method=POST path=/booking/price-offers tag:agent
func (s *Service) CreatePriceOffer(ctx context.Context, params CreatePriceOfferParams) (*PriceOfferResponse, error) {
	snapshot, err := s.getSnapshot(ctx, params.SnapshotID)
	if err != nil {
		return nil, err
	}

	plan, err := findPlan(snapshot, params.RateQualifier, params.SupplierCode)
	if err != nil {
		return nil, err
	}

	authData := auth.GetAuthData()

	carDetailsJSON, err := json.Marshal(plan.CarDetails)
	if err != nil {
		rlog.Error("failed to marshal reservation car details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	var btErpPrice int
	var brokerErpPrice float64
	if params.IncludeERP {
		btErpPrice = plan.ChargedERPPriceWithVat
		brokerErpPrice = plan.SupplierErpPrice
	}

	totalPrice := pricing.CalculateTotalPrice(plan.CarPurchasePrice, plan.MarkupPercentage, brokerErpPrice, btErpPrice, plan.DiscountPercentage)

	priceOffer, err := s.query.CreatePriceOffer(ctx, db.CreatePriceOfferParams{
		AgentID:             authData.UserID,
		Name:                params.Name,
		PickupLocationID:    plan.PickupLocationCode,
		DropoffLocationID:   plan.DropoffLocationCode,
		PickupDate:          snapshot.PickupDate,
		ReturnDate:          snapshot.ReturnDate,
		PickupTime:          snapshot.PickupTime,
		DropoffTime:         snapshot.ReturnTime,
		DriverAge:           snapshot.DriverAge,
		SupplierCode:        params.SupplierCode,
		CarDetails:          carDetailsJSON,
		PlanInclusions:      plan.Inclusions,
		CurrencyCode:        plan.CurrencyCode,
		PurchasePrice:       db.NumericFromFloat64(plan.CarPurchasePrice),
		MarkupPercentage:    db.NumericFromFloat64(plan.MarkupPercentage),
		BrokerErpPrice:      db.NumericFromFloat64(brokerErpPrice),
		BtErpPrice:          int32(btErpPrice),
		TotalPrice:          int32(totalPrice),
		OfferedCurrencyCode: params.OfferedCurrencyCode,
		OfferedPrice:        params.OfferedPrice,
	})

	if err != nil {
		return nil, err
	}

	return &PriceOfferResponse{
		ID:    priceOffer.ID,
		Token: db.UuidToString(priceOffer.Token),
	}, nil
}

// GetPriceOfferResponse represents the public-facing details of a price offer, exposing only the offered price (no internal pricing breakdown).
type GetPriceOfferResponse struct {
	ID                  int64             `json:"id"`
	Status              string            `json:"status"`
	Name                string            `json:"name"`
	CarDetails          broker.CarDetails `json:"carDetails"`
	PlanInclusions      []string          `json:"planInclusions"`
	IsErpIncluded       bool              `json:"isErpIncluded"`
	CurrencyCode        string            `json:"currencyCode"`
	TotalPrice          int32             `json:"totalPrice"`
	PickupLocationName  string            `json:"pickupLocationName"`
	DropoffLocationName string            `json:"dropoffLocationName"`
	PickupDate          string            `json:"pickupDate"`
	ReturnDate          string            `json:"returnDate"`
	PickupTime          string            `json:"pickupTime"`
	DropoffTime         string            `json:"dropoffTime"`
	DriverAge           string            `json:"driverAge"`
	CreatedAt           string            `json:"createdAt"`
}

// GetClientPriceOffer retrieves the details of a price offer based on the provided token, it doesn't exposed the agent internal pricing details.
//
// encore:api public method=GET path=/booking/price-offers/:token
func (s *Service) GetClientPriceOffer(ctx context.Context, token string) (*GetPriceOfferResponse, error) {
	uuid := db.StringToUuid(token)
	if !uuid.Valid {
		return nil, api_errors.ErrNotFound
	}

	row, err := s.query.GetPriceOfferByToken(ctx, uuid)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
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

	isErpIncluded := (float64(row.BtErpPrice) + db.NumericToFloat64(row.BrokerErpPrice)) > 0

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
		PickupDate:          db.DateToString(row.PickupDate),
		ReturnDate:          db.DateToString(row.ReturnDate),
		PickupTime:          row.PickupTime,
		DropoffTime:         row.DropoffTime,
		DriverAge:           row.DriverAge,
		CreatedAt:           db.TimestamptzToString(row.CreatedAt),
	}, nil
}
