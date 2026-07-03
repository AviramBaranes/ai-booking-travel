package booking_handlers

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.dev/rlog"
)

type SnapshotPlanSummaryParams struct {
	SnapshotID    int64
	RateQualifier string
	SupplierCode  string
	PlanID        string
	IncludeERP    bool
}

type SnapshotPlanSummaryResult struct {
	TotalPrice          float64
	CurrencyCode        string
	CarModel            string
	PickupDate          string
	DropoffDate         string
	PickupTime          string
	DropoffTime         string
	PickupLocationName  string
	DropoffLocationName string
}

func (s *BookingService) GetSnapshotPlanSummary(ctx context.Context, p SnapshotPlanSummaryParams) (*SnapshotPlanSummaryResult, error) {
	snapshot, err := s.getSnapshot(ctx, p.SnapshotID)
	if err != nil {
		return nil, err
	}

	plan, err := findPlan(snapshot, p.RateQualifier, p.SupplierCode, p.PlanID)
	if err != nil {
		return nil, err
	}

	var brokerErpPrice, btErpPrice float64
	if p.IncludeERP {
		brokerErpPrice = plan.SupplierErpPrice
		btErpPrice = plan.ChargedERPPriceWithVat
	}
	totalPrice := pricing.CalculateTotalPrice(plan.CarPurchasePrice, plan.MarkupPercentage, brokerErpPrice, btErpPrice, plan.DiscountPercentage)

	pickupLocName, dropoffLocName, err := s.getLocationsNames(ctx, plan.PickupLocationCode, plan.DropoffLocationCode)
	if err != nil {
		rlog.Error("failed to get location names for plan summary", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &SnapshotPlanSummaryResult{
		TotalPrice:          totalPrice,
		CurrencyCode:        plan.CurrencyCode,
		CarModel:            plan.CarDetails.Model,
		PickupDate:          dbadapters.DateToString(snapshot.PickupDate),
		DropoffDate:         dbadapters.DateToString(snapshot.DropoffDate),
		PickupTime:          snapshot.PickupTime,
		DropoffTime:         snapshot.DropoffTime,
		PickupLocationName:  pickupLocName,
		DropoffLocationName: dropoffLocName,
	}, nil
}
