package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

const (
	htmlHeaders  = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	templatesDir = "./templates/"
)

type Sender struct {
	auth smtp.Auth
	from string
	adr  string
}

// NewSender creates a new Sender with the given email credentials and SMTP server information.
func NewSender(from, password, host string, port int) Sender {
	auth := smtp.PlainAuth(
		"",
		from,
		password,
		host,
	)

	return Sender{
		auth: auth,
		from: from,
		adr:  fmt.Sprintf("%s:%d", host, port),
	}
}

// SendEmail sends an email using the provided Sender, recipient list, subject, template, and data.
func SendEmail[T any](s Sender, to []string, subject string, t Template[T], data T) error {
	tmpl, err := template.ParseFiles(templatesDir + t.name + ".html")
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var renderedBody bytes.Buffer
	if err := tmpl.Execute(&renderedBody, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	msg := "Subject: " + subject + "\n" + htmlHeaders + "\n" + renderedBody.String()

	if err := smtp.SendMail(
		s.adr,
		s.auth,
		s.from,
		to,
		[]byte(msg),
	); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
