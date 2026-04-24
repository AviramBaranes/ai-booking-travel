package notifications

import (
	"encore.app/services/notifications/email"
	"encore.app/services/notifications/sms"
	"encore.dev/config"
	"encore.dev/rlog"
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
	emailPassword string
	smsToken      string
}

func initService() (*Service, error) {
	es, err := email.NewSender(
		cfg.EmailFrom(),
		secrets.emailPassword,
		cfg.EmailHost(),
		cfg.EmailPort(),
	)
	if err != nil {
		rlog.Error("failed to create email sender", "error", err)
		return nil, err
	}

	ss := sms.NewSender(
		secrets.smsToken,
		cfg.SMSSenderName(),
		cfg.SMSUsername(),
	)

	return &Service{
		emailSender: &es,
		smsSender:   &ss,
	}, nil
}
