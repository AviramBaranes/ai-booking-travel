package booking

import (
	"context"

	"encore.app/services/booking/handlers/booking_handlers"
)

type GetSnapshotPlanSummaryParams struct {
	SnapshotID    int64  `json:"snapshotId"`
	RateQualifier string `json:"rateQualifier"`
	SupplierCode  string `json:"supplierCode"`
	PlanID        string `json:"planId"`
	IncludeERP    bool   `json:"includeERP"`
}

type GetSnapshotPlanSummaryResponse struct {
	TotalPrice          float64 `json:"totalPrice"`
	CurrencyCode        string  `json:"currencyCode"`
	CarModel            string  `json:"carModel"`
	PickupDate          string  `json:"pickupDate"`
	DropoffDate         string  `json:"dropoffDate"`
	PickupTime          string  `json:"pickupTime"`
	DropoffTime         string  `json:"dropoffTime"`
	PickupLocationName  string  `json:"pickupLocationName"`
	DropoffLocationName string  `json:"dropoffLocationName"`
}

// GetSnapshotPlanSummary returns the server-authoritative pricing and booking summary for a specific plan within a snapshot.
// Used by the billing service to derive the payment amount without trusting the client.
//
// encore:api private method=POST path=/snapshot/plan-summary
func (s *Service) GetSnapshotPlanSummary(ctx context.Context, p *GetSnapshotPlanSummaryParams) (*GetSnapshotPlanSummaryResponse, error) {
	bs := booking_handlers.NewBookingService(s.query)
	result, err := bs.GetSnapshotPlanSummary(ctx, booking_handlers.SnapshotPlanSummaryParams{
		SnapshotID:    p.SnapshotID,
		RateQualifier: p.RateQualifier,
		SupplierCode:  p.SupplierCode,
		PlanID:        p.PlanID,
		IncludeERP:    p.IncludeERP,
	})
	if err != nil {
		return nil, err
	}
	return &GetSnapshotPlanSummaryResponse{
		TotalPrice:          result.TotalPrice,
		CurrencyCode:        result.CurrencyCode,
		CarModel:            result.CarModel,
		PickupDate:          result.PickupDate,
		DropoffDate:         result.DropoffDate,
		PickupTime:          result.PickupTime,
		DropoffTime:         result.DropoffTime,
		PickupLocationName:  result.PickupLocationName,
		DropoffLocationName: result.DropoffLocationName,
	}, nil
}
