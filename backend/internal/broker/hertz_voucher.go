package broker

// GenerateVoucher renders our own voucher, since Hertz has no voucher of its own to fetch.
func (h *Hertz) GenerateVoucher(d *VoucherData) (string, error) {
	return renderVoucher(d)
}
