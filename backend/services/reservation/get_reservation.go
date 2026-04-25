package reservation

import (
	"context"
	"encoding/json"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/pricing"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type GetReservationResponse struct {
	ID                  int64             `json:"id"`
	BrokerReservationID string            `json:"brokerReservationId"`
	ReservationStatus   string            `json:"reservationStatus"`
	PaymentStatus       string            `json:"paymentStatus"`
	CarDetails          broker.CarDetails `json:"carDetails"`
	PlanInclusions      []string          `json:"planInclusions"`
	CurrencyCode        string            `json:"currencyCode"`
	CurrencyRate        float64           `json:"currencyRate"`
	CarFullPrice        int               `json:"priceBefDesc"`
	DiscountAmount      int               `json:"discountAmount"`
	ErpPrice            int               `json:"erpPrice"`
	TotalPrice          int32             `json:"totalPrice"`
	PickupLocationName  string            `json:"pickupLocationName"`
	DropoffLocationName string            `json:"dropoffLocationName"`
	PickupDate          string            `json:"pickupDate"`
	ReturnDate          string            `json:"returnDate"`
	PickupTime          string            `json:"pickupTime"`
	DropoffTime         string            `json:"dropoffTime"`
	RentalDays          int32             `json:"rentalDays"`
	DriverTitle         string            `json:"driverTitle"`
	DriverFirstName     string            `json:"driverFirstName"`
	DriverLastName      string            `json:"driverLastName"`
	DriverAge           int32             `json:"driverAge"`
	Voucher             *string           `json:"voucher,omitempty" encore:"optional"`
	VoucheredAt         *string           `json:"voucheredAt,omitempty" encore:"optional"`
	CreatedAt           string            `json:"createdAt"`
}

// encore:api auth method=GET path=/reservations/:id
func (s *Service) GetReservation(ctx context.Context, id int64) (*GetReservationResponse, error) {
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

	var carDetails broker.CarDetails
	if err := json.Unmarshal(row.CarDetails, &carDetails); err != nil {
		rlog.Error("failed to unmarshal car details", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	rpd := calculatePriceDetails(row)

	voucheredAt := db.TimestamptzToString(row.VoucheredAt)
	return &GetReservationResponse{
		ID:                  row.ID,
		BrokerReservationID: row.BrokerReservationID,
		ReservationStatus:   string(row.ReservationStatus),
		PaymentStatus:       string(row.PaymentStatus),
		CarDetails:          carDetails,
		PlanInclusions:      row.PlanInclusions,
		CurrencyCode:        row.CurrencyCode,
		CurrencyRate:        db.NumericToFloat64(row.CurrencyRate),
		CarFullPrice:        rpd.carFullPrice,
		ErpPrice:            rpd.erpPrice,
		DiscountAmount:      rpd.discountAmount,
		TotalPrice:          row.TotalPrice,
		PickupDate:          db.DateToString(row.PickupDate),
		ReturnDate:          db.DateToString(row.ReturnDate),
		PickupTime:          row.PickupTime,
		DropoffTime:         row.DropoffTime,
		RentalDays:          row.RentalDays,
		DriverTitle:         row.DriverTitle,
		DriverFirstName:     row.DriverFirstName,
		DriverLastName:      row.DriverLastName,
		DriverAge:           row.DriverAge,
		CreatedAt:           db.TimestamptzToString(row.CreatedAt),
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
func calculatePriceDetails(reservation db.GetReservationByIDRow) reservationPriceDetails {
	pp := db.NumericToFloat64(reservation.PurchasePrice)
	mp := db.NumericToFloat64(reservation.MarkupPercentage)
	bErp := db.NumericToFloat64(reservation.BrokerErpPrice)
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
