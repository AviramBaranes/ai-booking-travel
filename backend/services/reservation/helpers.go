package reservation

import (
	"math"

	"encore.app/services/reservation/db"
)

// cancellationWindowHours is the number of hours before pickup within which a
// reservation cannot be cancelled. Must match the value in handlers/actions.
const cancellationWindowHours = 48

func nullStatusFromString(s string) db.NullReservationStatus {
	if s == "" {
		return db.NullReservationStatus{}
	}
	return db.NullReservationStatus{ReservationStatus: db.ReservationStatus(s), Valid: true}
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}
