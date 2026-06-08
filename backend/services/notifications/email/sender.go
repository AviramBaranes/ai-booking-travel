package email

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"

	"github.com/wneessen/go-mail"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Sender delivers a prepared mail message. Implemented by *SMTPSender in
// production and by fakes in tests.
type Sender interface {
	Send(ctx context.Context, msg *mail.Msg) error
}

//go:embed templates/*.html
var templatesFS embed.FS

type Attachment struct {
	Filename string
	Reader   io.Reader
}

type GmailAPISender struct {
	srv      *gmail.Service
	fromName string
	from     string
}

func NewGmailAPISender(ctx context.Context, serviceAccountJSON, fromName, from string) (*GmailAPISender, error) {
	config, err := google.JWTConfigFromJSON(
		[]byte(serviceAccountJSON),
		gmail.GmailSendScope,
	)
	if err != nil {
		return nil, fmt.Errorf("parsing service account json: %w", err)
	}

	config.Subject = from

	client := config.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("creating gmail service: %w", err)
	}

	return &GmailAPISender{
		srv:      srv,
		fromName: fromName,
		from:     from,
	}, nil
}

func (s *GmailAPISender) Send(ctx context.Context, msg *mail.Msg) error {
	if err := msg.FromFormat(s.fromName, s.from); err != nil {
		return fmt.Errorf("setting sender: %w", err)
	}

	var buf bytes.Buffer
	if _, err := msg.WriteTo(&buf); err != nil {
		return fmt.Errorf("building raw email: %w", err)
	}

	raw := base64.RawURLEncoding.EncodeToString(buf.Bytes())

	_, err := s.srv.Users.Messages.Send("me", &gmail.Message{
		Raw: raw,
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("sending gmail message: %w", err)
	}

	return nil
}

func SendEmail[T any](ctx context.Context, s Sender, to []string, subject string, t Template[T], data T, attachments []Attachment) error {
	tmpl, err := template.New(t.name+".html").ParseFS(templatesFS, "templates/"+t.name+".html", "templates/_layout.html")
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	msg := mail.NewMsg()
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

	if err := s.Send(ctx, msg); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
