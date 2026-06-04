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
	ID                  int64               `json:"id"`
	BrokerReservationID string              `json:"brokerReservationId"`
	ReservationStatus   string              `json:"reservationStatus"`
	PaymentStatus       string              `json:"paymentStatus"`
	CarDetails          broker.CarDetails   `json:"carDetails"`
	PlanInclusions      []string            `json:"planInclusions"`
	CurrencyCode        string              `json:"currencyCode"`
	CurrencyRate        float64             `json:"currencyRate"`
	CarFullPrice        int                 `json:"priceBefDesc"`
	DiscountAmount      int                 `json:"discountAmount"`
	ErpPrice            int                 `json:"erpPrice"`
	TotalPrice          int32               `json:"totalPrice"`
	PayAtPickup         actions.PayAtPickup `json:"payAtPickup"`
	FlightNumber        *string             `json:"flightNumber,omitempty" encore:"optional"`
	PickupLocationName  string              `json:"pickupLocationName"`
	DropoffLocationName string              `json:"dropoffLocationName"`
	PickupDate          string              `json:"pickupDate"`
	DropoffDate         string              `json:"dropoffDate"`
	PickupTime          string              `json:"pickupTime"`
	DropoffTime         string              `json:"dropoffTime"`
	RentalDays          int32               `json:"rentalDays"`
	DriverTitle         string              `json:"driverTitle"`
	DriverFirstName     string              `json:"driverFirstName"`
	DriverLastName      string              `json:"driverLastName"`
	DriverAge           int32               `json:"driverAge"`
	Voucher             *string             `json:"voucher,omitempty" encore:"optional"`
	VoucheredAt         *string             `json:"voucheredAt,omitempty" encore:"optional"`
	CreatedAt           string              `json:"createdAt"`
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
		TotalPrice:          row.TotalPrice,
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
	}, nil
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
	btErp := float64(reservation.BtErpPrice)

	carFullPrice := pricing.RoundToInt(pricing.ApplyMarkup(pp, mp))
	erpFullPrice := pricing.RoundToInt(pricing.ApplyMarkup(bErp, mp) + btErp)
	discountAmount := (erpFullPrice + carFullPrice) - int(reservation.TotalPrice)

	return reservationPriceDetails{
		carFullPrice:   carFullPrice,
		erpPrice:       erpFullPrice,
		discountAmount: discountAmount,
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
