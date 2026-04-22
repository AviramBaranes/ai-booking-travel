package booking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/internal/validation"
	auth "encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.app/services/notifications"
	"encore.app/services/reservation"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	errSnapshotNotFound = api_errors.NewErrorWithDetail(
		errs.NotFound,
		"Snapshot not found",
		api_errors.ErrorDetails{Code: api_errors.CodeSnapshotNotFound},
	)
	errPlanNotFound = api_errors.NewErrorWithDetail(
		errs.NotFound,
		"Plan not found",
		api_errors.ErrorDetails{Code: api_errors.CodePlanNotFound},
	)
	errBookingFailed = api_errors.NewErrorWithDetail(
		errs.Unknown,
		"Booking failed",
		api_errors.ErrorDetails{Code: api_errors.CodeBookingFailed},
	)
	errReservationCreationFailed = api_errors.NewErrorWithDetail(
		errs.Unknown,
		"Reservation creation failed",
		api_errors.ErrorDetails{Code: api_errors.CodeReservationCreationFailed},
	)

	errFlightNumberRequired = api_errors.NewErrorWithDetail(
		errs.InvalidArgument,
		"Flight number is required for this office",
		api_errors.ErrorDetails{Code: api_errors.CodeFlightNumberRequired},
	)
)

// BookResponse represents the response returned after a successful booking, including the booking reference number and any relevant details.
type BookRequest struct {
	SnapshotID      int64                `json:"snapshotId" validate:"required"`
	RateQualifier   string               `json:"rateQualifier" validate:"required"`
	SupplierCode    string               `json:"supplierCode" validate:"required"`
	PlanID          string               `json:"planId"`
	IncludeERP      bool                 `json:"includeERP"`
	SelectedAddOns  []broker.SelectAddOn `json:"selectedAddOns"`
	DriverTitle     string               `json:"driverTitle" validate:"required,notblank,oneof='Mr' 'Ms'"`
	DriverFirstName string               `json:"driverFirstName" validate:"required,uppercase_only"`
	DriverLastName  string               `json:"driverLastName" validate:"required,uppercase_only"`
	FlightNumber    string               `json:"flightNumber" encore:"optional"`
}

func (p BookRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type BookResponse struct {
	ReservationID int64 `json:"reservationId"`
}

//encore:api auth method=POST path=/booking tag:agent
func (s *Service) Book(ctx context.Context, params BookRequest) (*BookResponse, error) {
	snapshot, err := s.getSnapshot(ctx, params.SnapshotID)
	if err != nil {
		return nil, err
	}

	plan, err := findPlan(snapshot, params.RateQualifier, params.SupplierCode)
	if err != nil {
		return nil, err
	}

	confID, err := bookCarAtBroker(snapshot, plan, params)
	if err != nil {
		if errors.Is(err, broker.ErrFlightNumberRequired) {
			return nil, errFlightNumberRequired
		}
		return nil, errBookingFailed
	}

	reservationReq := s.buildCreateReservationRequest(snapshot, plan, params, confID)
	res, err := reservation.CreateReservation(ctx, reservationReq)
	if err != nil {
		rlog.Error("failed to create reservation after successful booking",
			"confirmationNumber", confID, "error", err)
		notifications.CriticalErrorEventTopic.Publish(ctx, &notifications.CriticalErrorEvent{
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

// buildCreateReservationRequest maps booking data into a CreateReservationRequest.
func (s *Service) buildCreateReservationRequest(
	snapshot db.AvailablePlansSnapshot,
	plan planPriceDetails,
	params BookRequest,
	confirmationNumber string,
) reservation.CreateReservationRequest {
	authData := auth.GetAuthData()

	rentalDays, _ := broker.CalculateDaysCount(
		snapshot.PickupDate, snapshot.PickupTime,
		snapshot.ReturnDate, snapshot.ReturnTime,
	)

	driverAge, _ := strconv.Atoi(snapshot.DriverAge)

	var btErpPrice int
	var brokerErpPrice float64
	if params.IncludeERP {
		btErpPrice = plan.ChargedERPPriceWithVat
		brokerErpPrice = plan.SupplierErpPrice
	}

	pickupLocName, dropoffLocName, err := s.getLocationsNames(context.Background(), plan.PickupLocationCode, plan.DropoffLocationCode)
	if err != nil {
		rlog.Error("failed to get location names for reservation request", "error", err)
	}

	return reservation.CreateReservationRequest{
		UserID:              authData.UserID,
		BrokerReservationID: confirmationNumber,
		Broker:              string(plan.Broker),
		SupplierCode:        plan.SupplierCode,
		CarDetails:          &plan.CarDetails,
		PlanInclusions:      plan.Inclusions,
		PickupDate:          snapshot.PickupDate,
		ReturnDate:          snapshot.ReturnDate,
		RentalDays:          rentalDays,
		DriverTitle:         params.DriverTitle,
		DriverFirstName:     params.DriverFirstName,
		DriverLastName:      params.DriverLastName,
		DriverAge:           driverAge,
		CountryCode:         snapshot.CountryCode,
		CurrencyCode:        plan.CurrencyCode,
		CurrencyRate:        plan.CurrencyRate,
		PurchasePrice:       plan.CarPurchasePrice,
		MarkupPercentage:    plan.MarkupPercentage,
		DiscountPercentage:  plan.DiscountPercentage,
		BrokerErpPrice:      brokerErpPrice,
		BtErpPrice:          btErpPrice,
		PickupTime:          snapshot.PickupTime,
		DropoffTime:         snapshot.ReturnTime,
		PickupLocationName:  pickupLocName,
		DropoffLocationName: dropoffLocName,
	}
}

// getSnapshot retrieves the snapshot row for the given snapshot ID.
func (s *Service) getSnapshot(ctx context.Context, snapshotID int64) (db.AvailablePlansSnapshot, error) {
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

// findPlan finds the plan in the snapshot that matches the given rate qualifier and supplier code.
func findPlan(snapshot db.AvailablePlansSnapshot, rateQualifier, supplierCode string) (planPriceDetails, error) {
	var plans []planPriceDetails
	if err := json.Unmarshal(snapshot.Plans, &plans); err != nil {
		rlog.Error("failed to unmarshal plans JSON", "error", err)
		return planPriceDetails{}, api_errors.ErrInternalError
	}

	for _, plan := range plans {
		if plan.RateQualifier == rateQualifier && plan.SupplierCode == supplierCode {
			return plan, nil
		}
	}
	return planPriceDetails{}, errPlanNotFound
}

// bookCarAtBroker performs the actual booking with the broker using the provided plan details and booking request parameters.
func bookCarAtBroker(snapshot db.AvailablePlansSnapshot, plan planPriceDetails, params BookRequest) (string, error) {
	b, err := getBrokerByPlan(plan)
	if err != nil {
		rlog.Error("failed to get broker for plan", "RateQualifier", plan.RateQualifier, "error", err)
		return "", err
	}

	res, err := b.Book(broker.BookingParams{
		RateQualifier:   plan.RateQualifier,
		SupplierCode:    plan.SupplierCode,
		Acriss:          plan.CarDetails.Acriss,
		PlanID:          params.PlanID,
		PickupLocation:  plan.PickupLocationCode,
		DropoffLocation: plan.DropoffLocationCode,
		IncludeERP:      params.IncludeERP,
		SelectedAddOns:  params.SelectedAddOns,
		DriverTitle:     params.DriverTitle,
		DriverFirstName: params.DriverFirstName,
		DriverLastName:  params.DriverLastName,
		FlightNumber:    params.FlightNumber,
		DriverAge:       snapshot.DriverAge,
		PickupDate:      snapshot.PickupDate,
		DropoffDate:     snapshot.ReturnDate,
		PickupTime:      snapshot.PickupTime,
		DropoffTime:     snapshot.ReturnTime,
		CountryCode:     snapshot.CountryCode,
	})
	if err != nil {
		rlog.Error("failed to book car at broker", "broker", b.Name(), "error", err)
		return "", err
	}

	return res.ConfirmationNumber, nil
}

// getBrokerByPlan returns the broker implementation based on the broker specified in the plan details.
func getBrokerByPlan(plan planPriceDetails) (broker.Booker, error) {
	switch plan.Broker {
	case broker.BrokerHertz:
		return broker.NewHertz(), nil
	case broker.BrokerFlex:
		return broker.NewFlex(), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}

// getLocationsNames retrieves the pickup and dropoff location names based on the broker location IDs in the reservation request.
func (s *Service) getLocationsNames(ctx context.Context, pickupBrokerLocationID, dropoffBrokerLocationID string) (string, string, error) {
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
