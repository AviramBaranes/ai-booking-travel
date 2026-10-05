package broker

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed voucher.html
var voucherTemplateHTML string

var voucherTemplate = template.Must(template.New("voucher").Parse(voucherTemplateHTML))

// renderVoucher renders our own voucher, for the brokers that have no voucher of their own to
// fetch. Every such broker shares it, so a voucher looks the same whoever supplied the car.
func renderVoucher(d *VoucherData) (string, error) {
	var htmlVoucher bytes.Buffer
	if err := voucherTemplate.Execute(&htmlVoucher, d); err != nil {
		return "", fmt.Errorf("executing voucher html template %w", err)
	}

	return htmlVoucher.String(), nil
}
