package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Sender interface {
	SendSMS(phoneNumber, message string) error
}

// Sender019 is the production Sender implementation backed by the 019sms.co.il HTTP API.
type Sender019 struct {
	token      string
	senderName string
	username   string
	client     *http.Client
}

const (
	smsAPIBaseURL = "https://www.019sms.co.il/api"
	timeout       = 10 * time.Second
)

// NewSender creates a new Sender019 with the provided token, sender name, and username.
func NewSender(token, senderName, username string) *Sender019 {
	return &Sender019{
		token:      token,
		senderName: senderName,
		username:   username,
		client:     &http.Client{Timeout: timeout},
	}
}

// SendSMS sends an SMS message to the specified phone number with the given message content.
func (s *Sender019) SendSMS(phoneNumber string, message string) error {
	reqBody := smsRequest{
		SMS: smsRequestSMS{
			User: smsRequestSMSUser{
				Username: s.username,
			},
			Source: s.senderName,
			Destinations: []smsRequestSMSDestination{
				{Phone: phoneNumber},
			},
			Message: message,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}

	req, err := http.NewRequest("POST", smsAPIBaseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
