package booking

import "encore.app/internal/broker"

// AvailableVehicle represents a vehicle that is available for rent, including details about the car, the rental plans, add-ons, location details, and price details.
type AvailableVehicle struct {
	Broker          broker.Name            `json:"broker"`
	CarDetails      broker.CarDetails      `json:"carDetails"`
	Plans           []Plan                 `json:"plans"`
	AddOns          []broker.AddOn         `json:"addOns"`
	LocationDetails broker.LocationDetails `json:"locationDetails"`
	PriceDetails    broker.PriceDetails    `json:"priceDetails"`
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
