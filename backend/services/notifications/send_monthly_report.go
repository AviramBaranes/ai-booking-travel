package notifications

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"encore.app/services/notifications/email"
)

type SendMonthlyReportRequest struct {
	ContactName  string
	ContactEmail string
	ExcelBase64  string
}

// encore:api private
func (s *Service) SendMonthlyReport(ctx context.Context, p SendMonthlyReportRequest) error {
	excelBytes, err := base64.StdEncoding.DecodeString(p.ExcelBase64)
	if err != nil {
		return fmt.Errorf("decoding excel report: %w", err)
	}

	month := time.Now().AddDate(0, -1, 0).Format("01-2006")
	filename := fmt.Sprintf("monthly_report_%s.xlsx", month)
	subject := fmt.Sprintf("ריכוז חודשי %s", time.Now().AddDate(0, -1, 0).Format("01/2006"))

	attachments := []email.Attachment{
		{
			Filename: filename,
			Reader:   bytes.NewReader(excelBytes),
		},
	}

	return email.SendEmail(
		ctx,
		*s.emailSender,
		[]string{p.ContactEmail},
		subject,
		email.MonthlyReportTemplate,
		email.MonthlyReportData{ContactName: p.ContactName},
		attachments,
	)
}
