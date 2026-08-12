package availability

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/booking/db"
)

// PlanPriceDetails holds the full pricing breakdown for a single plan, stored as a snapshot for future reference.
type PlanPriceDetails struct {
	PlanID                 int               `json:"planId"`
	RateQualifier          string            `json:"rateQualifier"`
	SupplierCode           string            `json:"supplierCode"`
	SupplierName           string            `json:"supplierName"` //keys the plan into the snapshot's suppliers_info, which holds the terms and station details
	Broker                 broker.Name       `json:"broker"`
	PickupLocationCode     string            `json:"pickupLocationCode"` //we store the pickup location code in the plan and not as column in the snapshot because the same snapshot can be used for different suppliers (different location codes)
	DropoffLocationCode    string            `json:"dropoffLocationCode"`
	CurrencyCode           string            `json:"currencyCode"`
	CurrencyRate           float64           `json:"currencyRate"`
	DiscountPercentage     float64           `json:"discountPercentage"`
	CouponName             string            `json:"couponName"` //the name of the coupon that produced DiscountPercentage, empty when no coupon was applied
	CarPurchasePrice       float64           `json:"carPurchasePrice"`
	SupplierErpPrice       float64           `json:"supplierErpPrice"`
	MarkupPercentage       float64           `json:"markupPercentage"`
	ChargedERPPriceWithVat float64           `json:"chargedErpPriceWithVat"`
	CarDetails             broker.CarDetails `json:"carDetails"`
	AvailableAddOns        []broker.AddOn    `json:"availableAddOns"`
	Deposit                int               `json:"deposit"`
	DepositCurrency        string            `json:"depositCurrency"`
	Excess                 int               `json:"excess"`
	ExcessCurrency         string            `json:"excessCurrency"`
	Fees                   broker.Fees       `json:"fees"`
	Inclusions             []string          `json:"inclusions"`
}

// storePlansDetails stores the given plan details in the database and returns the ID of the inserted snapshot.
// suppliersInfo is stored once per snapshot rather than per plan, since every plan of a given supplier shares it.
func (s *AvailabilityService) storePlansDetails(ctx context.Context, plans []PlanPriceDetails, suppliersInfo []broker.SupplierInfo, reqParams SearchAvailabilityParams, countryCode string) (int64, error) {
	plansJson, err := json.Marshal(plans)
	if err != nil {
		return 0, fmt.Errorf("marshaling plans details: %w", err)
	}

	suppliersInfoJson, err := json.Marshal(suppliersInfo)
	if err != nil {
		return 0, fmt.Errorf("marshaling suppliers info: %w", err)
	}

	ID, err := s.query.InsertAvailablePlansSnapshot(ctx, db.InsertAvailablePlansSnapshotParams{
		Plans:         plansJson,
		SuppliersInfo: suppliersInfoJson,
		DriverAge:     strconv.Itoa(reqParams.DriverAge),
		PickupDate:    dbadapters.DateFromString(reqParams.PickupDate),
		PickupTime:    reqParams.PickupTime,
		DropoffDate:   dbadapters.DateFromString(reqParams.DropoffDate),
		DropoffTime:   reqParams.DropoffTime,
		CountryCode:   countryCode,
	})
	if err != nil {
		return 0, fmt.Errorf("inserting available plans snapshot: %w", err)
	}

	return ID, nil
}
