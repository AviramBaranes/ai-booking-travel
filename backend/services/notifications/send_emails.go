package notifications

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.app/services/accounts"
	"encore.app/services/notifications/email"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var _ = pubsub.NewSubscription(
	EmailRequestedTopic,
	"send-email",
	pubsub.SubscriptionConfig[*EmailEvent]{
		Handler: pubsub.MethodHandler((*Service).HandleEmailEvent),
	},
)

// HandleEmailEvent dispatches the incoming EmailEvent to the appropriate sender.
func (s *Service) HandleEmailEvent(ctx context.Context, event *EmailEvent) error {
	switch event.Type {
	case EmailEventTypeCriticalError:
		return s.sendCriticalErrorEmail(ctx, event.Payload)
	case EmailEventTypeCancellation:
		return s.sendCancellationEmail(ctx, event.Payload)
	case EmailEventTypeNewOrder:
		return s.sendNewOrderEmail(ctx, event.Payload)
	case EmailEventTypeOpenOrderAlert:
		return s.sendOpenOrderAlertEmail(ctx, event.Payload)
	default:
		rlog.Error("unknown email event type", "type", event.Type)
		return fmt.Errorf("unknown email event type: %s", event.Type)
	}
}

func (s *Service) sendCriticalErrorEmail(ctx context.Context, raw json.RawMessage) error {
	var p CriticalErrorEmailPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshaling critical error payload: %w", err)
	}

	adminEmails, err := accounts.ListAdminsEmails(ctx)
	if err != nil {
		rlog.Error("failed to query admin emails", "error", err)
		return err
	}

	if err := email.SendEmail(
		ctx,
		s.emailSender,
		adminEmails.Emails,
		p.Subject,
		email.CriticalErrorTemplate,
		email.CriticalErrorData{Message: p.Message},
		nil,
	); err != nil {
		rlog.Error("failed to send critical error email", "error", err)
		return err
	}
	return nil
}

func (s *Service) sendCancellationEmail(ctx context.Context, raw json.RawMessage) error {
	var p CancellationEmailPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshaling cancellation email payload: %w", err)
	}

	userEmail, err := accounts.GetUserEmail(ctx, accounts.GetUserEmailParams{UserID: p.UserID})
	if err != nil {
		rlog.Error("failed to resolve recipient for cancellation email", "error", err, "user_id", p.UserID)
		return err
	}

	if err := email.SendEmail(
		ctx,
		s.emailSender,
		[]string{userEmail.Email},
		fmt.Sprintf("Global Rentals reservation number %s- %s has been cancelled", p.BookingReferenceID, p.DriverFullName),
		email.CancellationEmailTemplate,
		email.CancellationEmailData{
			BookingReferenceID: p.BookingReferenceID,
			DriverFullName:     p.DriverFullName,
		},
		nil,
	); err != nil {
		rlog.Error("failed to send cancellation email", "error", err, "booking_reference_id", p.BookingReferenceID)
		return err
	}
	return nil
}

func (s *Service) sendNewOrderEmail(ctx context.Context, raw json.RawMessage) error {
	var p NewOrderEmailPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshaling new order email payload: %w", err)
	}

	userEmail, err := accounts.GetUserEmail(ctx, accounts.GetUserEmailParams{UserID: p.UserID})
	if err != nil {
		rlog.Error("failed to resolve recipient for new order email", "error", err, "user_id", p.UserID)
		return err
	}

	if err := email.SendEmail(
		ctx,
		s.emailSender,
		[]string{userEmail.Email},
		fmt.Sprintf("Attached Global Rentals reservation number %s- %s", p.BookingReferenceID, p.DriverFullName),
		email.NewOrderEmailTemplate,
		email.NewOrderEmailData{
			BookingReferenceID: p.BookingReferenceID,
			DriverFullName:     p.DriverFullName,
		},
		nil,
	); err != nil {
		rlog.Error("failed to send new order email", "error", err, "booking_reference_id", p.BookingReferenceID)
		return err
	}
	return nil
}

func (s *Service) sendOpenOrderAlertEmail(ctx context.Context, raw json.RawMessage) error {
	var p OpenOrderAlertEmailPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshaling open order alert payload: %w", err)
	}

	userEmail, err := accounts.GetUserEmail(ctx, accounts.GetUserEmailParams{UserID: p.UserID})
	if err != nil {
		rlog.Error("failed to resolve recipient for open order alert email", "error", err, "user_id", p.UserID)
		return err
	}

	if err := email.SendEmail(
		ctx,
		s.emailSender,
		[]string{userEmail.Email},
		fmt.Sprintf("%s - %s - הזמנה פתוחה מס׳", p.BookingReferenceID, p.DriverFullName),
		email.OpenOrderAlertEmailTemplate,
		email.OpenOrderAlertEmailData{
			BookingReferenceID: p.BookingReferenceID,
			DriverFullName:     p.DriverFullName,
		},
		nil,
	); err != nil {
		rlog.Error("failed to send open order alert email", "error", err, "booking_reference_id", p.BookingReferenceID)
		return err
	}
	return nil
}
