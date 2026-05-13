package booking

import (
	"context"
	"encoding/json"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/pricing"
	"encore.app/internal/validation"
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
	Name                string `json:"name" validate:"required,notblank"`
	OfferedCurrencyCode string `json:"offeredCurrencyCode" validate:"required,len=3,uppercase_only"`
	OfferedPrice        int32  `json:"offeredPrice" validate:"required,gt=0"`
}

func (p CreatePriceOfferParams) Validate() error {
	return validation.ValidateStruct(p)
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

	rentalDays, err := calculateSnapshotRentalDays(snapshot)
	if err != nil {
		rlog.Error("failed to calculate rental days for snapshot", "snapshotId", params.SnapshotID, "error", err)
		return nil, api_errors.ErrInternalError
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
		RentalDays:          int32(rentalDays),
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

// GetAgentPriceOfferResponse represents the agent-facing details of a price offer, including internal pricing details.
type GetAgentPriceOfferResponse struct {
	ID                  int64             `json:"id"`
	Token               string            `json:"token"`
	Status              string            `json:"status"`
	Name                string            `json:"name"`
	CarDetails          broker.CarDetails `json:"carDetails"`
	PlanInclusions      []string          `json:"planInclusions"`
	SupplierCode        string            `json:"supplierCode"`
	CurrencyCode        string            `json:"currencyCode"`
	CarFullPrice        int               `json:"priceBefDesc"`
	ErpPrice            int               `json:"erpPrice"`
	TotalPrice          int32             `json:"totalPrice"`
	OfferedCurrencyCode string            `json:"offeredCurrencyCode"`
	OfferedPrice        int32             `json:"offeredPrice"`
	PickupLocationName  string            `json:"pickupLocationName"`
	DropoffLocationName string            `json:"dropoffLocationName"`
	PickupDate          string            `json:"pickupDate"`
	ReturnDate          string            `json:"returnDate"`
	PickupTime          string            `json:"pickupTime"`
	DropoffTime         string            `json:"dropoffTime"`
	RentalDays          int32             `json:"rentalDays"`
	DriverAge           string            `json:"driverAge"`
	CreatedAt           string            `json:"createdAt"`
}

// GetClientPriceOffer retrieves the details of a price offer based on the provided token, it doesn't exposed the agent internal pricing details.
//
// encore:api public method=GET path=/booking/price-offers/client/:token
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

// GetAgentPriceOffer retrieves the details of a price offer for the authenticated agent, including internal pricing details.
//
// encore:api auth method=GET path=/booking/price-offers/agent/:id tag:agent
func (s *Service) GetAgentPriceOffer(ctx context.Context, id int64) (*GetAgentPriceOfferResponse, error) {
	authData := auth.GetAuthData()

	row, err := s.query.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
		ID:      id,
		AgentID: authData.UserID,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
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
		Token:               db.UuidToString(row.Token),
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
		PickupDate:          db.DateToString(row.PickupDate),
		ReturnDate:          db.DateToString(row.ReturnDate),
		RentalDays:          row.RentalDays,
		PickupTime:          row.PickupTime,
		DropoffTime:         row.DropoffTime,
		DriverAge:           row.DriverAge,
		CreatedAt:           db.TimestamptzToString(row.CreatedAt),
	}, nil
}

// priceOfferPriceDetails holds the calculated price details for a price offer.
type priceOfferPriceDetails struct {
	carFullPrice int
	erpPrice     int
}

// calculatePriceOfferDetails calculates the price details for a price offer based on the given parameters.
func calculatePriceOfferDetails(offer db.GetPriceOfferByIdRow) priceOfferPriceDetails {
	pp := db.NumericToFloat64(offer.PurchasePrice)
	mp := db.NumericToFloat64(offer.MarkupPercentage)
	bErp := db.NumericToFloat64(offer.BrokerErpPrice)
	btErp := float64(offer.BtErpPrice)

	carFullPrice := pricing.RoundToInt(pricing.ApplyMarkup(pp, mp))
	erpFullPrice := pricing.RoundToInt(pricing.ApplyMarkup(bErp, mp) + btErp)

	return priceOfferPriceDetails{
		carFullPrice: carFullPrice,
		erpPrice:     erpFullPrice,
	}
}

type ListPriceOffersRequest struct {
	Name   string `query:"name" encore:"optional"`
	Status string `query:"status" encore:"optional"`
	Page   int32  `query:"page" validate:"required,gte=1"`
}

func (p ListPriceOffersRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type PriceOfferSummary struct {
	ID                  int64  `json:"id"`
	Status              string `json:"status"`
	Name                string `json:"name"`
	PickupLocationName  string `json:"pickupLocationName"`
	DropoffLocationName string `json:"dropoffLocationName"`
	PickupDate          string `json:"pickupDate"`
	ReturnDate          string `json:"returnDate"`
	PickupTime          string `json:"pickupTime"`
	DropoffTime         string `json:"dropoffTime"`
	CurrencyCode        string `json:"currencyCode"`
	TotalPrice          int32  `json:"totalPrice"`
	OfferedCurrencyCode string `json:"offeredCurrencyCode"`
	OfferedPrice        int32  `json:"offeredPrice"`
	CreatedAt           string `json:"createdAt"`
}

type ListPriceOffersResponse struct {
	PriceOffers []PriceOfferSummary `json:"priceOffers"`
	Total       int64               `json:"total"`
}

const listPriceOffersLimit int32 = 8

// ListPriceOffers returns a paginated list of the authenticated agent's price offers.
//
// encore:api auth method=GET path=/booking/price-offers tag:agent
func (s *Service) ListPriceOffers(ctx context.Context, params ListPriceOffersRequest) (*ListPriceOffersResponse, error) {
	rows, err := s.listPriceOffersByAgent(ctx, params)
	if err != nil {
		return nil, err
	}

	total, err := s.countPriceOffersByAgent(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ListPriceOffersResponse{
		PriceOffers: mapRowsToPriceOfferSummaries(rows),
		Total:       total,
	}, nil
}

// listPriceOffersByAgent returns a paginated list of price offers for the authenticated agent.
func (s *Service) listPriceOffersByAgent(ctx context.Context, p ListPriceOffersRequest) ([]db.ListPriceOffersByAgentRow, error) {
	authData := auth.GetAuthData()
	offset := (p.Page - 1) * listPriceOffersLimit

	rows, err := s.query.ListPriceOffersByAgent(ctx, db.ListPriceOffersByAgentParams{
		AgentID:    authData.UserID,
		Status:     nullOfferStatusFromString(p.Status),
		NameSearch: nilIfEmpty(p.Name),
		PageSize:   listPriceOffersLimit,
		PageOffset: offset,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return []db.ListPriceOffersByAgentRow{}, nil
		}
		rlog.Error("failed to list price offers", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return rows, nil
}

// countPriceOffersByAgent returns the total number of price offers for the authenticated agent.
func (s *Service) countPriceOffersByAgent(ctx context.Context, p ListPriceOffersRequest) (int64, error) {
	authData := auth.GetAuthData()

	count, err := s.query.CountPriceOffersByAgent(ctx, db.CountPriceOffersByAgentParams{
		AgentID:    authData.UserID,
		Status:     nullOfferStatusFromString(p.Status),
		NameSearch: nilIfEmpty(p.Name),
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return 0, nil
		}
		rlog.Error("failed to count price offers", "error", err)
		return 0, api_errors.ErrInternalError
	}

	return count, nil
}

// mapRowsToPriceOfferSummaries maps database rows to price offer summaries.
func mapRowsToPriceOfferSummaries(rows []db.ListPriceOffersByAgentRow) []PriceOfferSummary {
	summaries := make([]PriceOfferSummary, len(rows))
	for i, r := range rows {
		summaries[i] = PriceOfferSummary{
			ID:                  r.ID,
			Status:              string(r.Status),
			Name:                r.Name,
			PickupLocationName:  r.PickupLocation,
			DropoffLocationName: r.DropoffLocation,
			PickupDate:          db.DateToString(r.PickupDate),
			ReturnDate:          db.DateToString(r.ReturnDate),
			PickupTime:          r.PickupTime,
			DropoffTime:         r.DropoffTime,
			CurrencyCode:        r.CurrencyCode,
			TotalPrice:          r.TotalPrice,
			OfferedCurrencyCode: r.OfferedCurrencyCode,
			OfferedPrice:        r.OfferedPrice,
			CreatedAt:           db.TimestamptzToString(r.CreatedAt),
		}
	}
	return summaries
}

func nullOfferStatusFromString(s string) db.NullOfferStatus {
	if s == "" {
		return db.NullOfferStatus{}
	}
	return db.NullOfferStatus{OfferStatus: db.OfferStatus(s), Valid: true}
}

type UpdatePriceOfferParams struct {
	Status              *string `json:"status" encore:"optional" validate:"omitempty,oneof=open booked declined"`
	Name                *string `json:"name" encore:"optional" validate:"omitempty,notblank"`
	OfferedCurrencyCode *string `json:"offeredCurrencyCode" encore:"optional" validate:"omitempty,len=3,uppercase_only"`
	OfferedPrice        *int32  `json:"offeredPrice" encore:"optional" validate:"omitempty,gt=0"`
}

func (p UpdatePriceOfferParams) Validate() error {
	return validation.ValidateStruct(p)
}

// UpdatePriceOffer updates a price offer's mutable fields for the authenticated agent.
//
// encore:api auth method=PATCH path=/booking/price-offers/:id tag:agent
func (s *Service) UpdatePriceOffer(ctx context.Context, id int64, params UpdatePriceOfferParams) error {
	authData := auth.GetAuthData()

	var status db.NullOfferStatus
	if params.Status != nil {
		status = nullOfferStatusFromString(*params.Status)
	}

	err := s.query.UpdatePriceOffer(ctx, db.UpdatePriceOfferParams{
		ID:                  id,
		AgentID:             authData.UserID,
		Status:              status,
		Name:                params.Name,
		OfferedCurrencyCode: params.OfferedCurrencyCode,
		OfferedPrice:        params.OfferedPrice,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to update price offer", "id", id, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
