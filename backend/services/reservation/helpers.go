package reservation

import (
	"math"

	"encore.app/services/reservation/db"
)

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
