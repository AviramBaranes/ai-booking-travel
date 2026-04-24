package email

import (
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/wneessen/go-mail"
)

const templatesDir = "./templates/"

type Sender struct {
	client *mail.Client
	from   string
}

// NewSender creates a new Sender with the given email credentials and SMTP server information.
func NewSender(from, password, host string, port int) (Sender, error) {
	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(from),
		mail.WithPassword(password),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)

	if err != nil {
		return Sender{}, fmt.Errorf("creating mail client: %w", err)
	}

	return Sender{
		client: client,
		from:   from,
	}, nil
}

type Attachment struct {
	Filename string
	Reader   io.Reader
}

// SendEmail sends an email using the provided Sender, recipient list, subject, template, and data.
func SendEmail[T any](ctx context.Context, s Sender, to []string, subject string, t Template[T], data T, attachments []Attachment) error {
	tmpl, err := template.ParseFiles(templatesDir + t.name + ".html")
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	msg := mail.NewMsg()
	if err := msg.From(s.from); err != nil {
		return fmt.Errorf("setting sender: %w", err)
	}
	if err := msg.To(to...); err != nil {
		return fmt.Errorf("setting recipients: %w", err)
	}

	msg.Subject(subject)
	if err := msg.SetBodyHTMLTemplate(tmpl, data); err != nil {
		return fmt.Errorf("setting email body: %w", err)
	}

	for _, attachment := range attachments {
		if err := msg.AttachReader(attachment.Filename, attachment.Reader); err != nil {
			return fmt.Errorf("adding attachment: %w", err)
		}
	}

	if err := s.client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
