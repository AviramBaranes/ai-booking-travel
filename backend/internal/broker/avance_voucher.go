package broker

// GenerateVoucher renders our own voucher, since Wheelsys has no voucher of its own to fetch.
func (a *Avance) GenerateVoucher(d *VoucherData) (string, error) {
	return renderVoucher(d)
}
