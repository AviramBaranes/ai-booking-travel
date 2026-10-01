package broker

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	avanceResConfirmed = "RES"
	avanceStatusOK     = "OK"

	// avanceAdditionalDriver is the extra whose first unit Avance gives free. Our add-on counts only
	// the paid drivers, so booking sends one more than was selected.
	avanceAdditionalDriver = "ADD"
)

// ErrBookingLeftOpen means the broker created a booking we could not confirm and then failed to
// cancel, so it is still open at the broker and must be cancelled by hand.
var ErrBookingLeftOpen = errors.New("booking left open at the broker")

// avanceAPIError is an error Wheelsys reported inside a successful HTTP response.
type avanceAPIError struct {
	page    string
	code    string
	message string
}

func (e *avanceAPIError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("avance %s returned %s", e.page, e.code)
	}
	return fmt.Sprintf("avance %s returned %s: %s", e.page, e.code, e.message)
}

// Book reserves the car at the price the search quoted, by passing the quote back as QUOTEREF.
// Wheelsys only accepts it with the dates, stations and product the quote was made with.
func (a *Avance) Book(p BookingParams) (BookingResponse, error) {
	q, err := avanceBookingQuery(p)
	if err != nil {
		return BookingResponse{}, err
	}

	body, err := a.get("new-res", q)
	if err != nil {
		return BookingResponse{}, err
	}

	res, err := parseAvanceReservation("new-res", body)
	if err != nil {
		return BookingResponse{}, err
	}

	if res.Reservation.ResStatus == avanceResConfirmed && res.Reservation.IRN != "" {
		return BookingResponse{ConfirmationNumber: res.Reservation.IRN}, nil
	}

	// Anything but a confirmed reservation, typically a group that went on request between quote
	// and booking, is released: the customer has paid for a confirmed car.
	irn := res.Reservation.IRN
	if irn != "" {
		if err := a.Cancel(irn, "", ""); err != nil {
			return BookingResponse{}, fmt.Errorf("%w: avance new-res %s came back %q and could not be cancelled: %v", ErrBookingLeftOpen, irn, res.Reservation.ResStatus, err)
		}
	}

	return BookingResponse{}, fmt.Errorf("avance new-res %s came back %q and was cancelled", irn, res.Reservation.ResStatus)
}

func avanceBookingQuery(p BookingParams) (url.Values, error) {
	quoteID, group, ok := strings.Cut(p.RateQualifier, ":")
	if !ok || quoteID == "" || group == "" {
		return nil, fmt.Errorf("avance book: malformed rate qualifier %q", p.RateQualifier)
	}

	cdp, err := avanceCDPForPlan(p.PlanID)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("DATE_FROM", formatDate(p.PickupDate))
	q.Set("TIME_FROM", p.PickupTime)
	q.Set("DATE_TO", formatDate(p.DropoffDate))
	q.Set("TIME_TO", p.DropoffTime)
	q.Set("PICKUP_STATION", p.PickupLocation)
	q.Set("RETURN_STATION", p.DropoffLocation)
	q.Set("CDP", cdp)
	q.Set("GROUP", group)
	q.Set("QUOTEREF", quoteID)
	q.Set("VOUCHERNO", avanceVoucherNo(p.RateQualifier))
	q.Set("CUSTFIRST_NAME", truncateRunes(p.DriverFirstName, 30))
	q.Set("CUSTLAST_NAME", truncateRunes(p.DriverLastName, 25))
	if p.FlightNumber != "" {
		q.Set("PICKUP_INFO", truncateRunes(p.FlightNumber, 50))
	}
	if age, err := strconv.Atoi(p.DriverAge); err == nil && age > 0 {
		q.Set("DRIVERAGE", p.DriverAge)
	}

	for _, ao := range p.SelectedAddOns {
		if ao.Quantity <= 0 {
			continue
		}

		code, ok := avanceAddOnCode(ao.ID)
		if !ok {
			return nil, fmt.Errorf("avance book: unknown add-on id %d", ao.ID)
		}

		qty := ao.Quantity
		if code == avanceAdditionalDriver {
			qty++
		}
		q.Set(code, strconv.Itoa(qty))
	}

	return q, nil
}

// avanceVoucherNo derives our reference for the booking from the quote and group it books. Wheelsys
// rejects a reference it has seen before, so a repeated call for the same plan cannot book twice.
func avanceVoucherNo(rateQualifier string) string {
	sum := sha256.Sum256([]byte(rateQualifier))
	return "AIBT" + strings.ToUpper(hex.EncodeToString(sum[:]))[:20]
}

func avanceCDPForPlan(planID string) (string, error) {
	id, err := strconv.Atoi(planID)
	if err != nil {
		return "", fmt.Errorf("avance book: invalid plan id %q: %w", planID, err)
	}

	for _, prod := range avanceProducts {
		if prod.planID == id {
			return prod.cdp, nil
		}
	}

	return "", fmt.Errorf("avance book: unknown plan id %d", id)
}

func avanceAddOnCode(id int) (string, bool) {
	for code, addOnID := range avanceAddOnIDs {
		if addOnID == id {
			return code, true
		}
	}
	return "", false
}

// parseAvanceReservation reads a new-res or cancel-res response. The documentation only says a
// failure is an ERR code in status, so an <errors> block, as the quote returns, is handled too.
func parseAvanceReservation(page string, body []byte) (*avanceReservationResXML, error) {
	var res avanceReservationResXML
	if err := xml.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("avance %s unmarshal response: %w", page, err)
	}

	if len(res.Errors) > 0 {
		e := res.Errors[0]
		return nil, &avanceAPIError{page: page, code: e.Code, message: strings.TrimSpace(e.Message)}
	}

	if res.Reservation.Status != avanceStatusOK {
		return nil, &avanceAPIError{page: page, code: res.Reservation.Status}
	}

	return &res, nil
}

func truncateRunes(s string, n int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= n {
		return string(r)
	}
	return string(r[:n])
}
