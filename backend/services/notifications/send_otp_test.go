package notifications

import (
	"context"
	"errors"
	"testing"

	"encore.app/services/accounts"
)

type fakeSMSSender struct {
	phoneNumber string
	message     string
	err         error
}

func (f *fakeSMSSender) SendSMS(phoneNumber, message string) error {
	f.phoneNumber = phoneNumber
	f.message = message
	return f.err
}

func TestSendCustomerLoginOTPSMS(t *testing.T) {
	ctx := context.Background()

	t.Run("sends Hebrew OTP message when lang is he", func(t *testing.T) {
		fake := &fakeSMSSender{}
		s := &Service{smsSender: fake}

		event := &accounts.CustomerLoginOTPRequestedEvent{
			PhoneNumber: "+972500000000",
			OTP:         "123456",
			LangCode:    "he",
		}

		if err := s.SendCustomerLoginOTPSMS(ctx, event); err != nil {
			t.Fatalf("SendCustomerLoginOTPSMS: %v", err)
		}

		if fake.phoneNumber != event.PhoneNumber {
			t.Errorf("phoneNumber = %q, want %q", fake.phoneNumber, event.PhoneNumber)
		}
		want := "קוד האימות שלך הוא: 123456"
		if fake.message != want {
			t.Errorf("message = %q, want %q", fake.message, want)
		}
	})

	t.Run("sends English OTP message for non-he lang", func(t *testing.T) {
		fake := &fakeSMSSender{}
		s := &Service{smsSender: fake}

		event := &accounts.CustomerLoginOTPRequestedEvent{
			PhoneNumber: "+15555555555",
			OTP:         "987654",
			LangCode:    "en",
		}

		if err := s.SendCustomerLoginOTPSMS(ctx, event); err != nil {
			t.Fatalf("SendCustomerLoginOTPSMS: %v", err)
		}

		want := "Your verification code is: 987654"
		if fake.message != want {
			t.Errorf("message = %q, want %q", fake.message, want)
		}
	})

	t.Run("propagates SMS send failure", func(t *testing.T) {
		sendErr := errors.New("sms gateway down")
		fake := &fakeSMSSender{err: sendErr}
		s := &Service{smsSender: fake}

		event := &accounts.CustomerLoginOTPRequestedEvent{
			PhoneNumber: "+15555555555",
			OTP:         "111111",
			LangCode:    "en",
		}

		err := s.SendCustomerLoginOTPSMS(ctx, event)
		if !errors.Is(err, sendErr) {
			t.Fatalf("err = %v, want %v", err, sendErr)
		}
	})
}
