package reservation

import (
	"context"
	"fmt"
	"net/http"

	// "time"

	// dbadapters "encore.app/internal/db_adapters"
	emailevents "encore.app/services/notifications/events"
	actions "encore.app/services/reservation/handlers/actions"
	"encore.dev/beta/errs"
	"encore.dev/cron"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var _ = cron.NewJob("alert-open-reservations", cron.JobConfig{
	Title:    "Alert open reservations",
	Schedule: "0 5 * * *",
	Endpoint: AlertOpenReservations,
})

var emailRequestedTopic = pubsub.TopicRef[pubsub.Publisher[*emailevents.EmailEvent]](emailevents.EmailRequestedTopic)
var emailPublisher = emailevents.NewPublisher(emailRequestedTopic)

func (s *Service) newActionService() *actions.ActionService {
	return actions.NewActionService(s.query, s.pool, s.cancellationTopic, s.currencyCache, s.pdfConverter, actions.Config{
		VAT: cfg.VAT(),
	},
	)
}

// encore:api private
func (s *Service) CreateReservation(ctx context.Context, p actions.CreateReservationParams) (*actions.CreateReservationResponse, error) {
	return s.newActionService().CreateReservation(ctx, p)
}

// encore:api auth method=POST path=/api/reservation/:id/cancel tag:agent_customer
func (s *Service) CancelReservation(ctx context.Context, id int64) error {
	return s.newActionService().CancelReservation(ctx, id)
}

// encore:api auth method=POST path=/reservations/:id/voucher tag:agent
func (s *Service) ApplyVoucher(ctx context.Context, id int64, p actions.ApplyVoucherParams) error {
	return s.newActionService().ApplyVoucher(ctx, id, p)
}

// DownloadVoucher streams back the voucher PDF of one of the caller's own reservations, so a user
// can fetch the voucher again without waiting on the email sent when it was first applied.
//
//encore:api auth method=GET path=/reservations/:id/voucher/download tag:agent_customer raw
func (s *Service) DownloadVoucher(w http.ResponseWriter, req *http.Request) {
	id, err := rawPathParamInt64("id")
	if err != nil {
		errs.HTTPError(w, err)
		return
	}

	s.newActionService().DownloadVoucherRaw(w, req, id)
}

// encore:api private
func (s *Service) ResolveReservations(ctx context.Context, p actions.ResolveReservationsParams) error {
	return s.newActionService().ResolveReservations(ctx, p)
}

// encore:api private
func (s *Service) ResolvePenalties(ctx context.Context, p actions.ResolvePenaltiesParams) error {
	return s.newActionService().ResolvePenalties(ctx, p)
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
		// if err != nil {
		// 	rlog.Error("failed to parse pickup date/time", "error", err, "reservationId", r.ID)
		// 	continue
		// }

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

// encore:api private
func (s *Service) VoucherReservationAfterPayment(ctx context.Context, p actions.VoucherReservationAfterPaymentParams) (*actions.VoucherReservationAfterPaymentResponse, error) {
	return s.newActionService().VoucherReservationAfterPayment(ctx, p)
}

// encore:api private
func (s *Service) SaveInvoiceDocNum(ctx context.Context, p actions.SaveInvoiceDocNumParams) error {
	return s.newActionService().SaveInvoiceDocNum(ctx, p)
}
