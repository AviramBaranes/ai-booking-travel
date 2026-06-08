package notifications

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.app/services/accounts"
	"encore.app/services/notifications/email"
	emailevents "encore.app/services/notifications/events"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var _ = pubsub.NewSubscription(
	emailevents.EmailRequestedTopic,
	"send-email",
	pubsub.SubscriptionConfig[*emailevents.EmailEvent]{
		Handler: pubsub.MethodHandler((*Service).HandleEmailEvent),
	},
)

// HandleEmailEvent dispatches the incoming EmailEvent to the appropriate sender.
func (s *Service) HandleEmailEvent(ctx context.Context, event *emailevents.EmailEvent) error {
	switch event.Type {
	case emailevents.EmailEventTypeCriticalError:
		return s.sendCriticalErrorEmail(ctx, event.Payload)
	case emailevents.EmailEventTypeCancellation:
		return s.sendCancellationEmail(ctx, event.Payload)
	case emailevents.EmailEventTypeLateCancellationAlert:
		return s.sendLateCancellationAlertEmail(ctx, event.Payload)
	case emailevents.EmailEventTypeNewOrder:
		return s.sendNewOrderEmail(ctx, event.Payload)
	case emailevents.EmailEventTypeOpenOrderAlert:
		return s.sendOpenOrderAlertEmail(ctx, event.Payload)
	case emailevents.EmailEventTypePasswordReset:
		return s.sendPasswordResetEmail(ctx, event.Payload)
	case emailevents.EmailEventPriceOfferApproved:
		return s.sendPriceOfferApprovedEmail(ctx, event.Payload)
	default:
		rlog.Error("unknown email event type", "type", event.Type)
		return fmt.Errorf("unknown email event type: %s", event.Type)
	}
}

func (s *Service) sendCriticalErrorEmail(ctx context.Context, raw json.RawMessage) error {
	var p emailevents.CriticalErrorEmailPayload
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
		s.reservationsEmailSender,
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
	var p emailevents.CancellationEmailPayload
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
		s.reservationsEmailSender,
		[]string{userEmail.Email},
		fmt.Sprintf("AI Booking Travel reservation number %s- %s has been cancelled", p.BookingReferenceID, p.DriverFullName),
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

func (s *Service) sendLateCancellationAlertEmail(ctx context.Context, raw json.RawMessage) error {
	var p emailevents.LateCancellationAlertEmailPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshaling late cancellation alert payload: %w", err)
	}

	adminEmails, err := accounts.ListAdminsEmails(ctx)
	if err != nil {
		rlog.Error("failed to query admin emails for late cancellation alert", "error", err)
		return err
	}

	accountsLookup, err := accounts.GetAccountsLookup(ctx, accounts.GetAccountsLookupParams{
		OrganizationIDs: optionalIDSlice(p.OrganizationID),
		OfficeIDs:       optionalIDSlice(p.OfficeID),
		UserIDs:         []int64{p.AgentID},
	})
	if err != nil {
		rlog.Error("failed to resolve account names for late cancellation alert", "error", err, "reservation_id", p.ReservationID)
		return err
	}

	userNames := accountNamesByID(accountsLookup.Users)
	officeNames := accountNamesByID(accountsLookup.Offices)
	organizationNames := accountNamesByID(accountsLookup.Organizations)

	if err := email.SendEmail(
		ctx,
		s.reservationsEmailSender,
		adminEmails.Emails,
		fmt.Sprintf("בוצע ביטול פחות מ-48 שעות לפני האיסוף - הזמנה %d", p.ReservationID),
		email.LateCancellationAlertEmailTemplate,
		email.LateCancellationAlertEmailData{
			ReservationID:       p.ReservationID,
			BrokerReservationID: p.BrokerReservationID,
			AgentLabel:          accountLabel(p.AgentID, userNames),
			OfficeLabel:         optionalAccountLabel(p.OfficeID, officeNames),
			OrganizationLabel:   optionalAccountLabel(p.OrganizationID, organizationNames),
		},
		nil,
	); err != nil {
		rlog.Error("failed to send late cancellation alert email", "error", err, "reservation_id", p.ReservationID)
		return err
	}

	return nil
}

func (s *Service) sendNewOrderEmail(ctx context.Context, raw json.RawMessage) error {
	var p emailevents.NewOrderEmailPayload
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
		s.reservationsEmailSender,
		[]string{userEmail.Email},
		fmt.Sprintf("Attached AI Booking Travel reservation number %s- %s", p.BookingReferenceID, p.DriverFullName),
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

func optionalIDSlice(id *int64) []int64 {
	if id == nil {
		return nil
	}
	return []int64{*id}
}

func accountNamesByID(rows []accounts.AccountName) map[int64]string {
	names := make(map[int64]string, len(rows))
	for _, row := range rows {
		names[row.ID] = row.Name
	}
	return names
}

func accountLabel(id int64, names map[int64]string) string {
	if name, ok := names[id]; ok && name != "" {
		return fmt.Sprintf("%d - %s", id, name)
	}
	return fmt.Sprintf("%d", id)
}

func optionalAccountLabel(id *int64, names map[int64]string) *string {
	if id == nil {
		return nil
	}
	label := accountLabel(*id, names)
	return &label
}

func (s *Service) sendOpenOrderAlertEmail(ctx context.Context, raw json.RawMessage) error {
	var p emailevents.OpenOrderAlertEmailPayload
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
		s.reservationsEmailSender,
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

func (s *Service) sendPasswordResetEmail(ctx context.Context, raw json.RawMessage) error {
	var p emailevents.PasswordResetEmailPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshaling password reset email payload: %w", err)
	}

	resetURL := fmt.Sprintf("%s?token=%s", cfg.PasswordResetTokenURL(), p.TokenHash)
	if err := email.SendEmail(
		ctx,
		s.reservationsEmailSender,
		[]string{p.Email},
		"בקשה לאיפוס סיסמה - AI Booking Travel",
		email.PasswordResetEmailTemplate,
		email.PasswordResetEmailData{ResetURL: resetURL},
		nil,
	); err != nil {
		rlog.Error("failed to send password reset email", "error", err, "email", p.Email)
		return err
	}
	return nil
}

func (s *Service) sendPriceOfferApprovedEmail(ctx context.Context, raw json.RawMessage) error {
	var p emailevents.PriceOfferApprovedEmailPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("unmarshaling price offer approved email payload: %w", err)
	}

	userEmail, err := accounts.GetUserEmail(ctx, accounts.GetUserEmailParams{UserID: p.AgentID})
	if err != nil {
		rlog.Error("failed to resolve recipient for price offer approved email", "error", err, "agent_id", p.AgentID)
		return err
	}

	offerURL := fmt.Sprintf("%s/%d", cfg.PriceOfferURL(), p.PriceOfferID)

	if err := email.SendEmail(
		ctx,
		s.reservationsEmailSender,
		[]string{userEmail.Email},
		"מחיר הצעה אושר - AI Booking Travel",
		email.PriceOfferApprovedEmailTemplate,
		email.PriceOfferApprovedEmailData{
			PriceOfferID:   p.PriceOfferID,
			PriceOfferName: p.PriceOfferName,
			Price:          p.Price,
			Currency:       p.Currency,
			URL:            offerURL,
		},
		nil,
	); err != nil {
		rlog.Error("failed to send price offer approved email", "error", err, "email", userEmail.Email)
		return err
	}
	return nil
}
