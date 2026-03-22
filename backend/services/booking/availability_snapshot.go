package booking

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.app/internal/broker"
	"encore.app/services/booking/db"
)

// planPriceDetails holds the full pricing breakdown for a single plan, stored as a snapshot for future reference.
type planPriceDetails struct {
	PlanID                    int         `json:"planId"`
	CarModel                  string      `json:"carModel"`
	Broker                    broker.Name `json:"broker"`
	RateQualifier             string      `json:"rateQualifier"`
	SupplierCode              string      `json:"supplierCode"`
	CurrencyCode              string      `json:"currencyCode"`
	CurrencyRate              float64     `json:"currencyRate"`
	CarPurchasePrice          float64     `json:"carPurchasePrice"`
	CarSellPriceWithVat       int         `json:"carSellPriceWithVat"`
	CarPurchasePriceWithErp   float64     `json:"carPurchasePriceWithErp"`
	CarSellPriceWithErpAndVat int         `json:"carSellPriceWithErpAndVat"`
	DiscountPercentage        int         `json:"discountPercentage"`
	ChargedERPPriceWithVat    int         `json:"chargedErpPriceWithVat"`
}

// storePlansDetails stores the given plan details in the database and returns the ID of the inserted snapshot.
func storePlansDetails(ctx context.Context, q db.Querier, plans []planPriceDetails) (int64, error) {
	plansJson, err := json.Marshal(plans)
	if err != nil {
		return 0, fmt.Errorf("marshaling plans details: %w", err)
	}

	ID, err := q.InsertAvailablePlansSnapshot(ctx, plansJson)
	if err != nil {
		return 0, fmt.Errorf("inserting available plans snapshot: %w", err)
	}

	return ID, nil
}
