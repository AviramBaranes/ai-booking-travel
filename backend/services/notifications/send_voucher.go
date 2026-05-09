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
)

type SendVoucherRequest struct {
	RecipientEmail string
	VoucherNumber  string
	VoucherHTML    string
	Broker         VoucherBroker
}

// encore:api private
func (s *Service) SendVoucher(ctx context.Context, p SendVoucherRequest) error {
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
		s.emailSender,
		[]string{p.RecipientEmail},
		fmt.Sprintf("שובר השכרת רכב – %s", p.VoucherNumber),
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
			{Filename: "flex-terms.pdf", Reader: bytes.NewReader(terms)},
			{Filename: "flex-erp-letter.pdf", Reader: bytes.NewReader(erp)},
		}, nil
	case VoucherBrokerHertz:
		// No static attachments for Hertz yet.
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown broker: %s", b)
	}
}

