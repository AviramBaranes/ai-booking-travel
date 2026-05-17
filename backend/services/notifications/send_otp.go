package notifications

import (
	"context"
	"fmt"

	auth_handler "encore.app/services/accounts/handlers/auth_handler"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

const (
	otpMessageHe = "קוד האימות שלך הוא: %s"
	otpMessageEn = "Your verification code is: %s"
)

var _ = pubsub.NewSubscription(
	auth_handler.CustomerLoginOTPRequestedTopic,
	"send-customer-login-otp-sms",
	pubsub.SubscriptionConfig[*auth_handler.CustomerLoginOTPRequestedEvent]{
		Handler: pubsub.MethodHandler((*Service).SendCustomerLoginOTPSMS),
	},
)

func (s *Service) SendCustomerLoginOTPSMS(ctx context.Context, event *auth_handler.CustomerLoginOTPRequestedEvent) error {
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
