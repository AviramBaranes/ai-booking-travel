package reservation

import (
	"context"
	"fmt"

	// "time"

	// dbadapters "encore.app/internal/db_adapters"
	emailevents "encore.app/services/notifications/events"
	actions "encore.app/services/reservation/handlers/actions"
	"encore.dev/cron"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var _ = cron.NewJob("alert-open-reservations", cron.JobConfig{
	Title:    "Alert open reservations",
	Schedule: "0 7 * * *",
	Endpoint: AlertOpenReservations,
})

var emailRequestedTopic = pubsub.TopicRef[pubsub.Publisher[*emailevents.EmailEvent]](emailevents.EmailRequestedTopic)
var emailPublisher = emailevents.NewPublisher(emailRequestedTopic)

func (s *Service) newActionService() *actions.ActionService {
	return actions.NewActionService(s.query, s.pool, s.cancellationTopic, actions.Config{
		VAT:        cfg.VAT(),
		IcountCID:  cfg.Icount.CID(),
		IcountUser: cfg.Icount.User(),
	})
}

// encore:api private
func (s *Service) CreateReservation(ctx context.Context, p actions.CreateReservationParams) (*actions.CreateReservationResponse, error) {
	return s.newActionService().CreateReservation(ctx, p)
}

// encore:api auth method=POST path=/api/reservation/:id/cancel tag:agent
func (s *Service) CancelReservation(ctx context.Context, id int64) error {
	return s.newActionService().CancelReservation(ctx, id)
}

// encore:api auth method=POST path=/reservations/:id/voucher tag:agent
func (s *Service) ApplyVoucher(ctx context.Context, id int64, p actions.ApplyVoucherParams) error {
	return s.newActionService().ApplyVoucher(ctx, id, p)
}

// encore:api private
func (s *Service) ResolveReservations(ctx context.Context, p actions.ResolveReservationsParams) error {
	return s.newActionService().ResolveReservations(ctx, p)
}

// encore:api private
func (s *Service) AlertOpenReservations(ctx context.Context) error {
	reservations, err := s.query.GetOpenReservationsPickingUpWithinWeek(ctx)
	if err != nil {
		rlog.Error("failed to query open reservations picking up within a week", "error", err)
		return err
	}

	for _, r := range reservations {
		// pickupDateTime, err := dbadapters.CombineDateTime(r.PickupDate, r.PickupTime)
		if err != nil {
			rlog.Error("failed to parse pickup date/time", "error", err, "reservationId", r.ID)
			continue
		}

		// if time.Until(pickupDateTime) > cancellationWindowHours*time.Hour {
		// More than 48h until pickup — send open order alert
		if _, err := emailPublisher.Publish(ctx, emailevents.EmailEventTypeOpenOrderAlert, emailevents.OpenOrderAlertEmailPayload{
			UserID:             r.UserID,
			BookingReferenceID: r.BrokerReservationID,
			DriverFullName:     fmt.Sprintf("%s %s %s", r.DriverTitle, r.DriverFirstName, r.DriverLastName),
		}); err != nil {
			rlog.Error("failed to publish open order alert email", "error", err, "reservationId", r.ID)
		}
		// } else {
		// 	// Within 48h of pickup and still not vouchered — auto-cancel
		// 	if err := CancelReservation(ctx, r.ID); err != nil {
		// 		rlog.Error("failed to auto-cancel open reservation", "error", err, "reservationId", r.ID)
		// 	}
		// }
	}

	return nil
}

// encore:api private
func (s *Service) SeedReservations(ctx context.Context) (*actions.SeedReservationsResponse, error) {
	return s.newActionService().SeedReservations(ctx)
}
