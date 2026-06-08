package events

import (
	"context"
	"encoding/json"
	"fmt"

	"encore.dev/pubsub"
)

// EmailEventType identifies which kind of email to send.
type EmailEventType string

const (
	EmailEventTypeCriticalError         EmailEventType = "critical_error"
	EmailEventTypeCancellation          EmailEventType = "cancellation"
	EmailEventTypeLateCancellationAlert EmailEventType = "late_cancellation_alert"
	EmailEventTypeNewOrder              EmailEventType = "new_order"
	EmailEventTypeOpenOrderAlert        EmailEventType = "open_order_alert"
	EmailEventTypePasswordReset         EmailEventType = "password_reset"
	EmailEventPriceOfferApproved        EmailEventType = "price_offer_approved"
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

// NewEmailEvent marshals payload into an EmailEvent without publishing it.
func NewEmailEvent[T any](eventType EmailEventType, payload T) (*EmailEvent, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling email event payload: %w", err)
	}

	return &EmailEvent{
		Type:    eventType,
		Payload: raw,
	}, nil
}

type Publisher struct {
	topic pubsub.Publisher[*EmailEvent]
}

func NewPublisher(topic pubsub.Publisher[*EmailEvent]) Publisher {
	return Publisher{topic: topic}
}

func (p Publisher) Publish(ctx context.Context, eventType EmailEventType, payload any) (string, error) {
	event, err := NewEmailEvent(eventType, payload)
	if err != nil {
		return "", err
	}

	return p.topic.Publish(ctx, event)
}

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

// LateCancellationAlertEmailPayload is the payload for EmailEventTypeLateCancellationAlert.
// Recipients are resolved dynamically by the subscriber (admin emails).
type LateCancellationAlertEmailPayload struct {
	ReservationID       int64  `json:"reservationId"`
	BrokerReservationID string `json:"brokerReservationId"`
	AgentID             int64  `json:"agentId"`
	OfficeID            *int64 `json:"officeId,omitempty"`
	OrganizationID      *int64 `json:"organizationId,omitempty"`
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

// PasswordResetEmailPayload is the payload for PasswordResetEmailEventType.
type PasswordResetEmailPayload struct {
	Email     string `json:"email"`
	TokenHash string `json:"tokenHash" encore:"sensitive"`
}

// PriceOfferApprovedEmailPayload is the payload for EmailEventPriceOfferApproved.
type PriceOfferApprovedEmailPayload struct {
	AgentID        int64   `json:"agentId"`
	PriceOfferID   int64   `json:"priceOfferId"`
	PriceOfferName string  `json:"priceOfferName"`
	Price          float64 `json:"price"`
	Currency       string  `json:"currency"`
}
