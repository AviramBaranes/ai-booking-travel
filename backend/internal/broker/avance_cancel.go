package broker

import (
	"errors"
	"net/url"
)

// avanceErrAlreadyCancelled lets a repeated cancellation succeed.
const avanceErrAlreadyCancelled = "ERR/104"

// Cancel cancels a reservation by its Wheelsys number. A car already collected (ERR/105) can no
// longer be cancelled and is returned as an error.
func (a *Avance) Cancel(bookingID, _, _ string) error {
	q := url.Values{}
	q.Set("irn", bookingID)

	body, err := a.get("cancel-res", q)
	if err != nil {
		return err
	}

	_, err = parseAvanceReservation("cancel-res", body)
	var apiErr *avanceAPIError
	if errors.As(err, &apiErr) && apiErr.code == avanceErrAlreadyCancelled {
		return nil
	}

	return err
}
