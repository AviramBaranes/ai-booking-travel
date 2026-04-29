package notifications

import (
	"context"
	"errors"
	"strings"
	"testing"

	"encore.app/services/accounts"
	"encore.dev/et"
	"github.com/wneessen/go-mail"
)

func TestSendCriticalErrorEmail(t *testing.T) {
	ctx := context.Background()
	event := &CriticalErrorEvent{
		Subject: "DB exploded",
		Message: "primary node unreachable",
	}

	t.Run("sends email to all admins with subject and message in body", func(t *testing.T) {
		et.MockEndpoint(accounts.ListAdminsEmails, func(_ context.Context) (*accounts.ListAdminsEmailsResponse, error) {
			return &accounts.ListAdminsEmailsResponse{
				Emails: []string{"admin1@test.com", "admin2@test.com"},
			}, nil
		})

		fake := &fakeEmailSender{}
		s := &Service{emailSender: fake}

		if err := s.SendCriticalErrorEmail(ctx, event); err != nil {
			t.Fatalf("SendCriticalErrorEmail: %v", err)
		}

		if fake.msg == nil {
			t.Fatal("expected msg to be captured")
		}

		recipients, err := fake.msg.GetRecipients()
		if err != nil {
			t.Fatalf("GetRecipients: %v", err)
		}
		want := []string{"<admin1@test.com>", "<admin2@test.com>"}
		if len(recipients) != len(want) || recipients[0] != want[0] || recipients[1] != want[1] {
			t.Errorf("recipients = %v, want %v", recipients, want)
		}

		// Subject is plain ASCII so no MIME-encoding to decode.
		if subjects := fake.msg.GetGenHeader(mail.HeaderSubject); len(subjects) != 1 || subjects[0] != event.Subject {
			t.Errorf("subject = %v, want %q", subjects, event.Subject)
		}

		var raw strings.Builder
		if _, err := fake.msg.WriteTo(&raw); err != nil {
			t.Fatalf("writing msg: %v", err)
		}
		if !strings.Contains(raw.String(), event.Message) {
			t.Errorf("rendered message does not contain %q", event.Message)
		}
	})

	t.Run("returns error and does not send when listing admins fails", func(t *testing.T) {
		et.MockEndpoint(accounts.ListAdminsEmails, func(_ context.Context) (*accounts.ListAdminsEmailsResponse, error) {
			return nil, errors.New("db down")
		})

		fake := &fakeEmailSender{}
		s := &Service{emailSender: fake}

		err := s.SendCriticalErrorEmail(ctx, event)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if fake.msg != nil {
			t.Error("email should not have been sent")
		}
	})

	t.Run("propagates email send failure", func(t *testing.T) {
		et.MockEndpoint(accounts.ListAdminsEmails, func(_ context.Context) (*accounts.ListAdminsEmailsResponse, error) {
			return &accounts.ListAdminsEmailsResponse{Emails: []string{"admin@test.com"}}, nil
		})

		sendErr := errors.New("smtp boom")
		fake := &fakeEmailSender{err: sendErr}
		s := &Service{emailSender: fake}

		err := s.SendCriticalErrorEmail(ctx, event)
		if !errors.Is(err, sendErr) {
			t.Fatalf("err = %v, want %v", err, sendErr)
		}
	})
}
