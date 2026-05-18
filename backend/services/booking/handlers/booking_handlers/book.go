package booking_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	availability "encore.app/services/booking/handlers/availability"
	"encore.app/services/notifications"
	"encore.app/services/reservation"
	"encore.dev/rlog"
)

type BookParams struct {
	SnapshotID      int64                `json:"snapshotId" validate:"required"`
	RateQualifier   string               `json:"rateQualifier" validate:"required"`
	SupplierCode    string               `json:"supplierCode" validate:"required"`
	PlanID          string               `json:"planId"`
	IncludeERP      bool                 `json:"includeERP"`
	SelectedAddOns  []broker.SelectAddOn `json:"selectedAddOns"`
	DriverTitle     string               `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Mrs' 'Ms' 'Miss' 'Dr'"`
	DriverFirstName string               `json:"driverFirstName" validate:"required,uppercase_only"`
	DriverLastName  string               `json:"driverLastName" validate:"required,uppercase_only"`
	FlightNumber    string               `json:"flightNumber" encore:"optional"`
}

func (p BookParams) Validate() error {
	return validation.ValidateStruct(p)
}

type BookResponse struct {
	ReservationID int64 `json:"reservationId"`
}

func (s *BookingService) Book(ctx context.Context, p BookParams) (*BookResponse, error) {
	snapshot, err := s.getSnapshot(ctx, p.SnapshotID)
	if err != nil {
		return nil, err
	}

	plan, err := findPlan(snapshot, p.RateQualifier, p.SupplierCode)
	if err != nil {
		return nil, err
	}

	confID, err := bookCarAtBroker(snapshot, plan, p)
	if err != nil {
		if errors.Is(err, broker.ErrFlightNumberRequired) {
			return nil, errFlightNumberRequired
		}
		return nil, errBookingFailed
	}

	reservationReq := s.buildCreateReservationParams(snapshot, plan, p, confID)
	res, err := reservation.CreateReservation(ctx, reservationReq)
	if err != nil {
		rlog.Error("failed to create reservation after successful booking",
			"confirmationNumber", confID, "error", err)
		notifications.CriticalErrorEventTopic.Publish(ctx, &notifications.CriticalErrorEvent{
			Subject: "Reservation creation failed after successful booking",
			Message: fmt.Sprintf("failed to create reservation after successful booking, confirmationNumber: %s, error: %v", confID, err),
		})
		return nil, errReservationCreationFailed
	}

	err = s.query.DeleteSnapshotByID(ctx, snapshot.ID)
	if err != nil {
		rlog.Error("failed to delete snapshot after successful booking", "snapshotID", snapshot.ID, "error", err)
	}

	return &BookResponse{ReservationID: res.ID}, nil
}

func (s *BookingService) buildCreateReservationParams(
	snapshot db.AvailablePlansSnapshot,
	plan availability.PlanPriceDetails,
	p BookParams,
	confirmationNumber string,
) reservation.CreateReservationParams {
	authData := auth.GetAuthData()
	rentalDays, _ := calculateSnapshotRentalDays(snapshot)
	driverAge, _ := strconv.Atoi(snapshot.DriverAge)

	var btErpPrice int
	var brokerErpPrice float64
	if p.IncludeERP {
		btErpPrice = plan.ChargedERPPriceWithVat
		brokerErpPrice = plan.SupplierErpPrice
	}

	pickupLocName, dropoffLocName, err := s.getLocationsNames(context.Background(), plan.PickupLocationCode, plan.DropoffLocationCode)
	if err != nil {
		rlog.Error("failed to get location names for reservation request", "error", err)
	}

	return reservation.CreateReservationParams{
		UserID:              authData.UserID,
		BrokerReservationID: confirmationNumber,
		Broker:              string(plan.Broker),
		SupplierCode:        plan.SupplierCode,
		CarDetails:          &plan.CarDetails,
		PlanInclusions:      plan.Inclusions,
		PickupDate:          dbadapters.DateToString(snapshot.PickupDate),
		DropoffDate:         dbadapters.DateToString(snapshot.DropoffDate),
		RentalDays:          rentalDays,
		DriverTitle:         p.DriverTitle,
		DriverFirstName:     p.DriverFirstName,
		DriverLastName:      p.DriverLastName,
		DriverAge:           driverAge,
		CountryCode:         snapshot.CountryCode,
		CurrencyCode:        plan.CurrencyCode,
		CurrencyRate:        plan.CurrencyRate,
		PurchasePrice:       plan.CarPurchasePrice,
		MarkupPercentage:    plan.MarkupPercentage,
		DiscountPercentage:  int(plan.DiscountPercentage),
		BrokerErpPrice:      brokerErpPrice,
		BtErpPrice:          btErpPrice,
		PickupTime:          snapshot.PickupTime,
		DropoffTime:         snapshot.DropoffTime,
		PickupLocationName:  pickupLocName,
		DropoffLocationName: dropoffLocName,
	}
}

func calculateSnapshotRentalDays(snapshot db.AvailablePlansSnapshot) (int, error) {
	return broker.CalculateDaysCount(
		dbadapters.DateToString(snapshot.PickupDate), snapshot.PickupTime,
		dbadapters.DateToString(snapshot.DropoffDate), snapshot.DropoffTime,
	)
}

func (s *BookingService) getSnapshot(ctx context.Context, snapshotID int64) (db.AvailablePlansSnapshot, error) {
	snapshot, err := s.query.GetSnapshotByID(ctx, snapshotID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return db.AvailablePlansSnapshot{}, errSnapshotNotFound
		}
		rlog.Error("failed to get snapshot", "snapshotID", snapshotID, "error", err)
		return db.AvailablePlansSnapshot{}, api_errors.ErrInternalError
	}

	return snapshot, nil
}

func findPlan(snapshot db.AvailablePlansSnapshot, rateQualifier, supplierCode string) (availability.PlanPriceDetails, error) {
	var plans []availability.PlanPriceDetails
	if err := json.Unmarshal(snapshot.Plans, &plans); err != nil {
		rlog.Error("failed to unmarshal plans JSON", "error", err)
		return availability.PlanPriceDetails{}, api_errors.ErrInternalError
	}

	for _, plan := range plans {
		if plan.RateQualifier == rateQualifier && plan.SupplierCode == supplierCode {
			return plan, nil
		}
	}
	return availability.PlanPriceDetails{}, errPlanNotFound
}

func bookCarAtBroker(snapshot db.AvailablePlansSnapshot, plan availability.PlanPriceDetails, p BookParams) (string, error) {
	b, err := getBrokerByPlan(plan)
	if err != nil {
		rlog.Error("failed to get broker for plan", "RateQualifier", plan.RateQualifier, "error", err)
		return "", err
	}

	res, err := b.Book(broker.BookingParams{
		RateQualifier:   plan.RateQualifier,
		SupplierCode:    plan.SupplierCode,
		Acriss:          plan.CarDetails.Acriss,
		PlanID:          p.PlanID,
		PickupLocation:  plan.PickupLocationCode,
		DropoffLocation: plan.DropoffLocationCode,
		IncludeERP:      p.IncludeERP,
		SelectedAddOns:  p.SelectedAddOns,
		DriverTitle:     p.DriverTitle,
		DriverFirstName: p.DriverFirstName,
		DriverLastName:  p.DriverLastName,
		FlightNumber:    p.FlightNumber,
		DriverAge:       snapshot.DriverAge,
		PickupDate:      dbadapters.DateToString(snapshot.PickupDate),
		DropoffDate:     dbadapters.DateToString(snapshot.DropoffDate),
		PickupTime:      snapshot.PickupTime,
		DropoffTime:     snapshot.DropoffTime,
		CountryCode:     snapshot.CountryCode,
	})
	if err != nil {
		rlog.Error("failed to book car at broker", "broker", b.Name(), "error", err)
		return "", err
	}

	return res.ConfirmationNumber, nil
}

func getBrokerByPlan(plan availability.PlanPriceDetails) (broker.Booker, error) {
	switch plan.Broker {
	case broker.BrokerHertz:
		return broker.NewHertz(), nil
	case broker.BrokerFlex:
		return broker.NewFlex(), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}

func (s *BookingService) getLocationsNames(ctx context.Context, pickupBrokerLocationID, dropoffBrokerLocationID string) (string, string, error) {
	pickupLoc, err := s.query.GetLocationByBrokerLocationID(ctx, pickupBrokerLocationID)
	if err != nil {
		rlog.Error("failed to get pickup location name", "brokerLocationID", pickupBrokerLocationID, "error", err)
		return "", "", fmt.Errorf("query pickup location %w", err)
	}

	if pickupBrokerLocationID == dropoffBrokerLocationID {
		return pickupLoc.Name, pickupLoc.Name, nil
	}

	dropoffLoc, err := s.query.GetLocationByBrokerLocationID(ctx, dropoffBrokerLocationID)
	if err != nil {
		rlog.Error("failed to get dropoff location name", "brokerLocationID", dropoffBrokerLocationID, "error", err)
		return "", "", fmt.Errorf("query dropoff location %w", err)
	}

	return pickupLoc.Name, dropoffLoc.Name, nil
}
