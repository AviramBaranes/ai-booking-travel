package price_offer

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/availability"
	"encore.app/services/reservation"
	"encore.dev/rlog"
)

// RenewPriceOfferResponse is returned after attempting to renew a price offer.
type RenewPriceOfferResponse struct {
	Found bool `json:"found"`
}

// RenewPriceOffer refreshes the stored pricing details for a price offer if the original plan is still available.
// snapshotID is 0 when the caller determined that no availability was found.
func (s *PriceOfferService) RenewPriceOffer(ctx context.Context, id int64) (*RenewPriceOfferResponse, error) {
	authData := auth.GetAuthData()

	offer, err := s.query.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
		ID:      id,
		AgentID: authData.UserID,
	})
	if err != nil {
		if isNotFound(err) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get price offer for renewal", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if !offer.RenewedAt.Valid {
		rlog.Error("price offer has invalid renewed_at", "id", id)
		return nil, api_errors.ErrInternalError
	}
	if time.Since(offer.RenewedAt.Time) < 15*time.Minute {
		return nil, ErrOfferRenewalTooSoon
	}

	driverAge, err := strconv.Atoi(offer.DriverAge)
	if err != nil {
		rlog.Error("failed to parse price offer driver age", "id", id, "driverAge", offer.DriverAge, "error", err)
		return nil, api_errors.ErrInternalError
	}

	snapshotID, err := s.searchAvailability(ctx,
		offer.PickupLocationID, offer.DropoffLocationID,
		dbadapters.DateToString(offer.PickupDate), dbadapters.DateToString(offer.DropoffDate),
		offer.PickupTime, offer.DropoffTime,
		driverAge,
	)
	if err != nil {
		return nil, err
	}
	if snapshotID == 0 {
		return s.markRenewedPriceOfferUnavailable(ctx, offer)
	}

	snapshot, err := s.getSnapshot(ctx, snapshotID)
	if err != nil {
		return nil, err
	}

	plan, err := findRenewalPlan(snapshot, offer)
	if err != nil {
		if err == errPlanNotFound {
			return s.markRenewedPriceOfferUnavailable(ctx, offer)
		}
		return nil, err
	}

	if err := s.renewPriceOfferDetails(ctx, offer, plan); err != nil {
		return nil, err
	}

	return &RenewPriceOfferResponse{Found: true}, nil
}

func findRenewalPlan(snapshot db.AvailablePlansSnapshot, offer db.GetPriceOfferByIdRow) (availability.PlanPriceDetails, error) {
	var offerCarDetails broker.CarDetails
	if err := json.Unmarshal(offer.CarDetails, &offerCarDetails); err != nil {
		rlog.Error("failed to unmarshal price offer car details", "id", offer.ID, "error", err)
		return availability.PlanPriceDetails{}, api_errors.ErrInternalError
	}

	var plans []availability.PlanPriceDetails
	if err := json.Unmarshal(snapshot.Plans, &plans); err != nil {
		rlog.Error("failed to unmarshal plans JSON", "error", err)
		return availability.PlanPriceDetails{}, api_errors.ErrInternalError
	}

	for _, plan := range plans {
		if plan.SupplierCode == offer.SupplierCode &&
			plan.CarDetails.Model == offerCarDetails.Model &&
			plan.CarDetails.SupplierName == offerCarDetails.SupplierName &&
			plan.CarDetails.Acriss == offerCarDetails.Acriss {
			return plan, nil
		}
	}

	return availability.PlanPriceDetails{}, errPlanNotFound
}

func (s *PriceOfferService) markRenewedPriceOfferUnavailable(ctx context.Context, offer db.GetPriceOfferByIdRow) (*RenewPriceOfferResponse, error) {
	err := s.query.RenewPriceOfferUnavailable(ctx, db.RenewPriceOfferUnavailableParams{
		ID:      offer.ID,
		AgentID: offer.AgentID,
	})
	if err != nil {
		rlog.Error("failed to mark renewed price offer unavailable", "id", offer.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}
	return &RenewPriceOfferResponse{Found: false}, nil
}

func (s *PriceOfferService) renewPriceOfferDetails(ctx context.Context, offer db.GetPriceOfferByIdRow, plan availability.PlanPriceDetails) error {
	carDetailsJSON, err := json.Marshal(plan.CarDetails)
	if err != nil {
		rlog.Error("failed to marshal renewed price offer car details", "id", offer.ID, "error", err)
		return api_errors.ErrInternalError
	}

	payAtPickupJSON, err := renewPayAtPickup(offer, plan)
	if err != nil {
		return err
	}

	var brokerErpPrice, btErpPrice float64
	if isPriceOfferErpIncluded(offer) {
		btErpPrice = plan.ChargedERPPriceWithVat
		brokerErpPrice = plan.SupplierErpPrice
	}

	totalPrice := pricing.CalculateTotalPrice(plan.CarPurchasePrice, plan.MarkupPercentage, brokerErpPrice, btErpPrice, plan.DiscountPercentage)

	err = s.query.RenewPriceOfferDetails(ctx, db.RenewPriceOfferDetailsParams{
		ID:               offer.ID,
		AgentID:          offer.AgentID,
		CarDetails:       carDetailsJSON,
		PlanInclusions:   plan.Inclusions,
		CurrencyCode:     plan.CurrencyCode,
		PurchasePrice:    dbadapters.NumericFromFloat64(plan.CarPurchasePrice),
		MarkupPercentage: dbadapters.NumericFromFloat64(plan.MarkupPercentage),
		BrokerErpPrice:   dbadapters.NumericFromFloat64(brokerErpPrice),
		BtErpPrice:       dbadapters.NumericFromFloat64(btErpPrice),
		TotalPrice:       dbadapters.NumericFromFloat64(totalPrice),
		PayAtPickup:      payAtPickupJSON,
	})
	if err != nil {
		rlog.Error("failed to renew price offer details", "id", offer.ID, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}

// renewPayAtPickup reconciles the originally requested pay-at-pickup details against the
// currently available add-ons, deposit and fees from the plan. Add-ons that are no longer
// available are dropped, and add-ons whose requested quantity exceeds the currently allowed
// quantity are clamped to the maximum allowed.
func renewPayAtPickup(offer db.GetPriceOfferByIdRow, plan availability.PlanPriceDetails) ([]byte, error) {
	var originalPap reservation.PayAtPickup
	if err := json.Unmarshal(offer.PayAtPickup, &originalPap); err != nil {
		rlog.Error("failed to unmarshal price offer pay at pickup", "id", offer.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	renewedAddons := make([]reservation.SelectedAddon, 0, len(originalPap.SelectedAddons))
	for _, requested := range originalPap.SelectedAddons {
		planAddOn, ok := findAvailableAddOn(plan.AvailableAddOns, requested.ID)
		if !ok || planAddOn.AllowedQuantity <= 0 {
			continue
		}

		quantity := requested.Quantity
		if quantity > planAddOn.AllowedQuantity {
			quantity = planAddOn.AllowedQuantity
		}

		renewedAddons = append(renewedAddons, reservation.SelectedAddon{
			ID:       requested.ID,
			Name:     requested.Name,
			Price:    requested.Price,
			Quantity: quantity,
		})
	}

	renewedPap := reservation.PayAtPickup{
		Deposit:         plan.Deposit,
		DepositCurrency: plan.DepositCurrency,
		Fees:            plan.Fees,
		SelectedAddons:  renewedAddons,
	}

	payAtPickupJSON, err := json.Marshal(renewedPap)
	if err != nil {
		rlog.Error("failed to marshal renewed price offer pay at pickup", "id", offer.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return payAtPickupJSON, nil
}

func findAvailableAddOn(availableAddOns []broker.AddOn, id int) (broker.AddOn, bool) {
	for _, addOn := range availableAddOns {
		if addOn.ID == id {
			return addOn, true
		}
	}
	return broker.AddOn{}, false
}

func isPriceOfferErpIncluded(offer db.GetPriceOfferByIdRow) bool {
	return dbadapters.NumericToFloat64(offer.BtErpPrice) != 0 || dbadapters.NumericToFloat64(offer.BrokerErpPrice) != 0
}
