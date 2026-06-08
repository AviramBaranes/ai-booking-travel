package events

import "encore.app/internal/broker"

type ReservationPDFData struct {
	BrokerReservationID string            `json:"brokerReservationId"`
	CarDetails          broker.CarDetails `json:"carDetails"`
	PlanInclusions      []string          `json:"planInclusions"`
	CurrencyCode        string            `json:"currencyCode"`
	CarFullPrice        int               `json:"priceBefDesc"`
	DiscountAmount      int               `json:"discountAmount"`
	ErpPrice            int               `json:"erpPrice"`
	TotalPrice          int               `json:"totalPrice"`
	PayAtPickup         PayAtPickup       `json:"payAtPickup"`
	FlightNumber        *string           `json:"flightNumber,omitempty"`
	PickupLocationName  string            `json:"pickupLocationName"`
	DropoffLocationName string            `json:"dropoffLocationName"`
	PickupDate          string            `json:"pickupDate"`
	DropoffDate         string            `json:"dropoffDate"`
	PickupTime          string            `json:"pickupTime"`
	DropoffTime         string            `json:"dropoffTime"`
	RentalDays          int32             `json:"rentalDays"`
	DriverTitle         string            `json:"driverTitle"`
	DriverFirstName     string            `json:"driverFirstName"`
	DriverLastName      string            `json:"driverLastName"`
	DriverAge           int32             `json:"driverAge"`
	Voucher             *string           `json:"voucher,omitempty"`
	VoucheredAt         *string           `json:"voucheredAt,omitempty"`
	CreatedAt           string            `json:"createdAt"`
}

type PayAtPickup struct {
	Fees           broker.Fees     `json:"fees"`
	SelectedAddons []SelectedAddon `json:"selectedAddons,omitempty"`
}

type SelectedAddon struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}
