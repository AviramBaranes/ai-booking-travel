package notifications

import (
	"bytes"
	"context"
	"embed"
	"fmt"

	"encore.app/services/notifications/email"
	"encore.dev/rlog"
)

//go:embed assets/*.pdf
var assetsFS embed.FS

// VoucherBroker identifies which broker issued the voucher, used to select the correct static attachments.
type VoucherBroker string

const (
	VoucherBrokerFlex  VoucherBroker = "flex"
	VoucherBrokerHertz VoucherBroker = "hertz"

	TermsFileName = "תנאים-כלליים.pdf"
	ERPFileName   = "תנאי-כיסוי-מלא.pdf"
)

type SendVoucherParams struct {
	RecipientEmail     string
	BookingReferenceID string
	DriverFullName     string
	VoucherNumber      string
	VoucherHTML        string
	Broker             VoucherBroker
}

// encore:api private
func (s *Service) SendVoucher(ctx context.Context, p SendVoucherParams) error {
	pdfBytes, err := s.pdfConverter.ConvertHTMLToPDF(p.VoucherHTML)
	if err != nil {
		rlog.Error("converting voucher html to pdf", "error", err, "voucher", p.VoucherNumber)
		return fmt.Errorf("converting voucher to PDF: %w", err)
	}

	attachments := []email.Attachment{
		{
			Filename: fmt.Sprintf("voucher_%s.pdf", p.VoucherNumber),
			Reader:   bytes.NewReader(pdfBytes),
		},
	}

	staticAttachments, err := brokerAttachments(p.Broker)
	if err != nil {
		rlog.Error("loading broker attachments", "error", err, "broker", p.Broker)
		return fmt.Errorf("loading broker attachments: %w", err)
	}
	attachments = append(attachments, staticAttachments...)

	return email.SendEmail(
		ctx,
		s.reservationsEmailSender,
		[]string{p.RecipientEmail},
		fmt.Sprintf("Attached AI Booking Travel voucher number %s %s", p.BookingReferenceID, p.DriverFullName),
		email.VoucherEmailTemplate,
		email.VoucherEmailData{VoucherNumber: p.VoucherNumber},
		attachments,
	)
}

func brokerAttachments(b VoucherBroker) ([]email.Attachment, error) {
	switch b {
	case VoucherBrokerFlex:
		terms, err := assetsFS.ReadFile("assets/flex-terms.pdf")
		if err != nil {
			return nil, fmt.Errorf("reading flex-terms.pdf: %w", err)
		}
		erp, err := assetsFS.ReadFile("assets/flex-erp-letter.pdf")
		if err != nil {
			return nil, fmt.Errorf("reading flex-erp-letter.pdf: %w", err)
		}
		return []email.Attachment{
			{Filename: TermsFileName, Reader: bytes.NewReader(terms)},
			{Filename: ERPFileName, Reader: bytes.NewReader(erp)},
		}, nil
	case VoucherBrokerHertz:
		// No static attachments for Hertz yet.
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown broker: %s", b)
	}
}
