package broker

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"encore.dev/rlog"
)

const (
	// hertzGRCommitted is the status of a confirmed reservation. Anything else, such as "On Hold"
	// while Autohellas' reservations desk decides, is not what the customer paid for.
	hertzGRCommitted = "Committed"

	// hertzGRAlreadyCancelled is the error code for a reservation that is already cancelled.
	hertzGRAlreadyCancelled = "95"

	// hertzGRPriceTolerance absorbs rounding between the quoted and the booked total.
	hertzGRPriceTolerance = 0.01
)

// Book reserves the car group the plan quoted. Autohellas has no quote to lock the price, so the
// reservation is priced again; one that comes back unconfirmed, or dearer than quoted, is cancelled
// and the booking fails.
func (h *HertzGR) Book(p BookingParams) (BookingResponse, error) {
	rateCode, group, ok := strings.Cut(p.RateQualifier, ":")
	if !ok || rateCode == "" || group == "" {
		return BookingResponse{}, fmt.Errorf("hertz gr book: malformed rate qualifier %q", p.RateQualifier)
	}

	ref, err := hertzGRReference()
	if err != nil {
		return BookingResponse{}, err
	}

	rq := hertzGRVehResRQ{
		Xmlns:   otaNamespace,
		Version: "1.0",
		POS:     hertzGRRequestPOS(),
		Core: hertzGRVehResRQCore{
			Status: "All",
			RentalCore: hertzGRVehRentalCore{
				PickUpDateTime: p.PickupDate + "T" + p.PickupTime + ":00",
				ReturnDateTime: p.DropoffDate + "T" + p.DropoffTime + ":00",
				PickUpLocation: hertzGRLocationCodeRef{LocationCode: p.PickupLocation},
				ReturnLocation: hertzGRLocationCodeRef{LocationCode: p.DropoffLocation},
			},
			Customer: hertzGRCustomer{PersonName: hertzGRPersonName{
				GivenName: strings.TrimSpace(p.DriverFirstName),
				Surname:   strings.TrimSpace(p.DriverLastName),
			}},
			VehPref:       hertzGRVehPref{VendorCarType: group},
			RateQualifier: hertzGRRateQualifier{RateQualifier: rateCode},
		},
	}
	// Our reference goes in a second POS source, as Autohellas requires on every reservation.
	rq.POS.Sources = append(rq.POS.Sources, hertzGRSource{
		ISOCountry:  hertzGRISOCountry,
		RequestorID: hertzGRRequestorID{Type: "16", ID: ref},
	})
	if age, err := strconv.Atoi(p.DriverAge); err == nil && age > 0 {
		rq.Core.DriverType = &hertzGRDriverType{Age: age}
	}
	if flight := strings.TrimSpace(p.FlightNumber); flight != "" {
		rq.Info = &hertzGRVehResRQInfo{ArrivalDetails: &hertzGRArrivalDetails{TransportationCode: "14", Number: flight}}
	}

	var rs hertzGRVehResRS
	if err := h.call("VehicleReservation", rq, &rs); err != nil {
		return BookingResponse{}, err
	}
	if len(rs.Errors) > 0 {
		e := rs.Errors[0]
		return BookingResponse{}, fmt.Errorf("hertz gr reservation returned %s: %s", e.Code, e.message())
	}

	seg := rs.Core.Segment
	confID := strings.TrimSpace(seg.ConfID.ID)
	if confID == "" {
		return BookingResponse{}, fmt.Errorf("hertz gr reservation %s came back %q without a reservation number", ref, rs.Core.ReservationStatus)
	}

	reason := ""
	switch {
	case rs.Core.ReservationStatus != hertzGRCommitted:
		reason = fmt.Sprintf("came back %q", rs.Core.ReservationStatus)
	case seg.TotalCharge.EstimatedTotalAmount <= 0:
		rlog.Warn("hertz gr reservation has no total to check against the quote", "confID", confID)
	default:
		if booked := hertzGRPrepaid(seg.TotalCharge, seg.Fees); booked > p.QuotedPrice+hertzGRPriceTolerance {
			reason = fmt.Sprintf("priced %.2f %s, above the %.2f quoted", booked, seg.TotalCharge.CurrencyCode, p.QuotedPrice)
		}
	}

	if reason == "" {
		return BookingResponse{ConfirmationNumber: confID}, nil
	}

	if err := h.cancel(confID, p.DriverFirstName, p.DriverLastName); err != nil {
		return BookingResponse{}, fmt.Errorf("%w: hertz gr reservation %s %s and could not be cancelled: %v", ErrBookingLeftOpen, confID, reason, err)
	}
	return BookingResponse{}, fmt.Errorf("hertz gr reservation %s %s and was cancelled", confID, reason)
}

// hertzGRPrepaid is what we pay for a car: the grand total less what the customer settles at the
// station. Search prices the plan the same way, so the two compare.
func hertzGRPrepaid(total hertzGRTotalCharge, fees []hertzGRCharge) float64 {
	payOnArrival := 0.0
	for _, f := range fees {
		if f.Purpose == hertzGRPayOnArrival {
			payOnArrival += f.Amount
		}
	}
	return total.EstimatedTotalAmount - payOnArrival
}

// hertzGRReference is our reference for a reservation, at most the 20 characters Autohellas allows.
// It is random: a retried booking is a new reservation, and the reference only has to be unique.
func hertzGRReference() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("hertz gr book: generate reference: %w", err)
	}
	return "AIBT" + strings.ToUpper(hex.EncodeToString(b)), nil
}

// Cancel cancels a reservation by Autohellas' reservation number. Autohellas checks the driver's
// surname against the reservation; a reservation already cancelled counts as cancelled.
func (h *HertzGR) Cancel(bookingID, lastName, supplierCode string) error {
	return h.cancel(bookingID, "", lastName)
}

func (h *HertzGR) cancel(confID, givenName, surname string) error {
	var rq hertzGRVehCancelRQ
	rq.Xmlns = otaNamespace
	rq.Version = "1.0"
	rq.POS = hertzGRRequestPOS()
	rq.Core.CancelType = "Cancel"
	rq.Core.UniqueID = hertzGRUniqueID{Type: "14", ID: confID}
	rq.Core.PersonName = hertzGRPersonName{GivenName: strings.TrimSpace(givenName), Surname: strings.TrimSpace(surname)}

	var rs hertzGRVehCancelRS
	if err := h.call("VehicleCancelReservation", rq, &rs); err != nil {
		return err
	}
	if len(rs.Errors) > 0 {
		e := rs.Errors[0]
		if e.Code == hertzGRAlreadyCancelled {
			return nil
		}
		return fmt.Errorf("hertz gr cancel %s returned %s: %s", confID, e.Code, e.message())
	}

	return nil
}

// GenerateVoucher renders our own voucher, since Autohellas has no voucher of its own to fetch.
func (h *HertzGR) GenerateVoucher(d *VoucherData) (string, error) {
	return renderVoucher(d)
}
