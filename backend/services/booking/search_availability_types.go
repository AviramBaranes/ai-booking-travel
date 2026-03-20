package booking

import "encore.app/internal/broker"

// AvailableVehicle represents a vehicle that is available for rent, including details about the car, the rental plans, add-ons, location details, and price details.
type AvailableVehicle struct {
	Broker          broker.Name     `json:"broker"`
	CarDetails      CarDetails      `json:"carDetails"`
	Plans           []Plan          `json:"plans"`
	AddOns          []AddOn         `json:"addOns"`
	LocationDetails LocationDetails `json:"locationDetails"`
	PriceDetails    PriceDetails    `json:"priceDetails"`
}

// PriceDetails represents the pricing details of a rental, including the currency, drop charge, and drop charge currency.
type PriceDetails struct {
	Currency           string `json:"currency"`
	DropCharge         int    `json:"dropCharge"`
	DropChargeCurrency string `json:"dropChargeCurrency"`
}

// LocationDetails represents the details of a rental location, including delivery collection, pickup and return branch addresses and phone numbers, location type, and pickup notes.
type LocationDetails struct {
	DeliveryCollection  string `json:"deliveryCollection"`
	PickupBranchAddress string `json:"pickupBranchAddress"`
	ReturnBranchAddress string `json:"returnBranchAddress"`
	PickupBranchPhone   string `json:"pickupBranchPhone"`
	ReturnBranchPhone   string `json:"returnBranchPhone"`
	PickupNotes         string `json:"pickupNotes"`
	LocationType        string `json:"locationType"`
}

// CarDetails represents the details of a car available for rent, including its ID, model, car group, image URL, supplier, car type, car size, ACRISS code, whether it has AC and auto gear, and the number of seats, bags, and doors.
type CarDetails struct {
	Model        string `json:"model"`
	CarGroup     string `json:"carGroup"`
	ImageURL     string `json:"imageUrl"`
	SupplierName string `json:"supplierName"`
	CarType      string `json:"carType"`
	Acriss       string `json:"acriss"`
	HasAC        bool   `json:"hasAC"`
	IsAutoGear   bool   `json:"isAutoGear"`
	Seats        int    `json:"seats"`
	Bags         int    `json:"bags"`
	Doors        int    `json:"doors"`
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

// AddOn represents an additional service or product that can be added to a rental, including its ID, name, price, allowed quantity, and rental period.
type AddOn struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Price           int    `json:"price"`
	Currency        string `json:"currency"`
	AllowedQuantity int    `json:"allowedQuantity"`
	Period          string `json:"period"`
}
