package booking_handlers

import (
	"context"
	"errors"
	"strconv"

	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/booking/db"
	availability "encore.app/services/booking/handlers/availability"
	"encore.app/services/reservation"
	"encore.dev/rlog"
)

type CustomerBookParams struct {
	BookParams
	UserID int64
}

func (s *BookingService) CustomerBook(ctx context.Context, p CustomerBookParams) (*BookResponse, error) {
	snapshot, err := s.getSnapshot(ctx, p.SnapshotID)
	if err != nil {
		return nil, err
	}

	plan, err := findPlan(snapshot, p.RateQualifier, p.SupplierCode, p.PlanID)
	if err != nil {
		return nil, err
	}

	confID, err := bookCarAtBroker(snapshot, plan, p.BookParams)
	if err != nil {
		if errors.Is(err, broker.ErrFlightNumberRequired) {
			rlog.Error("booking failed due to missing flight number", "snapshotID", snapshot.ID, "rateQualifier", p.RateQualifier, "supplierCode", p.SupplierCode, "planID", p.PlanID, "error", err)
			return nil, ErrBookingFailed
		}
		return nil, ErrBookingFailed
	}

	reservationReq := s.buildCreateCustomerReservationParams(snapshot, plan, p, confID)
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

func (s *BookingService) buildCreateCustomerReservationParams(
	snapshot db.AvailablePlansSnapshot,
	plan availability.PlanPriceDetails,
	p CustomerBookParams,
	confirmationNumber string,
) reservation.CreateReservationParams {
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

	return reservation.CreateReservationParams{
		UserID:              p.UserID,
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
		DiscountPercentage:  plan.DiscountPercentage,
		BrokerErpPrice:      brokerErpPrice,
		BtErpPrice:          btErpPrice,
		PickupTime:          snapshot.PickupTime,
		DropoffTime:         snapshot.DropoffTime,
		PickupLocationName:  pickupLocName,
		DropoffLocationName: dropoffLocName,
		FlightNumber:        p.FlightNumber,
		PayAtPickup:         GetPayAtPickup(p.SelectedAddOns, plan),
		Excess:              plan.Excess,
		ExcessCurrency:      plan.ExcessCurrency,
		PickupLocationCode:  plan.PickupLocationCode,
		DropoffLocationCode: plan.DropoffLocationCode,
		SupplierTerms:       supplierInfo.TermsAndConditions,
		PickupDetails:       supplierInfo.PickupDetails,
		DropoffDetails:      supplierInfo.DropoffDetails,
	}
}
