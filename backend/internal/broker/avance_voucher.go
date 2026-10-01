package broker

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed avance_voucher.html
var avanceVoucherTemplate string

var avanceVoucher = template.Must(template.New("avance_voucher").Parse(avanceVoucherTemplate))

// GenerateVoucher renders our own voucher, since Wheelsys has no voucher of its own to fetch.
func (a *Avance) GenerateVoucher(d *VoucherData) (string, error) {
	var htmlVoucher bytes.Buffer
	if err := avanceVoucher.Execute(&htmlVoucher, d); err != nil {
		return "", fmt.Errorf("executing avance voucher html template %w", err)
	}

	return htmlVoucher.String(), nil
}
