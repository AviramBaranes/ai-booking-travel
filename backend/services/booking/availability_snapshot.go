package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"encore.app/internal/broker"
	"encore.app/services/booking/db"
)

// planPriceDetails holds the full pricing breakdown for a single plan, stored as a snapshot for future reference.
type planPriceDetails struct {
	PlanID                 int               `json:"planId"`
	RateQualifier          string            `json:"rateQualifier"`
	SupplierCode           string            `json:"supplierCode"`
	Broker                 broker.Name       `json:"broker"`
	PickupLocationCode     string            `json:"pickupLocationCode"` //we store the pickup location code in the plan and not as column in the snapshot because the same snapshot can be used for different suppliers (different location codes)
	DropoffLocationCode    string            `json:"dropoffLocationCode"`
	CurrencyCode           string            `json:"currencyCode"`
	CurrencyRate           float64           `json:"currencyRate"`
	DiscountPercentage     int               `json:"discountPercentage"`
	CarPurchasePrice       float64           `json:"carPurchasePrice"`
	SupplierErpPrice       float64           `json:"supplierErpPrice"`
	MarkupPercentage       float64           `json:"markupPercentage"`
	ChargedERPPriceWithVat int               `json:"chargedErpPriceWithVat"`
	CarDetails             broker.CarDetails `json:"carDetails"`
	Inclusions             []string          `json:"inclusions"`
}

// storePlansDetails stores the given plan details in the database and returns the ID of the inserted snapshot.
func (s Service) storePlansDetails(ctx context.Context, plans []planPriceDetails, reqParams SearchAvailabilityRequest, countryCode string) (int64, error) {
	plansJson, err := json.Marshal(plans)
	if err != nil {
		return 0, fmt.Errorf("marshaling plans details: %w", err)
	}

	ID, err := s.query.InsertAvailablePlansSnapshot(ctx, db.InsertAvailablePlansSnapshotParams{
		Plans:       plansJson,
		DriverAge:   strconv.Itoa(reqParams.DriverAge),
		PickupDate:  reqParams.PickupDate,
		PickupTime:  reqParams.PickupTime,
		ReturnDate:  reqParams.DropoffDate,
		ReturnTime:  reqParams.DropoffTime,
		CountryCode: countryCode,
	})
	if err != nil {
		return 0, fmt.Errorf("inserting available plans snapshot: %w", err)
	}

	return ID, nil
}
