package notifications

import (
	"context"

	"encore.app/services/accounts"
	"encore.app/services/notifications/email"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

// CriticalErrorEvent holds the subject and message for a critical error notification.
type CriticalErrorEvent struct {
	Subject string
	Message string
}

// CriticalErrorEventTopic is the pub/sub topic for critical error events.
var CriticalErrorEventTopic = pubsub.NewTopic[*CriticalErrorEvent]("critical-error-event", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

var _ = pubsub.NewSubscription(CriticalErrorEventTopic, "send-critical-error-email", pubsub.SubscriptionConfig[*CriticalErrorEvent]{
	Handler: pubsub.MethodHandler((*Service).SendCriticalErrorEmail),
})

// SendCriticalErrorEmail sends a critical error notification email to administrators.
func (s *Service) SendCriticalErrorEmail(ctx context.Context, event *CriticalErrorEvent) error {
	template := email.CriticalErrorData{
		Message: event.Message,
	}

	adminEmails, err := accounts.ListAdminsEmails(ctx)
	if err != nil {
		rlog.Error("failed to query admin emails", "error", err)
		return err
	}

	err = email.SendEmail(
		*s.emailSender,
		adminEmails.Emails,
		event.Subject,
		email.CriticalErrorTemplate,
		template,
	)

	if err != nil {
		rlog.Error("failed to send critical error email", "error", err)
		return err
	}

	return nil
}
