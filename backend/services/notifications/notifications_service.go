package notifications

import (
	"context"
	"fmt"
	"net/http"

	"encore.app/internal/pdf"
	"encore.app/services/notifications/email"
	"encore.app/services/notifications/sms"
	"encore.dev"
	"encore.dev/config"
	"encore.dev/rlog"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// encore:service
type Service struct {
	emailSender  email.Sender
	smsSender    sms.Sender
	pdfConverter pdf.PDFConverter
}

// Config holds the configuration for the notifications service, including email settings.
type Config struct {
	EmailFrom             config.String
	EmailHost             config.String
	EmailPort             config.Int
	SMSUsername           config.String
	SMSSenderName         config.String
	GotenbergURL          config.String
	PasswordResetTokenURL config.String
}

var cfg = config.Load[*Config]()

var secrets struct {
	GoogleServiceAccountJSON string
	smsToken                 string
}

func initService() (*Service, error) {
	ctx := context.Background()

	es, err := email.NewGmailAPISender(
		ctx,
		secrets.GoogleServiceAccountJSON,
		cfg.EmailFrom(),
	)
	rlog.Info("email sender created", "from", cfg.EmailFrom())
	if err != nil {
		rlog.Error("failed to create email sender", "error", err)
		return nil, err
	}

	ss := sms.NewSender(
		secrets.smsToken,
		cfg.SMSSenderName(),
		cfg.SMSUsername(),
	)

	gotenbergClient, err := gotenbergHTTPClient(ctx, cfg.GotenbergURL())
	if err != nil {
		rlog.Error("failed to create gotenberg http client", "error", err)
		return nil, err
	}

	return &Service{
		emailSender:  es,
		smsSender:    ss,
		pdfConverter: pdf.NewPdfConverterWithHTTPClient(cfg.GotenbergURL(), gotenbergClient),
	}, nil
}

func gotenbergHTTPClient(ctx context.Context, audience string) (*http.Client, error) {
	meta := encore.Meta()
	if meta.Environment.Type == encore.EnvTest || meta.Environment.Cloud == encore.CloudLocal {
		return nil, nil
	}

	if meta.Environment.Cloud != encore.CloudGCP {
		return idtoken.NewClient(
			ctx,
			audience,
			option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(secrets.GoogleServiceAccountJSON)),
		)
	}

	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("creating gotenberg identity token client: %w", err)
	}

	return client, nil
}
