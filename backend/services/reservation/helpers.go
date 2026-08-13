package reservation

import (
	"math"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/services/reservation/db"
	"encore.dev"
)

// rawPathParamInt64 reads a numeric path parameter of the request being served. Raw endpoints get
// the bare *http.Request, which carries no parsed path parameters, so they have to come off the
// Encore request metadata instead.
func rawPathParamInt64(name string) (int64, error) {
	req := encore.CurrentRequest()
	if req == nil {
		return 0, api_errors.ErrInternalError
	}

	id, err := strconv.ParseInt(req.PathParams.Get(name), 10, 64)
	if err != nil {
		return 0, api_errors.ErrInvalidValue
	}

	return id, nil
}

func nullStatusFromString(s string) db.NullReservationStatus {
	if s == "" {
		return db.NullReservationStatus{}
	}
	return db.NullReservationStatus{ReservationStatus: db.ReservationStatus(s), Valid: true}
}

func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}
