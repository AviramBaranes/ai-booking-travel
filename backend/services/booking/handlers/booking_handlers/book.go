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
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation"
	"encore.app/services/reservation/handlers/actions"
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
	FlightNumber    *string              `json:"flightNumber" encore:"optional" validate:"omitempty"`
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

	plan, err := findPlan(snapshot, p.RateQualifier, p.SupplierCode, p.PlanID)
	if err != nil {
		return nil, err
	}

	confID, err := bookCarAtBroker(snapshot, plan, p)
	if err != nil {
		if errors.Is(err, broker.ErrFlightNumberRequired) {
			return nil, errFlightNumberRequired
		}
		return nil, ErrBookingFailed
	}

	reservationReq := s.buildCreateReservationParams(snapshot, plan, p, confID)
	rID, err := s.createReservation(ctx, reservationReq)
	if err != nil {
		return nil, err
	}

	err = s.query.DeleteSnapshotByID(ctx, snapshot.ID)
	if err != nil {
		rlog.Error("failed to delete snapshot after successful booking", "snapshotID", snapshot.ID, "error", err)
	}

	return &BookResponse{ReservationID: rID}, nil
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

	var btErpPrice, brokerErpPrice float64
	if p.IncludeERP {
		btErpPrice = plan.ChargedERPPriceWithVat
		brokerErpPrice = plan.SupplierErpPrice
	}

	pickupLocName, dropoffLocName, err := s.getLocationsNames(context.Background(), plan.PickupLocationCode, plan.DropoffLocationCode)
	if err != nil {
		rlog.Error("failed to get location names for reservation request", "error", err)
	}

	supplierInfo := FindSupplierInfo(snapshot, plan)

	var officeID, organizationID *int64
	var isOrganizationOrganic *bool
	if authData.OrganizationContext != nil {
		officeID = &authData.OrganizationContext.OfficeID
		organizationID = &authData.OrganizationContext.OrganizationID
		isOrganizationOrganic = &authData.OrganizationContext.IsOrganic
	}

	return reservation.CreateReservationParams{
		UserID:                authData.UserID,
		OfficeID:              officeID,
		OrganizationID:        organizationID,
		AdminRefID:            authData.AdminRefID,
		IsOrganizationOrganic: isOrganizationOrganic,
		BrokerReservationID:   confirmationNumber,
		Broker:                string(plan.Broker),
		SupplierCode:          plan.SupplierCode,
		CarDetails:            &plan.CarDetails,
		PlanInclusions:        plan.Inclusions,
		PickupDate:            dbadapters.DateToString(snapshot.PickupDate),
		DropoffDate:           dbadapters.DateToString(snapshot.DropoffDate),
		RentalDays:            rentalDays,
		DriverTitle:           p.DriverTitle,
		DriverFirstName:       p.DriverFirstName,
		DriverLastName:        p.DriverLastName,
		DriverAge:             driverAge,
		CountryCode:           snapshot.CountryCode,
		CurrencyCode:          plan.CurrencyCode,
		CurrencyRate:          plan.CurrencyRate,
		PurchasePrice:         plan.CarPurchasePrice,
		MarkupPercentage:      plan.MarkupPercentage,
		DiscountPercentage:    plan.DiscountPercentage,
		BrokerErpPrice:        brokerErpPrice,
		BtErpPrice:            btErpPrice,
		PickupTime:            snapshot.PickupTime,
		DropoffTime:           snapshot.DropoffTime,
		PickupLocationName:    pickupLocName,
		DropoffLocationName:   dropoffLocName,
		FlightNumber:          p.FlightNumber,
		PayAtPickup:           GetPayAtPickup(p.SelectedAddOns, plan),
		Excess:                plan.Excess,
		ExcessCurrency:        plan.ExcessCurrency,
		PickupLocationCode:    plan.PickupLocationCode,
		DropoffLocationCode:   plan.DropoffLocationCode,
		SupplierTerms:         supplierInfo.TermsAndConditions,
		PickupDetails:         supplierInfo.PickupDetails,
		DropoffDetails:        supplierInfo.DropoffDetails,
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

func findPlan(snapshot db.AvailablePlansSnapshot, rateQualifier, supplierCode, planID string) (availability.PlanPriceDetails, error) {
	var plans []availability.PlanPriceDetails
	if err := json.Unmarshal(snapshot.Plans, &plans); err != nil {
		rlog.Error("failed to unmarshal plans JSON", "error", err)
		return availability.PlanPriceDetails{}, api_errors.ErrInternalError
	}

	for _, plan := range plans {
		currPlanID := strconv.FormatInt(int64(plan.PlanID), 10)
		if plan.RateQualifier == rateQualifier && plan.SupplierCode == supplierCode && currPlanID == planID {
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

	var flightNumber string
	if p.FlightNumber != nil {
		flightNumber = *p.FlightNumber
	}
	res, err := b.Book(broker.BookingParams{
		RateQualifier:   plan.RateQualifier,
		SupplierCode:    plan.SupplierCode,
		Acriss:          plan.CarDetails.FullAcriss,
		PlanID:          p.PlanID,
		PickupLocation:  plan.PickupLocationCode,
		DropoffLocation: plan.DropoffLocationCode,
		IncludeERP:      p.IncludeERP,
		SelectedAddOns:  p.SelectedAddOns,
		DriverTitle:     p.DriverTitle,
		DriverFirstName: p.DriverFirstName,
		DriverLastName:  p.DriverLastName,
		FlightNumber:    flightNumber,
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

func GetPayAtPickup(selectedAddOnsReq []broker.SelectAddOn, plan availability.PlanPriceDetails) reservation.PayAtPickup {
	pap := reservation.PayAtPickup{
		Deposit:         plan.Deposit,
		DepositCurrency: plan.DepositCurrency,
		Fees:            plan.Fees,
	}

	selectedAddOns := make([]reservation.SelectedAddon, 0, len(selectedAddOnsReq))
	for _, selected := range selectedAddOnsReq {
		var addon broker.AddOn
		for _, planAddOn := range plan.AvailableAddOns {
			if planAddOn.ID == selected.ID {
				addon = planAddOn
				break
			}
		}

		selectedAddOns = append(selectedAddOns, reservation.SelectedAddon{
			ID:       addon.ID,
			Name:     addon.Name,
			Price:    addon.Price,
			Quantity: selected.Quantity,
		})
	}

	pap.SelectedAddons = selectedAddOns
	return pap
}

func (s *BookingService) createReservation(
	ctx context.Context,
	reservationReq actions.CreateReservationParams,
) (int64, error) {
	res, err := reservation.CreateReservation(ctx, reservationReq)
	if err != nil {
		rlog.Error("failed to create reservation after successful booking",
			"confirmationNumber", reservationReq.BrokerReservationID, "error", err)
		if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
			Subject: "Reservation creation failed after successful booking",
			Message: fmt.Sprintf("failed to create reservation after successful booking, confirmationNumber: %s, error: %v", reservationReq.BrokerReservationID, err),
		}); publishErr != nil {
			rlog.Error("failed to publish critical error email event", "confirmationNumber", reservationReq.BrokerReservationID, "error", publishErr)
		}
		return 0, ErrReservationCreationFailed
	}

	return res.ID, nil
}
