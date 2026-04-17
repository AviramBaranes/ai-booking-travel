package reservation

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/beta/auth"
	"encore.dev/rlog"
)

// ApplyVoucherRequest is the request payload type for the apply voucher EP
type ApplyVoucherRequest struct {
	Voucher string `json:"voucher" validate:"required,notblank"`
}

func (r ApplyVoucherRequest) Validate() error {
	return validation.ValidateStruct(r)
}

// ApplyVoucher is the EP for applying a voucher on an agent order
//
// encore:api auth method=POST path=/reservations/:id/voucher tag:agent
func (s *Service) ApplyVoucher(ctx context.Context, id int64, p ApplyVoucherRequest) error {
	authData := auth.Data().(*accounts.AuthData)
	modifiedRows, err := s.query.ApplyVoucher(ctx, db.ApplyVoucherParams{
		ID:            id,
		UserID:        authData.UserID,
		VoucherNumber: &p.Voucher,
	})

	if err != nil {
		rlog.Error("applying voucher", "error", err, "id", id, "voucher", p.Voucher)
		return api_errors.ErrInternalError
	}

	if modifiedRows == 0 {
		return api_errors.ErrNotFound
	}

	return nil
}
