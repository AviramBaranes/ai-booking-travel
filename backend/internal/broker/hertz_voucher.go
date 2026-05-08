package broker

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed hertz_voucher.html
var hertzVoucherTemplate string

func (h Hertz) GenerateVoucher(d *VoucherData) (string, error) {
	templ, err := template.New("hertz_voucher").Parse(hertzVoucherTemplate)
	if err != nil {
		return "", fmt.Errorf("generating voucher html template %w", err)
	}

	var htmlVoucher bytes.Buffer
	if err = templ.Execute(&htmlVoucher, d); err != nil {
		return "", fmt.Errorf("executing voucher html template %w", err)
	}

	return htmlVoucher.String(), nil
}
