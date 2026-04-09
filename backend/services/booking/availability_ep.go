package booking

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	"encore.dev/config"
	"encore.dev/rlog"
)

// AvailableVehicle represents a vehicle that is available for rent, including details about the car, the rental plans, add-ons, location details, and price details.
type AvailableVehicle struct {
	Broker          broker.Name            `json:"broker"`
	CarDetails      broker.CarDetails      `json:"carDetails"`
	Plans           []Plan                 `json:"plans"`
	AddOns          []broker.AddOn         `json:"addOns"`
	LocationDetails broker.LocationDetails `json:"locationDetails"`
	PriceDetails    broker.PriceDetails    `json:"priceDetails"`
	Signals         *BookingSignals        `json:"signals,omitempty" encore:"optional"`
}

// BookingSignals holds UI-facing demand and inventory indicators for a vehicle card.
type BookingSignals struct {
	LiveViewers    int      `json:"liveViewers"`
	RemainingCount int      `json:"remainingCount"`
	Tags           []string `json:"tags,omitempty"`
}

// Plan represents a rental plan, including its ID, name, description, full price, discount, and other pricing details.
type Plan struct {
	PlanID         int      `json:"planId"`
	PlanName       string   `json:"planName"`
	FullPrice      int      `json:"fullPrice"`
	Discount       int      `json:"discount"`
	Price          int      `json:"price"`
	ErpPrice       int      `json:"erpPrice"`
	PlanInclusions []string `json:"planInclusions"`
	Info           []string `json:"info"`
	RateQualifier  string   `json:"rateQualifier"`
	SupplierCode   string   `json:"supplierCode"`
}

// AvailableVehiclesConfig holds markup percentages and ERP day-charge values per broker.
type AvailableVehiclesConfig struct {
	HertzErpDayChargeUS config.Int
	HertzErpDayChargeCA config.Int
	FlexErpDayCharge    config.Int
	MarkUpGross         config.Float64
	MarkUpNet           config.Float64
}

// avCfg is the loaded AvailableVehiclesConfig for this service.
var avCfg = config.Load[*AvailableVehiclesConfig]()

// SearchAvailabilityRequest represents the request for searching availability of vehicles.
type SearchAvailabilityRequest struct {
	PickupLocationID  int64  `query:"pickupLocationId" validate:"required"`
	DropoffLocationID int64  `query:"dropoffLocationId" validate:"omitempty"`
	PickupTime        string `query:"pickupTime" validate:"required,datetime=15:04"`
	DropoffTime       string `query:"dropoffTime" validate:"required,datetime=15:04"`
	PickupDate        string `query:"pickupDate" validate:"required,datetime=2006-01-02"`
	DropoffDate       string `query:"dropoffDate" validate:"required,datetime=2006-01-02"`
	DriverAge         int    `query:"driverAge" validate:"required,gte=18"`
	CouponCode        string `query:"couponCode" validate:"omitempty"`
}

// Validate validates the fields of SearchAvailabilityRequest.
func (p SearchAvailabilityRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// SearchAvailabilityResponse represents the response for searching availability of vehicles.
type SearchAvailabilityResponse struct {
	SnapshotID          int64              `json:"snapshotId"`
	PickupLocationName  string             `json:"pickupLocationName"`
	DropoffLocationName string             `json:"dropoffLocationName"`
	AvailableVehicles   []AvailableVehicle `json:"availableVehicles"`
}

// SearchAvailability handles the http request for searching availability of vehicles.
// encore:api public method=GET path=/booking/availability
func (s *Service) SearchAvailability(ctx context.Context, p SearchAvailabilityRequest) (*SearchAvailabilityResponse, error) {
	locs, err := getLocations(ctx, s.query, p)
	if err != nil {
		return nil, err
	}

	couponDiscount, err := s.getCouponDiscount(ctx, p.CouponCode)
	if err != nil {
		return nil, err
	}

	rawVehicles, err := searchAvailabilityAcrossBrokers(p, locs)
	if err != nil {
		return nil, err
	}

	if len(rawVehicles) == 0 {
		return emptySearchAvailabilityResponse(), nil
	}

	artifacts, err := s.buildAvailabilityArtifacts(ctx, p, locs, rawVehicles, couponDiscount)
	if err != nil {
		return nil, err
	}
	if len(artifacts.plansDetails) == 0 {
		return emptySearchAvailabilityResponse(), nil
	}

	sortAvailableVehiclesByCheapestPlan(artifacts.availableCars)

	snapshotID, err := s.storePlansDetails(ctx, artifacts.plansDetails, p, extractCountryCode(locs))
	if err != nil {
		rlog.Error("failed to store plans details", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &SearchAvailabilityResponse{
		SnapshotID:          snapshotID,
		PickupLocationName:  extractPickupLocationName(locs),
		DropoffLocationName: extractDropoffLocationName(locs),
		AvailableVehicles:   artifacts.availableCars,
	}, nil
}

// emptySearchAvailabilityResponse returns an empty SearchAvailabilityResponse with no available vehicles.
func emptySearchAvailabilityResponse() *SearchAvailabilityResponse {
	return &SearchAvailabilityResponse{AvailableVehicles: []AvailableVehicle{}}
}
