package notifications

import (
	"encore.app/services/notifications/email"
	"encore.app/services/notifications/sms"
	"encore.dev/config"
)

// encore:service
type Service struct {
	emailSender *email.Sender
	smsSender   *sms.Sender
}

// Config holds the configuration for the notifications service, including email settings.
type Config struct {
	EmailFrom     config.String
	EmailHost     config.String
	EmailPort     config.Int
	SMSUsername   config.String
	SMSSenderName config.String
}

var cfg = config.Load[*Config]()

var secrets struct {
	EmailPassword string
	SmsToken      string
}

func initService() (*Service, error) {
	es := email.NewSender(
		cfg.EmailFrom(),
		secrets.EmailPassword,
		cfg.EmailHost(),
		cfg.EmailPort(),
	)

	ss := sms.NewSender(
		secrets.SmsToken,
		cfg.SMSSenderName(),
		cfg.SMSUsername(),
	)

	return &Service{
		emailSender: &es,
		smsSender:   &ss,
	}, nil
}
