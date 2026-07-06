package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"encore.app/internal/broker"
	"encore.app/internal/icount"
	"encore.app/services/billing/db"
	"encore.app/services/booking"
	"encore.app/services/booking/handlers/booking_handlers"
	emailevents "encore.app/services/notifications/events"
	"encore.dev/rlog"
)

const (
	CustomerPaymentIPNGatewayPath = "/customer-payment-ipn-gateway"
)

// encore:api public raw method=POST path=/customer-payment-ipn-gateway
func (s *Service) CustomerPaymentIPNGateway(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqData, err := parseRequest(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pp, err := s.query.GetPendingPaymentByID(ctx, reqData.ID)
	if err != nil {
		rlog.Error("failed to get pending payment by ID", "error", err, "pendingPaymentID", reqData.ID)
		http.NotFound(w, r)
		return
	}

	ic := icount.NewIcount()
	transaction, err := getTransaction(ic, reqData.confirmationCode)
	if err != nil {
		rlog.Error("failed to get transaction from iCount", "error", err, "confirmationCode", reqData.confirmationCode)
		sendBookingFailureEmail(ctx, &pp, err)
		http.NotFound(w, r)
		return
	}

	var SelectedAddons []broker.SelectAddOn
	err = json.Unmarshal(pp.SelectedAddons, &SelectedAddons)
	if err != nil {
		rlog.Error("failed to unmarshal selected addons", "error", err, "pendingPaymentID", reqData.ID)
	}

	resp, err := booking.CustomerBook(ctx, booking_handlers.CustomerBookParams{
		BookParams: booking_handlers.BookParams{
			SnapshotID:      pp.SnapshotID,
			RateQualifier:   pp.RateQualifier,
			SupplierCode:    pp.SupplierCode,
			PlanID:          pp.PlanID,
			IncludeERP:      pp.IncludeErp,
			SelectedAddOns:  SelectedAddons,
			DriverTitle:     pp.DriverTitle,
			DriverFirstName: pp.DriverFirstName,
			DriverLastName:  pp.DriverLastName,
			FlightNumber:    pp.FlightNumber,
		},
		UserID: pp.UserID,
	})
	if err != nil {
		rlog.Error("failed to book car after payment", "error", err, "pendingPaymentID", pp.ID)
		sendBookingFailureEmail(ctx, &pp, err)
		refundPayment(transaction, "failed to book car after payment")
	}

	if err = s.query.ResolvePendingPayment(ctx, db.ResolvePendingPaymentParams{
		ID:            pp.ID,
		ReservationID: &resp.ReservationID,
	}); err != nil {
		rlog.Error("failed to resolve pending payment", "error", err, "pendingPaymentID", pp.ID)
		sendBookingFailureEmail(ctx, &pp, err) // continue without resolving pending payment
	}

	BillingReservation, err := getBillingReservation(ctx, resp.ReservationID, reqData, transaction)
	if err != nil {
		rlog.Error("failed to get billing reservation after payment", "error", err, "reservationID", resp.ReservationID)
		sendVoucherReservationFailureEmail(ctx, resp.ReservationID, err)
		return
	}

	err = createInvoice(ic, BillingReservation, 0, transaction)
	if err != nil {
		sendInvoiceCreationFailureEmail(ctx, resp.ReservationID, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendBookingFailureEmail(ctx context.Context, pp *db.PendingCustomerPayment, err error) {
	if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
		Subject: "Failed to book car after payment",
		Message: fmt.Sprintf("failed to book car after payment: customerID: %d, driver: %s %s, snapshotID: %d, rateQualifier: %s, supplierCode: %s, planID: %s, error: %v",
			pp.UserID, pp.DriverFirstName, pp.DriverLastName, pp.SnapshotID, pp.RateQualifier, pp.SupplierCode, pp.PlanID, err),
	}); publishErr != nil {
		rlog.Error("failed to publish critical error email event", "pendingPaymentID", pp.ID, "error", publishErr)
	}
}

func refundPayment(t *icount.Transaction, reason string) {
	ic := icount.NewIcount()
	_, err := ic.CancelDocument(icount.CancelDocumentParams{
		DocType:  t.DocType,
		DocNum:   t.DocNumber,
		RefundCC: true,
		Reason:   reason,
	})
	if err != nil {
		rlog.Error("failed to refund payment", "error", err, "docNum", t.DocNumber)
	}
}
