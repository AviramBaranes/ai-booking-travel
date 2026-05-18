package reservation

import (
	"context"
	"fmt"
	"time"

	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/notifications"
	"encore.dev/cron"
	"encore.dev/rlog"
)

var _ = cron.NewJob("alert-open-reservations", cron.JobConfig{
	Title:    "Alert open reservations",
	Schedule: "0 7 * * *",
	Endpoint: AlertOpenReservations,
})

//encore:api private
func (s *Service) AlertOpenReservations(ctx context.Context) error {
	reservations, err := s.query.GetOpenReservationsPickingUpWithinWeek(ctx)
	if err != nil {
		rlog.Error("failed to query open reservations picking up within a week", "error", err)
		return err
	}

	for _, r := range reservations {
		pickupDateTime, err := dbadapters.CombineDateTime(r.PickupDate, r.PickupTime)
		if err != nil {
			rlog.Error("failed to parse pickup date/time", "error", err, "reservationId", r.ID)
			continue
		}

		if time.Until(pickupDateTime) > cancellationWindowHours*time.Hour {
			// More than 48h until pickup — send open order alert
			if _, err := notifications.PublishEmailEvent(ctx, notifications.EmailEventTypeOpenOrderAlert, notifications.OpenOrderAlertEmailPayload{
				UserID:             r.UserID,
				BookingReferenceID: r.BrokerReservationID,
				DriverFullName:     fmt.Sprintf("%s %s %s", r.DriverTitle, r.DriverFirstName, r.DriverLastName),
			}); err != nil {
				rlog.Error("failed to publish open order alert email", "error", err, "reservationId", r.ID)
			}
		} else {
			// Within 48h of pickup and still not vouchered — auto-cancel
			if err := CancelReservation(ctx, r.ID); err != nil {
				rlog.Error("failed to auto-cancel open reservation", "error", err, "reservationId", r.ID)
			}
		}
	}

	return nil
}
