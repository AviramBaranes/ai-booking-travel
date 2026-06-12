package notifications

import (
	"context"
	"html/template"

	"encore.app/internal/api_errors"
	"encore.app/services/notifications/email"
	"encore.dev/rlog"
)

type SendCMSEmailParams struct {
	Token   string `header:"X-Translation-Token" encore:"sensitive"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// encore:api public method=POST path=/send-cms-email
func (s *Service) SendCMSEmail(ctx context.Context, p SendCMSEmailParams) error {
	if p.Token != secrets.cmsAPIKey {
		rlog.Warn("invalid cms token", "provided_token", p.Token)
		return api_errors.ErrNotFound
	}

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
