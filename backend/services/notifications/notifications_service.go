package notifications

import (
	"context"

	"encore.app/internal/pdf"
	"encore.app/services/notifications/email"
	"encore.app/services/notifications/sms"
	"encore.dev/config"
	"encore.dev/rlog"
)

// encore:service
type Service struct {
	reservationsEmailSender email.Sender
	accountsEmailSender     email.Sender
	smsSender               sms.Sender
	pdfConverter            pdf.PDFConverter
}

// Config holds the configuration for the notifications service, including email settings.
type Config struct {
	AccountsEmailFrom         config.String
	AccountsEmailFromName     config.String
	ReservationsEmailFrom     config.String
	ReservationsEmailFromName config.String
	SMSUsername               config.String
	SMSSenderName             config.String
	GotenbergURL              config.String
	PasswordResetTokenURL     config.String
	PriceOfferURL             config.String
}

var cfg = config.Load[*Config]()

var secrets struct {
	GoogleServiceAccountJSON string
	smsToken                 string
}

func initService() (*Service, error) {
	ctx := context.Background()

	reservationsSender, err := email.NewGmailAPISender(
		ctx,
		secrets.GoogleServiceAccountJSON,
		cfg.ReservationsEmailFromName(),
		cfg.ReservationsEmailFrom(),
	)
	rlog.Info("reservations email sender created", "fromName", cfg.ReservationsEmailFromName(), "from", cfg.ReservationsEmailFrom())
	if err != nil {
		rlog.Error("failed to create reservations email sender", "error", err)
		return nil, err
	}

	accountsSender, err := email.NewGmailAPISender(
		ctx,
		secrets.GoogleServiceAccountJSON,
		cfg.AccountsEmailFromName(),
		cfg.AccountsEmailFrom(),
	)
	rlog.Info("accounts email sender created", "fromName", cfg.AccountsEmailFromName(), "from", cfg.AccountsEmailFrom())
	if err != nil {
		rlog.Error("failed to create accounts email sender", "error", err)
		return nil, err
	}

	ss := sms.NewSender(
		secrets.smsToken,
		cfg.SMSSenderName(),
		cfg.SMSUsername(),
	)

	pdfConverter, err := pdf.NewDeployedPdfConverter(ctx, cfg.GotenbergURL())
	if err != nil {
		rlog.Error("failed to create gotenberg pdf converter", "error", err)
		return nil, err
	}

	return &Service{
		reservationsEmailSender: reservationsSender,
		accountsEmailSender:     accountsSender,
		smsSender:               ss,
		pdfConverter:            pdfConverter,
	}, nil
}
