package notifications

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	netmail "net/mail"
	"strings"
	"testing"

	"encore.app/services/accounts"
	accountsuser "encore.app/services/accounts/handlers/user"
	emailevents "encore.app/services/notifications/events"
	"encore.dev/et"
	"github.com/wneessen/go-mail"
)

// makeEmailEvent marshals payload into an EmailEvent for the given event type.
func makeEmailEvent(t *testing.T, eventType emailevents.EmailEventType, payload any) *emailevents.EmailEvent {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	return &emailevents.EmailEvent{Type: eventType, Payload: raw}
}

// decodeSubject MIME-decodes the Subject header from a mail message.
func decodeSubject(t *testing.T, msg *mail.Msg) string {
	t.Helper()
	dec := new(mime.WordDecoder)
	s, err := dec.DecodeHeader(strings.Join(msg.GetGenHeader(mail.HeaderSubject), " "))
	if err != nil {
		t.Fatalf("decode subject: %v", err)
	}
	return s
}

// renderBody renders the full message to a string for body assertions.
func renderBody(t *testing.T, msg *mail.Msg) string {
	t.Helper()
	var buf bytes.Buffer
	if _, err := msg.WriteTo(&buf); err != nil {
		t.Fatalf("writing msg: %v", err)
	}

	parsed, err := netmail.ReadMessage(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("parsing msg: %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(parsed.Header.Get("Content-Type"))
	if err == nil && strings.HasPrefix(mediaType, "multipart/") {
		reader := multipart.NewReader(parsed.Body, params["boundary"])
		for {
			part, err := reader.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				t.Fatalf("reading multipart body: %v", err)
			}

			partMediaType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
			if err == nil && partMediaType != "text/html" && partMediaType != "text/plain" {
				continue
			}

			return decodeBody(t, part.Header.Get("Content-Transfer-Encoding"), part)
		}
		t.Fatal("message body part not found")
	}

	return decodeBody(t, parsed.Header.Get("Content-Transfer-Encoding"), parsed.Body)
}

func decodeBody(t *testing.T, transferEncoding string, body io.Reader) string {
	t.Helper()

	decoded := body
	switch strings.ToLower(strings.TrimSpace(transferEncoding)) {
	case "base64":
		decoded = base64.NewDecoder(base64.StdEncoding, body)
	case "quoted-printable":
		decoded = quotedprintable.NewReader(body)
	}

	raw, err := io.ReadAll(decoded)
	if err != nil {
		t.Fatalf("reading decoded body: %v", err)
	}

	return string(raw)
}

func TestHandleEmailEvent(t *testing.T) {
	ctx := context.Background()

	t.Run("critical error", func(t *testing.T) {
		payload := emailevents.CriticalErrorEmailPayload{
			Subject: "DB exploded",
			Message: "primary node unreachable",
		}

		t.Run("sends to all admins with correct subject and body", func(t *testing.T) {
			et.MockEndpoint(accounts.ListAdminsEmails, func(_ context.Context) (*accounts.ListAdminsEmailsResponse, error) {
				return &accounts.ListAdminsEmailsResponse{
					Emails: []string{"admin1@test.com", "admin2@test.com"},
				}, nil
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeCriticalError, payload)); err != nil {
				t.Fatalf("HandleEmailEvent: %v", err)
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

			// Subject is plain ASCII — no MIME encoding expected.
			if subjects := fake.msg.GetGenHeader(mail.HeaderSubject); len(subjects) != 1 || subjects[0] != payload.Subject {
				t.Errorf("subject = %v, want %q", subjects, payload.Subject)
			}

			if body := renderBody(t, fake.msg); !strings.Contains(body, payload.Message) {
				t.Errorf("body does not contain %q", payload.Message)
			}
		})

		t.Run("returns error and does not send when listing admins fails", func(t *testing.T) {
			et.MockEndpoint(accounts.ListAdminsEmails, func(_ context.Context) (*accounts.ListAdminsEmailsResponse, error) {
				return nil, errors.New("db down")
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeCriticalError, payload)); err == nil {
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

			err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeCriticalError, payload))
			if !errors.Is(err, sendErr) {
				t.Fatalf("err = %v, want %v", err, sendErr)
			}
		})
	})

	t.Run("cancellation", func(t *testing.T) {
		payload := emailevents.CancellationEmailPayload{
			UserID:             42,
			BookingReferenceID: "GR-001",
			DriverFullName:     "John Doe",
		}
		wantSubject := fmt.Sprintf("AI Booking Travel reservation number %s- %s has been cancelled", payload.BookingReferenceID, payload.DriverFullName)

		t.Run("sends to resolved user email with correct subject and body", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, p accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				if p.UserID != payload.UserID {
					return nil, fmt.Errorf("unexpected UserID %d", p.UserID)
				}
				return &accountsuser.GetUserEmailResponse{Email: "john@test.com"}, nil
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeCancellation, payload)); err != nil {
				t.Fatalf("HandleEmailEvent: %v", err)
			}

			recipients, err := fake.msg.GetRecipients()
			if err != nil {
				t.Fatalf("GetRecipients: %v", err)
			}
			if len(recipients) != 1 || recipients[0] != "<john@test.com>" {
				t.Errorf("recipients = %v, want [<john@test.com>]", recipients)
			}

			if got := decodeSubject(t, fake.msg); got != wantSubject {
				t.Errorf("subject = %q, want %q", got, wantSubject)
			}

			if body := renderBody(t, fake.msg); !strings.Contains(body, payload.BookingReferenceID) {
				t.Errorf("body does not contain booking reference %q", payload.BookingReferenceID)
			}
		})

		t.Run("returns error and does not send when resolving user email fails", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, _ accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				return nil, errors.New("user not found")
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeCancellation, payload)); err == nil {
				t.Fatal("expected error, got nil")
			}
			if fake.msg != nil {
				t.Error("email should not have been sent")
			}
		})

		t.Run("propagates email send failure", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, _ accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				return &accountsuser.GetUserEmailResponse{Email: "john@test.com"}, nil
			})

			sendErr := errors.New("smtp boom")
			fake := &fakeEmailSender{err: sendErr}
			s := &Service{emailSender: fake}

			err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeCancellation, payload))
			if !errors.Is(err, sendErr) {
				t.Fatalf("err = %v, want %v", err, sendErr)
			}
		})
	})

	t.Run("late cancellation alert", func(t *testing.T) {
		officeID := int64(77)
		organizationID := int64(88)
		payload := emailevents.LateCancellationAlertEmailPayload{
			ReservationID:       501,
			BrokerReservationID: "BR-501",
			AgentID:             42,
			OfficeID:            &officeID,
			OrganizationID:      &organizationID,
		}
		wantSubject := fmt.Sprintf("בוצע ביטול פחות מ-48 שעות לפני האיסוף - הזמנה %d", payload.ReservationID)

		t.Run("sends to admins with resolved account names in the body", func(t *testing.T) {
			et.MockEndpoint(accounts.ListAdminsEmails, func(_ context.Context) (*accounts.ListAdminsEmailsResponse, error) {
				return &accounts.ListAdminsEmailsResponse{Emails: []string{"admin@test.com"}}, nil
			})

			et.MockEndpoint(accounts.GetAccountsLookup, func(_ context.Context, p accounts.GetAccountsLookupParams) (*accounts.GetAccountsLookupResponse, error) {
				if len(p.UserIDs) != 1 || p.UserIDs[0] != payload.AgentID {
					return nil, fmt.Errorf("unexpected user ids: %v", p.UserIDs)
				}
				if len(p.OfficeIDs) != 1 || p.OfficeIDs[0] != officeID {
					return nil, fmt.Errorf("unexpected office ids: %v", p.OfficeIDs)
				}
				if len(p.OrganizationIDs) != 1 || p.OrganizationIDs[0] != organizationID {
					return nil, fmt.Errorf("unexpected organization ids: %v", p.OrganizationIDs)
				}
				return &accounts.GetAccountsLookupResponse{
					Users:         []accounts.AccountName{{ID: payload.AgentID, Name: "Dana Levy"}},
					Offices:       []accounts.AccountName{{ID: officeID, Name: "Jerusalem Office"}},
					Organizations: []accounts.AccountName{{ID: organizationID, Name: "Global Rental Org"}},
				}, nil
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeLateCancellationAlert, payload)); err != nil {
				t.Fatalf("HandleEmailEvent: %v", err)
			}

			recipients, err := fake.msg.GetRecipients()
			if err != nil {
				t.Fatalf("GetRecipients: %v", err)
			}
			if len(recipients) != 1 || recipients[0] != "<admin@test.com>" {
				t.Errorf("recipients = %v, want [<admin@test.com>]", recipients)
			}

			if got := decodeSubject(t, fake.msg); got != wantSubject {
				t.Errorf("subject = %q, want %q", got, wantSubject)
			}

			body := renderBody(t, fake.msg)
			for _, want := range []string{
				"501",
				"BR-501",
				"42 - Dana Levy",
				"77 - Jerusalem Office",
				"88 - Global Rental Org",
			} {
				if !strings.Contains(body, want) {
					t.Errorf("body does not contain %q", want)
				}
			}
		})

		t.Run("returns error when account lookup fails", func(t *testing.T) {
			et.MockEndpoint(accounts.ListAdminsEmails, func(_ context.Context) (*accounts.ListAdminsEmailsResponse, error) {
				return &accounts.ListAdminsEmailsResponse{Emails: []string{"admin@test.com"}}, nil
			})

			et.MockEndpoint(accounts.GetAccountsLookup, func(_ context.Context, _ accounts.GetAccountsLookupParams) (*accounts.GetAccountsLookupResponse, error) {
				return nil, errors.New("lookup failed")
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeLateCancellationAlert, payload)); err == nil {
				t.Fatal("expected error, got nil")
			}
			if fake.msg != nil {
				t.Error("email should not have been sent")
			}
		})
	})

	t.Run("new order", func(t *testing.T) {
		payload := emailevents.NewOrderEmailPayload{
			UserID:             43,
			BookingReferenceID: "GR-002",
			DriverFullName:     "Jane Smith",
		}
		wantSubject := fmt.Sprintf("Attached AI Booking Travel reservation number %s- %s", payload.BookingReferenceID, payload.DriverFullName)

		t.Run("sends to resolved user email with correct subject and body", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, p accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				if p.UserID != payload.UserID {
					return nil, fmt.Errorf("unexpected UserID %d", p.UserID)
				}
				return &accountsuser.GetUserEmailResponse{Email: "jane@test.com"}, nil
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeNewOrder, payload)); err != nil {
				t.Fatalf("HandleEmailEvent: %v", err)
			}

			recipients, err := fake.msg.GetRecipients()
			if err != nil {
				t.Fatalf("GetRecipients: %v", err)
			}
			if len(recipients) != 1 || recipients[0] != "<jane@test.com>" {
				t.Errorf("recipients = %v, want [<jane@test.com>]", recipients)
			}

			if got := decodeSubject(t, fake.msg); got != wantSubject {
				t.Errorf("subject = %q, want %q", got, wantSubject)
			}

			if body := renderBody(t, fake.msg); !strings.Contains(body, payload.BookingReferenceID) {
				t.Errorf("body does not contain booking reference %q", payload.BookingReferenceID)
			}
		})

		t.Run("returns error and does not send when resolving user email fails", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, _ accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				return nil, errors.New("user not found")
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeNewOrder, payload)); err == nil {
				t.Fatal("expected error, got nil")
			}
			if fake.msg != nil {
				t.Error("email should not have been sent")
			}
		})

		t.Run("propagates email send failure", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, _ accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				return &accountsuser.GetUserEmailResponse{Email: "jane@test.com"}, nil
			})

			sendErr := errors.New("smtp boom")
			fake := &fakeEmailSender{err: sendErr}
			s := &Service{emailSender: fake}

			err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeNewOrder, payload))
			if !errors.Is(err, sendErr) {
				t.Fatalf("err = %v, want %v", err, sendErr)
			}
		})
	})

	t.Run("open order alert", func(t *testing.T) {
		payload := emailevents.OpenOrderAlertEmailPayload{
			UserID:             44,
			BookingReferenceID: "GR-003",
			DriverFullName:     "Bob Cohen",
		}
		wantSubject := fmt.Sprintf("%s - %s - הזמנה פתוחה מס׳", payload.BookingReferenceID, payload.DriverFullName)

		t.Run("sends to resolved user email with correct subject and body", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, p accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				if p.UserID != payload.UserID {
					return nil, fmt.Errorf("unexpected UserID %d", p.UserID)
				}
				return &accountsuser.GetUserEmailResponse{Email: "bob@test.com"}, nil
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeOpenOrderAlert, payload)); err != nil {
				t.Fatalf("HandleEmailEvent: %v", err)
			}

			recipients, err := fake.msg.GetRecipients()
			if err != nil {
				t.Fatalf("GetRecipients: %v", err)
			}
			if len(recipients) != 1 || recipients[0] != "<bob@test.com>" {
				t.Errorf("recipients = %v, want [<bob@test.com>]", recipients)
			}

			if got := decodeSubject(t, fake.msg); got != wantSubject {
				t.Errorf("subject = %q, want %q", got, wantSubject)
			}

			if body := renderBody(t, fake.msg); !strings.Contains(body, payload.BookingReferenceID) {
				t.Errorf("body does not contain booking reference %q", payload.BookingReferenceID)
			}
		})

		t.Run("returns error and does not send when resolving user email fails", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, _ accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				return nil, errors.New("user not found")
			})

			fake := &fakeEmailSender{}
			s := &Service{emailSender: fake}

			if err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeOpenOrderAlert, payload)); err == nil {
				t.Fatal("expected error, got nil")
			}
			if fake.msg != nil {
				t.Error("email should not have been sent")
			}
		})

		t.Run("propagates email send failure", func(t *testing.T) {
			et.MockEndpoint(accounts.GetUserEmail, func(_ context.Context, _ accounts.GetUserEmailParams) (*accountsuser.GetUserEmailResponse, error) {
				return &accountsuser.GetUserEmailResponse{Email: "bob@test.com"}, nil
			})

			sendErr := errors.New("smtp boom")
			fake := &fakeEmailSender{err: sendErr}
			s := &Service{emailSender: fake}

			err := s.HandleEmailEvent(ctx, makeEmailEvent(t, emailevents.EmailEventTypeOpenOrderAlert, payload))
			if !errors.Is(err, sendErr) {
				t.Fatalf("err = %v, want %v", err, sendErr)
			}
		})
	})
}
