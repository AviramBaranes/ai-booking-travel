package actions

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// ErrReservationNotVouchered is returned when the caller asks for the voucher of a reservation that
// does not have one: a booked reservation is not vouchered yet, and a canceled one never will be.
var ErrReservationNotVouchered = api_errors.NewErrorWithDetail(
	errs.FailedPrecondition,
	"Reservation is not vouchered",
	api_errors.ErrorDetails{Code: api_errors.CodeReservationNotVouchered},
)

// Voucher is a rendered voucher PDF together with the booking reference it was issued under.
type Voucher struct {
	BookingReferenceID string
	PDF                []byte
}

// DownloadVoucherRaw writes the voucher PDF to the response as a file download. It exists so the raw
// endpoint stays a thin shell around DownloadVoucher, which is where the logic lives.
func (s *ActionService) DownloadVoucherRaw(w http.ResponseWriter, req *http.Request, id int64) {
	voucher, err := s.DownloadVoucher(req.Context(), id)
	if err != nil {
		errs.HTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(voucher.PDF)))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "voucher_"+voucher.BookingReferenceID+".pdf"))

	if _, err := w.Write(voucher.PDF); err != nil {
		rlog.Error("writing voucher pdf response", "error", err, "id", id)
	}
}

// DownloadVoucher renders the voucher PDF for one of the caller's own reservations. The voucher is
// rebuilt from the reservation on every call rather than stored, the same way the copy emailed when
// the voucher was first applied is built.
func (s *ActionService) DownloadVoucher(ctx context.Context, id int64) (*Voucher, error) {
	authData := auth.Data().(*accounts.AuthData)
	reservation, err := s.getReservation(ctx, id, authData.UserID)
	if err != nil {
		return nil, err
	}

	if reservation.ReservationStatus != db.ReservationStatusVouchered {
		return nil, ErrReservationNotVouchered
	}

	b, err := getVoucherProvider(reservation.Broker)
	if err != nil {
		rlog.Error("getting voucher provider", "error", err, "id", id, "broker", reservation.Broker)
		return nil, api_errors.ErrInternalError
	}

	voucherData, err := toVoucherData(reservation)
	if err != nil {
		rlog.Error("converting reservation to voucher data", "error", err, "id", id)
		return nil, api_errors.ErrInternalError
	}

	htmlVoucher, err := b.GenerateVoucher(voucherData)
	if err != nil {
		rlog.Error("generating voucher html", "error", err, "id", id)
		return nil, api_errors.ErrInternalError
	}

	pdfBytes, err := s.pdfConverter.ConvertHTMLToPDF(htmlVoucher)
	if err != nil {
		rlog.Error("converting voucher html to pdf", "error", err, "id", id)
		return nil, api_errors.ErrInternalError
	}

	return &Voucher{BookingReferenceID: reservation.BrokerReservationID, PDF: pdfBytes}, nil
}
