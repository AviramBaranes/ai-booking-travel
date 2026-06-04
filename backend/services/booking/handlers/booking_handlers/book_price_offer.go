package booking_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.app/services/notifications"
	"encore.app/services/reservation"
	"encore.dev/rlog"
)

type BookPriceOfferParams struct {
	PriceOfferID    int64   `json:"priceOfferId" validate:"required"`
	DriverTitle     string  `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Mrs' 'Ms' 'Miss' 'Dr'"`
	DriverFirstName string  `json:"driverFirstName" validate:"required,uppercase_only"`
	DriverLastName  string  `json:"driverLastName" validate:"required,uppercase_only"`
	FlightNumber    *string `json:"flightNumber" encore:"optional"`
}

func (p BookPriceOfferParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *BookingService) BookPriceOffer(ctx context.Context, p BookPriceOfferParams) (*BookResponse, error) {
	authData := auth.GetAuthData()

	offer, err := s.getBookablePriceOffer(ctx, p.PriceOfferID, authData.UserID)
	if err != nil {
		return nil, err
	}

	bookingRes, offerCarDetails, err := bookPriceOfferAtBroker(offer, p)
	if err != nil {
		return nil, err
	}

	reservationReq, err := buildPriceOfferReservationRequest(authData.UserID, offer, p, bookingRes.ConfirmationNumber, offerCarDetails)
	if err != nil {
		return nil, err
	}

	reservation, err := reservation.CreateReservation(ctx, reservationReq)
	if err != nil {
		rlog.Error("failed to create reservation after successful booking",
			"confirmationNumber", bookingRes.ConfirmationNumber, "error", err)
		notifications.PublishEmailEvent(ctx, notifications.EmailEventTypeCriticalError, notifications.CriticalErrorEmailPayload{
			Subject: "Reservation creation failed after successful booking",
			Message: fmt.Sprintf("failed to create reservation after successful booking, confirmationNumber: %s, error: %v", bookingRes.ConfirmationNumber, err),
		})
		return nil, errReservationCreationFailed
	}

	s.markPriceOfferBooked(ctx, p.PriceOfferID, reservation.ID, authData.UserID)

	return &BookResponse{ReservationID: reservation.ID}, nil
}

func (s *BookingService) getBookablePriceOffer(ctx context.Context, priceOfferID int64, agentID int64) (db.GetPriceOfferByIdRow, error) {
	offer, err := s.query.GetPriceOfferById(ctx, db.GetPriceOfferByIdParams{
		ID:      priceOfferID,
		AgentID: agentID,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return db.GetPriceOfferByIdRow{}, api_errors.ErrNotFound
		}
		rlog.Error("failed to get price offer for renewal", "id", priceOfferID, "error", err)
		return db.GetPriceOfferByIdRow{}, api_errors.ErrInternalError
	}

	if err := ensurePriceOfferRenewedRecently(offer, priceOfferID); err != nil {
		return db.GetPriceOfferByIdRow{}, err
	}

	return offer, nil
}

func ensurePriceOfferRenewedRecently(offer db.GetPriceOfferByIdRow, priceOfferID int64) error {
	if !offer.RenewedAt.Valid {
		rlog.Error("price offer has invalid renewed_at", "id", priceOfferID)
		return api_errors.ErrInternalError
	}
	if time.Since(offer.RenewedAt.Time) > 15*time.Minute {
		return errOldOffer
	}

	return nil
}

func bookPriceOfferAtBroker(offer db.GetPriceOfferByIdRow, p BookPriceOfferParams) (broker.BookingResponse, broker.CarDetails, error) {
	b, err := getBroker(offer.Broker)
	if err != nil {
		rlog.Error("failed to get broker for price offer booking", "error", err)
		return broker.BookingResponse{}, broker.CarDetails{}, api_errors.ErrInternalError
	}

	offerCarDetails, err := unmarshalPriceOfferCarDetails(offer)
	if err != nil {
		return broker.BookingResponse{}, broker.CarDetails{}, err
	}

	bookingRes, err := b.Book(buildPriceOfferBookingParams(offer, p, offerCarDetails))
	if err != nil {
		rlog.Error("failed to book car at broker", "broker", b.Name(), "error", err)
		return broker.BookingResponse{}, broker.CarDetails{}, err
	}

	return bookingRes, offerCarDetails, nil
}

func unmarshalPriceOfferCarDetails(offer db.GetPriceOfferByIdRow) (broker.CarDetails, error) {
	var offerCarDetails broker.CarDetails
	if err := json.Unmarshal(offer.CarDetails, &offerCarDetails); err != nil {
		rlog.Error("failed to unmarshal price offer car details", "id", offer.ID, "error", err)
		return broker.CarDetails{}, api_errors.ErrInternalError
	}

	return offerCarDetails, nil
}

func isPriceOfferErpIncluded(offer db.GetPriceOfferByIdRow) bool {
	return offer.BtErpPrice != 0 || dbadapters.NumericToFloat64(offer.BrokerErpPrice) != 0
}

func buildPriceOfferBookingParams(offer db.GetPriceOfferByIdRow, p BookPriceOfferParams, offerCarDetails broker.CarDetails) broker.BookingParams {
	var flightNumber string
	if p.FlightNumber != nil {
		flightNumber = *p.FlightNumber
	}
	return broker.BookingParams{
		RateQualifier:   offer.RateQualifier,
		SupplierCode:    offer.SupplierCode,
		Acriss:          offerCarDetails.Acriss,
		PlanID:          offer.PlanID,
		PickupLocation:  offer.PickupBrokerLocationID,
		DropoffLocation: offer.DropoffBrokerLocationID,
		IncludeERP:      isPriceOfferErpIncluded(offer),
		SelectedAddOns:  []broker.SelectAddOn{},
		DriverTitle:     p.DriverTitle,
		DriverFirstName: p.DriverFirstName,
		DriverLastName:  p.DriverLastName,
		FlightNumber:    flightNumber,
		DriverAge:       offer.DriverAge,
		PickupDate:      dbadapters.DateToString(offer.PickupDate),
		DropoffDate:     dbadapters.DateToString(offer.DropoffDate),
		PickupTime:      offer.PickupTime,
		DropoffTime:     offer.DropoffTime,
		CountryCode:     offer.CurrencyCode,
	}
}

func buildPriceOfferReservationRequest(
	userID int64,
	offer db.GetPriceOfferByIdRow,
	p BookPriceOfferParams,
	confirmationNumber string,
	offerCarDetails broker.CarDetails,
) (reservation.CreateReservationParams, error) {
	driverAge, err := strconv.Atoi(offer.DriverAge)
	if err != nil {
		rlog.Error("failed to convert driver age to int", "driverAge", offer.DriverAge, "error", err)
		return reservation.CreateReservationParams{}, api_errors.ErrInternalError
	}

	authData := auth.GetAuthData()

	var officeID, organizationID *int64
	var isOrganizationOrganic *bool
	if authData.OrganizationContext != nil {
		officeID = &authData.OrganizationContext.OfficeID
		organizationID = &authData.OrganizationContext.OrganizationID
		isOrganizationOrganic = &authData.OrganizationContext.IsOrganic
	}

	return reservation.CreateReservationParams{
		UserID:                userID,
		OfficeID:              officeID,
		OrganizationID:        organizationID,
		AdminRefID:            authData.AdminRefID,
		IsOrganizationOrganic: isOrganizationOrganic,
		BrokerReservationID:   confirmationNumber,
		Broker:                string(offer.Broker),
		SupplierCode:          offer.SupplierCode,
		CarDetails:            &offerCarDetails,
		PlanInclusions:        offer.PlanInclusions,
		PickupDate:            dbadapters.DateToString(offer.PickupDate),
		DropoffDate:           dbadapters.DateToString(offer.DropoffDate),
		RentalDays:            int(offer.RentalDays),
		DriverTitle:           p.DriverTitle,
		DriverFirstName:       p.DriverFirstName,
		DriverLastName:        p.DriverLastName,
		DriverAge:             driverAge,
		CountryCode:           offer.CountryCode,
		CurrencyCode:          offer.CurrencyCode,
		CurrencyRate:          dbadapters.NumericToFloat64(offer.CurrencyRate),
		PurchasePrice:         dbadapters.NumericToFloat64(offer.PurchasePrice),
		MarkupPercentage:      dbadapters.NumericToFloat64(offer.MarkupPercentage),
		DiscountPercentage:    0,
		BrokerErpPrice:        dbadapters.NumericToFloat64(offer.BrokerErpPrice),
		BtErpPrice:            int(offer.BtErpPrice),
		PickupTime:            offer.PickupTime,
		DropoffTime:           offer.DropoffTime,
		PickupLocationName:    offer.PickupLocation,
		DropoffLocationName:   offer.DropoffLocation,
		FlightNumber:          p.FlightNumber,
		PayAtPickup:           unmarshalPayAtPickup(offer.PayAtPickup),
	}, nil
}

func (s *BookingService) markPriceOfferBooked(ctx context.Context, priceOfferID int64, reservationID int64, agentID int64) {
	err := s.query.UpdatePriceOffer(ctx, db.UpdatePriceOfferParams{
		ID:            priceOfferID,
		AgentID:       agentID,
		ReservationID: &reservationID,
		Status:        db.NullOfferStatus{OfferStatus: db.OfferStatusBooked, Valid: true},
	})
	if err != nil {
		rlog.Error("failed to mark price offer as booked", "id", priceOfferID, "reservationID", reservationID, "error", err)
	}
}

// unmarshalPayAtPickup unmarshals the PayAtPickup JSON from the database into a struct for the response.
func unmarshalPayAtPickup(papJson []byte) reservation.PayAtPickup {
	var pap reservation.PayAtPickup
	if err := json.Unmarshal(papJson, &pap); err != nil {
		rlog.Error("failed to unmarshal pay at pickup", "error", err)
		return reservation.PayAtPickup{}
	}
	return pap
}
