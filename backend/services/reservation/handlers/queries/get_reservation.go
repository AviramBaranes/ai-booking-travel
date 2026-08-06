package queries

import (
	"context"
	"encoding/json"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/handlers/actions"
	"encore.dev/rlog"
)

type GetReservationResponse struct {
	ID                  int64                           `json:"id"`
	BrokerReservationID string                          `json:"brokerReservationId"`
	ReservationStatus   string                          `json:"reservationStatus"`
	PaymentStatus       string                          `json:"paymentStatus"`
	CarDetails          broker.CarDetails               `json:"carDetails"`
	PlanInclusions      []string                        `json:"planInclusions"`
	CurrencyCode        string                          `json:"currencyCode"`
	CurrencyRate        float64                         `json:"currencyRate"`
	CarFullPrice        int                             `json:"priceBefDesc"`
	DiscountAmount      int                             `json:"discountAmount"`
	ErpPrice            int                             `json:"erpPrice"`
	TotalPrice          int                             `json:"totalPrice"`
	PayAtPickup         actions.PayAtPickup             `json:"payAtPickup"`
	FlightNumber        *string                         `json:"flightNumber,omitempty" encore:"optional"`
	PickupLocationName  string                          `json:"pickupLocationName"`
	DropoffLocationName string                          `json:"dropoffLocationName"`
	PickupDate          string                          `json:"pickupDate"`
	DropoffDate         string                          `json:"dropoffDate"`
	PickupTime          string                          `json:"pickupTime"`
	DropoffTime         string                          `json:"dropoffTime"`
	RentalDays          int32                           `json:"rentalDays"`
	DriverTitle         string                          `json:"driverTitle"`
	DriverFirstName     string                          `json:"driverFirstName"`
	DriverLastName      string                          `json:"driverLastName"`
	DriverAge           int32                           `json:"driverAge"`
	Voucher             *string                         `json:"voucher,omitempty" encore:"optional"`
	VoucheredAt         *string                         `json:"voucheredAt,omitempty" encore:"optional"`
	CreatedAt           string                          `json:"createdAt"`
	Excess              int32                           `json:"excess"`
	ExcessCurrency      string                          `json:"excessCurrency"`
	SupplierTerms       []broker.TermsAndConditionsItem `json:"supplierTerms,omitempty" encore:"optional"`
	PickupDetails       *broker.StationInfo             `json:"pickupDetails,omitempty" encore:"optional"`
	DropoffDetails      *broker.StationInfo             `json:"dropoffDetails,omitempty" encore:"optional"`
	TotalPriceFloat     float64
}

func (s *QueryService) GetReservation(ctx context.Context, id int64) (*GetReservationResponse, error) {
	row, err := s.query.GetReservationByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get reservation", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	authData := accounts.GetAuthData()
	if authData.UserID != row.UserID {
		return nil, api_errors.ErrNotFound
	}

	carDetails, payAtPickup, err := unmarshalJsons(row.CarDetails, row.PayAtPickup)
	if err != nil {
		rlog.Error("failed to unmarshal car details or pay at pickup", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	rpd := calculatePriceDetails(row)

	voucheredAt := dbadapters.TimestamptzToString(row.VoucheredAt)
	return &GetReservationResponse{
		ID:                  row.ID,
		BrokerReservationID: row.BrokerReservationID,
		ReservationStatus:   string(row.ReservationStatus),
		PaymentStatus:       string(row.PaymentStatus),
		CarDetails:          carDetails,
		PlanInclusions:      row.PlanInclusions,
		CurrencyCode:        row.CurrencyCode,
		CurrencyRate:        dbadapters.NumericToFloat64(row.CurrencyRate),
		CarFullPrice:        rpd.carFullPrice,
		ErpPrice:            rpd.erpPrice,
		DiscountAmount:      rpd.discountAmount,
		TotalPrice:          pricing.RoundToInt(dbadapters.NumericToFloat64(row.TotalPrice)),
		PayAtPickup:         payAtPickup,
		FlightNumber:        row.FlightNumber,
		PickupDate:          dbadapters.DateToString(row.PickupDate),
		DropoffDate:         dbadapters.DateToString(row.DropoffDate),
		PickupTime:          row.PickupTime,
		DropoffTime:         row.DropoffTime,
		RentalDays:          row.RentalDays,
		DriverTitle:         row.DriverTitle,
		DriverFirstName:     row.DriverFirstName,
		DriverLastName:      row.DriverLastName,
		DriverAge:           row.DriverAge,
		CreatedAt:           dbadapters.TimestamptzToString(row.CreatedAt),
		PickupLocationName:  row.PickupLocationName,
		DropoffLocationName: row.DropoffLocationName,
		Voucher:             row.VoucherNumber,
		VoucheredAt:         &voucheredAt,
		Excess:              row.Excess,
		ExcessCurrency:      row.ExcessCurrency,
		SupplierTerms:       s.getSupplierTerms(ctx, row),
		PickupDetails:       unmarshalStationInfo(row.PickupDetails),
		DropoffDetails:      unmarshalStationInfo(row.DropoffDetails),
		TotalPriceFloat:     dbadapters.NumericToFloat64(row.TotalPrice),
	}, nil
}

// getSupplierTerms loads the terms the reservation was booked under. They live in their own table
// because every reservation of the same supplier at the same pickup location shares them.
//
// Terms are informational, so a reservation without them - or a failure to load them - still returns
// the reservation rather than an error.
func (s *QueryService) getSupplierTerms(ctx context.Context, row db.Reservation) []broker.TermsAndConditionsItem {
	if row.SupplierTermsID == nil {
		return nil
	}

	termsJSON, err := s.query.GetSupplierTermsByID(ctx, *row.SupplierTermsID)
	if err != nil {
		rlog.Error("failed to get supplier terms", "id", row.ID, "supplierTermsID", *row.SupplierTermsID, "error", err)
		return nil
	}

	var terms []broker.TermsAndConditionsItem
	if err := json.Unmarshal(termsJSON, &terms); err != nil {
		rlog.Error("failed to unmarshal supplier terms", "id", row.ID, "error", err)
		return nil
	}

	return terms
}

// unmarshalStationInfo reads a pickup or dropoff station detail blob stored on the reservation.
// The details are informational, so a reservation booked before they were captured - or a malformed
// value - yields nil and is omitted from the response rather than failing the request.
func unmarshalStationInfo(detailsJSON []byte) *broker.StationInfo {
	if len(detailsJSON) == 0 {
		return nil
	}

	var details broker.StationInfo
	if err := json.Unmarshal(detailsJSON, &details); err != nil {
		rlog.Error("failed to unmarshal station details", "error", err)
		return nil
	}

	return &details
}

// reservationPriceDetails holds the calculated price details for a reservation.
type reservationPriceDetails struct {
	carFullPrice   int
	erpPrice       int
	discountAmount int
}

// calculatePriceDetails calculates the price details for a reservation based on the given parameters.
func calculatePriceDetails(reservation db.Reservation) reservationPriceDetails {
	pp := dbadapters.NumericToFloat64(reservation.PurchasePrice)
	mp := dbadapters.NumericToFloat64(reservation.MarkupPercentage)
	bErp := dbadapters.NumericToFloat64(reservation.BrokerErpPrice)
	btErp := dbadapters.NumericToFloat64(reservation.BtErpPrice)

	carFullPrice := pricing.ApplyMarkup(pp, mp)
	erpFullPrice := pricing.ApplyMarkup(bErp, mp) + btErp
	discountAmount := (erpFullPrice + carFullPrice) - dbadapters.NumericToFloat64(reservation.TotalPrice)

	return reservationPriceDetails{
		carFullPrice:   pricing.RoundToInt(carFullPrice),
		erpPrice:       pricing.RoundToInt(erpFullPrice),
		discountAmount: pricing.RoundToInt(discountAmount),
	}
}

func unmarshalJsons(carDetailsJson, payAtPickupJson []byte) (broker.CarDetails, actions.PayAtPickup, error) {
	var carDetails broker.CarDetails
	var payAtPickup actions.PayAtPickup
	if err := json.Unmarshal(carDetailsJson, &carDetails); err != nil {
		return carDetails, payAtPickup, err
	}

	if err := json.Unmarshal(payAtPickupJson, &payAtPickup); err != nil {
		return carDetails, payAtPickup, err
	}

	return carDetails, payAtPickup, nil
}
