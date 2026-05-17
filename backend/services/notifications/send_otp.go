package notifications

import (
	"context"
	"fmt"

	"encore.app/services/accounts/handlers/auth"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

const (
	otpMessageHe = "קוד האימות שלך הוא: %s"
	otpMessageEn = "Your verification code is: %s"
)

var _ = pubsub.NewSubscription(
	auth.CustomerLoginOTPRequestedTopic,
	"send-customer-login-otp-sms",
	pubsub.SubscriptionConfig[*auth.CustomerLoginOTPRequestedEvent]{
		Handler: pubsub.MethodHandler((*Service).SendCustomerLoginOTPSMS),
	},
)

func (s *Service) SendCustomerLoginOTPSMS(ctx context.Context, event *auth.CustomerLoginOTPRequestedEvent) error {
	template := otpMessageEn
	if event.LangCode == "he" {
		template = otpMessageHe
	}

	message := fmt.Sprintf(template, event.OTP)

	if err := s.smsSender.SendSMS(event.PhoneNumber, message); err != nil {
		rlog.Error("failed to send OTP SMS", "phone_number", event.PhoneNumber, "error", err)
		return err
	}

	return nil
}
