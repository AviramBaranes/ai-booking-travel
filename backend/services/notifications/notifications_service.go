package notifications

import (
	"encore.app/services/notifications/email"
	"encore.dev/config"
)

// encore:service
type Service struct {
	emailSender *email.Sender
}

// Config holds the configuration for the notifications service, including email settings.
type Config struct {
	EmailFrom config.String
	EmailHost config.String
	EmailPort config.Int
}

var cfg = config.Load[*Config]()

var secrets struct {
	EmailPassword string
}

func initService() (*Service, error) {
	s := email.NewSender(
		cfg.EmailFrom(),
		secrets.EmailPassword,
		cfg.EmailHost(),
		cfg.EmailPort(),
	)

	return &Service{
		emailSender: &s,
	}, nil
}
