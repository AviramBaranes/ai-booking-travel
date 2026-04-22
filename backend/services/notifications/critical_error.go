package notifications

import (
	"context"

	"encore.app/services/notifications/email"
	"encore.dev/pubsub"
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

	return email.SendEmail(
		*s.emailSender,
		[]string{}, //needs to implement admin emails fetching
		event.Subject,
		email.CriticalErrorTemplate,
		template,
	)
}
