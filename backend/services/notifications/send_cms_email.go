package notifications

import (
	"context"
	"html/template"

	"encore.app/services/notifications/email"
	"encore.dev/rlog"
)

type SendCMSEmailParams struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// encore:api public method=POST path=/send-cms-email tag:service_client
func (s *Service) SendCMSEmail(ctx context.Context, p SendCMSEmailParams) error {
	if err := email.SendEmail(
		ctx,
		s.accountsEmailSender,
		[]string{p.To},
		p.Subject,
		email.CMSEmailTemplate,
		email.CMSEmailData{Content: template.HTML(p.Content)},
		nil,
	); err != nil {
		rlog.Error("failed to send CMS email", "error", err)
		return err
	}

	return nil
}
