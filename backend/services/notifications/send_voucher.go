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

const (
	TermsFileName = "תנאים-כלליים.pdf"
	ERPFileName   = "תנאי-כיסוי-מלא.pdf"
)

type SendVoucherParams struct {
	RecipientEmail     string
	BookingReferenceID string
	DriverFullName     string
	VoucherNumber      string
	VoucherHTML        string
	// IncludeERP attaches our full coverage (ERP) terms, which only apply when it was purchased.
	IncludeERP bool
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

	staticAttachments, err := voucherAttachments(p.IncludeERP)
	if err != nil {
		rlog.Error("loading voucher attachments", "error", err, "voucher", p.VoucherNumber)
		return fmt.Errorf("loading voucher attachments: %w", err)
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

// voucherAttachments returns the static documents sent with every voucher, whatever the broker:
// our terms always, and the full coverage terms when the reservation includes it.
func voucherAttachments(includeERP bool) ([]email.Attachment, error) {
	terms, err := assetsFS.ReadFile("assets/terms.pdf")
	if err != nil {
		return nil, fmt.Errorf("reading terms.pdf: %w", err)
	}
	attachments := []email.Attachment{{Filename: TermsFileName, Reader: bytes.NewReader(terms)}}

	if includeERP {
		erp, err := assetsFS.ReadFile("assets/full-coverage.pdf")
		if err != nil {
			return nil, fmt.Errorf("reading full-coverage.pdf: %w", err)
		}
		attachments = append(attachments, email.Attachment{Filename: ERPFileName, Reader: bytes.NewReader(erp)})
	}

	return attachments, nil
}
