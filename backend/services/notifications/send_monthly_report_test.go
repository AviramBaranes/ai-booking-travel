package notifications

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"strings"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.dev/beta/errs"
	"github.com/wneessen/go-mail"
)

type fakeEmailSender struct {
	msg *mail.Msg
	err error
}

func (f *fakeEmailSender) Send(_ context.Context, msg *mail.Msg) error {
	f.msg = msg
	return f.err
}

func TestSendMonthlyReport(t *testing.T) {
	ctx := context.Background()

	t.Run("returns InvalidArgument when excel base64 fails to decode", func(t *testing.T) {
		fake := &fakeEmailSender{}
		s := &Service{emailSender: fake}

		err := s.SendMonthlyReport(ctx, SendMonthlyReportRequest{
			ContactName:  "Alice",
			ContactEmail: "alice@test.com",
			ExcelBase64:  "not valid base64!!!",
		})

		want := api_errors.NewError(errs.InvalidArgument, "invalid excel report encoding")
		api_errors.AssertApiError(t, want, err)

		if fake.msg != nil {
			t.Errorf("email should not have been sent")
		}
	})

	t.Run("sends email with correct recipient, subject, template data, and attachment", func(t *testing.T) {
		excelBytes := []byte("excel-bytes-payload")
		req := SendMonthlyReportRequest{
			ContactName:  "Alice",
			ContactEmail: "alice@test.com",
			ExcelBase64:  base64.StdEncoding.EncodeToString(excelBytes),
		}

		prevMonth := time.Now().AddDate(0, -1, 0)
		wantSubject := fmt.Sprintf("ריכוז חודשי %s", prevMonth.Format("01/2006"))
		wantFilename := fmt.Sprintf("monthly_report_%s.xlsx", prevMonth.Format("01-2006"))

		fake := &fakeEmailSender{}
		s := &Service{emailSender: fake}

		if err := s.SendMonthlyReport(ctx, req); err != nil {
			t.Fatalf("SendMonthlyReport: %v", err)
		}

		if fake.msg == nil {
			t.Fatal("expected msg to be captured")
		}

		recipients, err := fake.msg.GetRecipients()
		if err != nil {
			t.Fatalf("GetRecipients: %v", err)
		}
		want := "<" + req.ContactEmail + ">"
		if len(recipients) != 1 || recipients[0] != want {
			t.Errorf("recipients = %v, want [%q]", recipients, want)
		}

		dec := new(mime.WordDecoder)
		gotSubject, err := dec.DecodeHeader(strings.Join(fake.msg.GetGenHeader(mail.HeaderSubject), " "))
		if err != nil {
			t.Fatalf("decode subject: %v", err)
		}
		if gotSubject != wantSubject {
			t.Errorf("subject = %q, want %q", gotSubject, wantSubject)
		}

		// Render the message and assert it contains rendered template data.
		var raw bytes.Buffer
		if _, err := fake.msg.WriteTo(&raw); err != nil {
			t.Fatalf("writing msg: %v", err)
		}
		if !strings.Contains(raw.String(), req.ContactName) {
			t.Errorf("rendered message does not contain contact name %q", req.ContactName)
		}

		attachments := fake.msg.GetAttachments()
		if len(attachments) != 1 {
			t.Fatalf("attachments len = %d, want 1", len(attachments))
		}
		if attachments[0].Name != wantFilename {
			t.Errorf("attachment filename = %q, want %q", attachments[0].Name, wantFilename)
		}
		var attachBuf bytes.Buffer
		if _, err := attachments[0].Writer(&attachBuf); err != nil {
			t.Fatalf("reading attachment: %v", err)
		}
		if attachBuf.String() != string(excelBytes) {
			t.Errorf("attachment bytes = %q, want %q", attachBuf.String(), excelBytes)
		}
	})

	t.Run("propagates email send failure", func(t *testing.T) {
		sendErr := errors.New("smtp boom")
		fake := &fakeEmailSender{err: sendErr}
		s := &Service{emailSender: fake}

		err := s.SendMonthlyReport(ctx, SendMonthlyReportRequest{
			ContactName:  "Alice",
			ContactEmail: "alice@test.com",
			ExcelBase64:  base64.StdEncoding.EncodeToString([]byte("x")),
		})
		if !errors.Is(err, sendErr) {
			t.Fatalf("err = %v, want %v", err, sendErr)
		}
	})
}
