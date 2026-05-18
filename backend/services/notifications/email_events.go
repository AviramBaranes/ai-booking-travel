package notifications

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.dev/pubsub"
)

// EmailEventType identifies which kind of email to send.
type EmailEventType string

const (
	EmailEventTypeCriticalError  EmailEventType = "critical_error"
	EmailEventTypeCancellation   EmailEventType = "cancellation"
	EmailEventTypeNewOrder       EmailEventType = "new_order"
	EmailEventTypeOpenOrderAlert EmailEventType = "open_order_alert"
)

// EmailEvent is the generic envelope published to EmailRequestedTopic.
// Payload holds a JSON-encoded type-specific struct; use PublishEmailEvent to publish safely.
type EmailEvent struct {
	Type    EmailEventType  `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// EmailRequestedTopic is the single topic for all outbound email notifications.
var EmailRequestedTopic = pubsub.NewTopic[*EmailEvent]("email-requested", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

// PublishEmailEvent marshals payload into an EmailEvent and publishes it.
// The type parameter T is inferred from the payload argument, keeping call sites type-safe.
func PublishEmailEvent[T any](ctx context.Context, eventType EmailEventType, payload T) (string, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshaling email event payload: %w", err)
	}
	return EmailRequestedTopic.Publish(ctx, &EmailEvent{
		Type:    eventType,
		Payload: raw,
	})
}

// ---- Per-type payload structs ----

// CriticalErrorEmailPayload is the payload for EmailEventTypeCriticalError.
// Recipients are resolved dynamically by the subscriber (admin emails).
type CriticalErrorEmailPayload struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// CancellationEmailPayload is the payload for EmailEventTypeCancellation.
type CancellationEmailPayload struct {
	UserID             int64  `json:"userId"`
	BookingReferenceID string `json:"bookingReferenceId"`
	DriverFullName     string `json:"driverFullName"`
}

// NewOrderEmailPayload is the payload for EmailEventTypeNewOrder.
type NewOrderEmailPayload struct {
	UserID             int64  `json:"userId"`
	BookingReferenceID string `json:"bookingReferenceId"`
	DriverFullName     string `json:"driverFullName"`
}

// OpenOrderAlertEmailPayload is the payload for EmailEventTypeOpenOrderAlert.
type OpenOrderAlertEmailPayload struct {
	UserID             int64  `json:"userId"`
	BookingReferenceID string `json:"bookingReferenceId"`
	DriverFullName     string `json:"driverFullName"`
}
