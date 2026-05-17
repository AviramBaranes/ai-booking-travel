package price_offer_handlers

import (
	"context"
	"encoding/json"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
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

// PriceOfferResponse is returned after successfully creating a price offer.
type PriceOfferResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

// CreatePriceOffer creates a new price offer.
func (s *PriceOfferService) CreatePriceOffer(ctx context.Context, p CreatePriceOfferParams) (*PriceOfferResponse, error) {
	snapshot, err := s.getSnapshot(ctx, p.SnapshotID)
	if err != nil {
		return nil, err
	}

	plan, err := findPlan(snapshot, p.RateQualifier, p.SupplierCode)
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
	if p.IncludeERP {
		btErpPrice = plan.ChargedERPPriceWithVat
		brokerErpPrice = plan.SupplierErpPrice
	}

	rentalDays, err := calculateRentalDays(snapshot)
	if err != nil {
		rlog.Error("failed to calculate rental days for snapshot", "snapshotId", p.SnapshotID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	pickupLocationID, dropoffLocationID, err := s.getLocationIDs(ctx, plan.PickupLocationCode, plan.DropoffLocationCode, string(plan.Broker))
	if err != nil {
		rlog.Error("failed to get location IDs for price offer", "pickupBrokerLocationID", plan.PickupLocationCode, "dropoffBrokerLocationID", plan.DropoffLocationCode, "broker", plan.Broker, "error", err)
		return nil, api_errors.ErrInternalError
	}

	totalPrice := pricing.CalculateTotalPrice(plan.CarPurchasePrice, plan.MarkupPercentage, brokerErpPrice, btErpPrice, plan.DiscountPercentage)

	priceOffer, err := s.query.CreatePriceOffer(ctx, db.CreatePriceOfferParams{
		AgentID:             authData.UserID,
		Name:                p.Name,
		PickupLocationID:    pickupLocationID,
		DropoffLocationID:   dropoffLocationID,
		PickupDate:          snapshot.PickupDate,
		DropoffDate:         snapshot.DropoffDate,
		PickupTime:          snapshot.PickupTime,
		DropoffTime:         snapshot.DropoffTime,
		RentalDays:          int32(rentalDays),
		DriverAge:           snapshot.DriverAge,
		PlanID:              strconv.Itoa(plan.PlanID),
		Broker:              db.Broker(plan.Broker),
		RateQualifier:       p.RateQualifier,
		SupplierCode:        p.SupplierCode,
		CarDetails:          carDetailsJSON,
		PlanInclusions:      plan.Inclusions,
		CurrencyCode:        plan.CurrencyCode,
		CurrencyRate:        dbadapters.NumericFromFloat64(plan.CurrencyRate),
		PurchasePrice:       dbadapters.NumericFromFloat64(plan.CarPurchasePrice),
		MarkupPercentage:    dbadapters.NumericFromFloat64(plan.MarkupPercentage),
		BrokerErpPrice:      dbadapters.NumericFromFloat64(brokerErpPrice),
		BtErpPrice:          int32(btErpPrice),
		TotalPrice:          int32(totalPrice),
		OfferedCurrencyCode: p.OfferedCurrencyCode,
		OfferedPrice:        p.OfferedPrice,
	})
	if err != nil {
		return nil, err
	}

	return &PriceOfferResponse{
		ID:    priceOffer.ID,
		Token: dbadapters.UuidToString(priceOffer.Token),
	}, nil
}

// getLocationIDs resolves pickup and dropoff location IDs from broker location codes.
func (s *PriceOfferService) getLocationIDs(ctx context.Context, pickupBrokerLocationID, dropoffBrokerLocationID string, brokerName string) (int64, int64, error) {
	pickupLocationID, err := s.query.GetLocationIDByBrokerCode(ctx, db.GetLocationIDByBrokerCodeParams{
		Broker:           db.Broker(brokerName),
		BrokerLocationID: pickupBrokerLocationID,
	})
	if err != nil {
		return 0, 0, err
	}

	if pickupBrokerLocationID == dropoffBrokerLocationID {
		return pickupLocationID, pickupLocationID, nil
	}

	dropoffLocationID, err := s.query.GetLocationIDByBrokerCode(ctx, db.GetLocationIDByBrokerCodeParams{
		Broker:           db.Broker(brokerName),
		BrokerLocationID: dropoffBrokerLocationID,
	})
	if err != nil {
		return 0, 0, err
	}

	return pickupLocationID, dropoffLocationID, nil
}

// getSnapshot fetches an available plans snapshot by ID.
func (s *PriceOfferService) getSnapshot(ctx context.Context, snapshotID int64) (db.AvailablePlansSnapshot, error) {
	snapshot, err := s.query.GetSnapshotByID(ctx, snapshotID)
	if err != nil {
		if isNotFound(err) {
			return db.AvailablePlansSnapshot{}, errSnapshotNotFound
		}
		rlog.Error("failed to get snapshot", "snapshotID", snapshotID, "error", err)
		return db.AvailablePlansSnapshot{}, api_errors.ErrInternalError
	}
	return snapshot, nil
}

// findPlan locates the plan matching rateQualifier and supplierCode in the snapshot.
func findPlan(snapshot db.AvailablePlansSnapshot, rateQualifier, supplierCode string) (planDetails, error) {
	var plans []planDetails
	if err := json.Unmarshal(snapshot.Plans, &plans); err != nil {
		rlog.Error("failed to unmarshal plans JSON", "error", err)
		return planDetails{}, api_errors.ErrInternalError
	}

	for _, plan := range plans {
		if plan.RateQualifier == rateQualifier && plan.SupplierCode == supplierCode {
			return plan, nil
		}
	}
	return planDetails{}, errPlanNotFound
}

// calculateRentalDays calculates rental duration from snapshot dates.
func calculateRentalDays(snapshot db.AvailablePlansSnapshot) (int, error) {
	return broker.CalculateDaysCount(
		dbadapters.DateToString(snapshot.PickupDate), snapshot.PickupTime,
		dbadapters.DateToString(snapshot.DropoffDate), snapshot.DropoffTime,
	)
}
