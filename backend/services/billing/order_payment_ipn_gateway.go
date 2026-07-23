package billing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"encore.app/internal/icount"
	"encore.app/services/accounts"
	contact "encore.app/services/accounts/handlers/contact"
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation"
	"encore.app/services/reservation/handlers/actions"
	"encore.dev/rlog"
)

const (
	OrderPaymentIPNGatewayPath = "/order-payment-ipn-gateway"
)

// encore:api public raw method=POST path=/order-payment-ipn-gateway
func OrderPaymentIPNGateway(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqData, err := parseRequest(r)
	if err != nil {
		rlog.Error("failed to parse request", "error", err)
		http.NotFound(w, r)
		return
	}

	orderID := reqData.ID
	ic := icount.NewIcount()
	transaction, err := getTransaction(ic, reqData.confirmationCode)
	if err != nil {
		rlog.Error("failed to get transaction", "error", err, "orderID", orderID)
		sendInvoiceCreationFailureEmail(ctx, orderID, err)
		http.NotFound(w, r)
		return
	}

	billingReservation, err := voucherReservation(ctx, orderID, reqData, transaction)
	if err != nil {
		rlog.Error("failed to get billing reservation after payment", "error", err, "orderID", orderID)
		sendVoucherReservationFailureEmail(ctx, orderID, err)
		return
	}

	clientID, err := getIcountClientID(ctx, reqData)
	if err != nil {
		rlog.Error("failed to get iCount client ID", "error", err, "orderID", orderID)
		sendInvoiceCreationFailureEmail(ctx, orderID, err)
		return
	}

	docNum, err := createInvoice(ic, billingReservation, clientID, transaction)
	if err != nil {
		rlog.Error("failed to create invoice", "error", err, "orderID", orderID)
		sendInvoiceCreationFailureEmail(ctx, orderID, err)
	} else {
		if err := reservation.SaveInvoiceDocNum(ctx, actions.SaveInvoiceDocNumParams{
			ID:     billingReservation.ID,
			DocNum: docNum,
		}); err != nil {
			rlog.Error("failed to save invoice doc number", "error", err, "orderID", orderID)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func parseBillingEntity(query url.Values) (officeID *int64, organizationID *int64, err error) {
	officeIDStr := query.Get("office_id")
	if officeIDStr != "" {
		officeIDVal, err := strconv.ParseInt(officeIDStr, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid office_id: %w", err)
		}
		officeID = &officeIDVal
		return officeID, nil, nil
	}

	organizationIDStr := query.Get("organization_id")
	if organizationIDStr != "" {
		organizationIDVal, err := strconv.ParseInt(organizationIDStr, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid organization_id: %w", err)
		}
		organizationID = &organizationIDVal
	}

	return officeID, organizationID, nil
}

func sendInvoiceCreationFailureEmail(ctx context.Context, reservationID int64, err error) {
	if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
		Subject: "Failed to create invoice after payment",
		Message: fmt.Sprintf("failed to create invoice after payment: reservationId: %v, error: %v", reservationID, err),
	}); publishErr != nil {
		rlog.Error("failed to publish critical error email event", "reservationId", reservationID, "error", publishErr)
	}
}

func sendVoucherReservationFailureEmail(ctx context.Context, reservationID int64, err error) {
	if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
		Subject: "Failed to voucher reservation after payment",
		Message: fmt.Sprintf("failed to voucher reservation after payment: reservationId: %v, error: %v", reservationID, err),
	}); publishErr != nil {
		rlog.Error("failed to publish critical error email event", "reservationId", reservationID, "error", publishErr)
	}
}

type ipnReqData struct {
	ID               int64
	confirmationCode string
	CustomerEmail    string
	officeID         *int64
	organizationID   *int64
}

// parseRequest parses the incoming IPN request and extracts the necessary data for processing.
func parseRequest(r *http.Request) (*ipnReqData, error) {
	query := r.URL.Query()

	orderID, err := strconv.ParseInt(query.Get("id"), 10, 64)
	if err != nil {
		rlog.Error("invalid order id in query parameters", "error", err)
		return nil, err
	}

	if err := r.ParseForm(); err != nil {
		rlog.Error("failed to parse form", "error", err)
		return nil, err
	}

	cEmail := r.Form.Get("customer_email")
	confCode := r.Form.Get("confirmation_code")
	if confCode == "" {
		rlog.Error("missing confirmation code in form parameters")
		return nil, fmt.Errorf("missing confirmation code")
	}

	offID, orgID, err := parseBillingEntity(query)
	if err != nil {
		rlog.Error("failed to parse billing entity from query parameters", "error", err)
		return nil, err
	}

	return &ipnReqData{
		ID:               orderID,
		CustomerEmail:    cEmail,
		confirmationCode: confCode,
		officeID:         offID,
		organizationID:   orgID,
	}, nil
}

// getTransaction retrieves the transaction details from iCount using the provided confirmation code.
func getTransaction(ic *icount.Icount, confirmationCode string) (*icount.Transaction, error) {
	resp, err := ic.GetTransactions(confirmationCode)
	if err != nil {
		rlog.Error("failed to get transactions", "error", err)
		return nil, err
	}

	if !resp.Status || resp.ResultsCount == 0 || len(resp.ResultsList) == 0 {
		rlog.Error("transaction not successful", "reason", resp.Reason)
		return nil, fmt.Errorf("transaction not successful: %s", resp.Reason)
	}

	transaction := resp.ResultsList[0]
	return &transaction, nil
}

// voucherReservation voucher the reservation, updates its currency rate and statuses and sends back the billing reservation after payment using the provided request data and transaction details.
func voucherReservation(ctx context.Context, id int64, reqData *ipnReqData, transaction *icount.Transaction) (reservation.BillingReservation, error) {
	email := reqData.CustomerEmail
	if email == "" {
		email = transaction.ClientEmail
	}

	resp, err := reservation.VoucherReservationAfterPayment(ctx, reservation.VoucherReservationAfterPaymentParams{
		ReservationID:           id,
		PaymentConfirmationCode: reqData.confirmationCode,
		PaymentDocNum:           transaction.DocNumber,
		UserEmail:               email,
	})
	if err != nil {
		rlog.Error("failed to voucher reservation after payment", "error", err, "orderID", id)
		return reservation.BillingReservation{}, err
	}

	return resp.BillingReservation, nil
}

// getIcountClientID retrieves the iCount client ID based on the provided request data using the accounts service.
func getIcountClientID(ctx context.Context, reqData *ipnReqData) (int, error) {
	icountClientRes, err := accounts.GetIcountClientID(ctx, contact.GetIcountClientIDParams{
		OfficeID:       reqData.officeID,
		OrganizationID: reqData.organizationID,
	})
	if err != nil {
		rlog.Error("failed to get iCount client ID", "error", err, "orderID", reqData.ID)
		return 0, err
	}

	return int(icountClientRes.ClientID), nil
}

// createInvoice creates an invoice in iCount for the given billing reservation and transaction details.
func createInvoice(ic *icount.Icount, reservation reservation.BillingReservation, clientID int, transaction *icount.Transaction) (string, error) {
	currencyID, _ := icount.CurrencyIDsMap[reservation.CurrencyCode]
	items := buildReservationInvoiceItems(reservation, currencyID, reservation.CurrencyRate)
	invoiceRes, err := ic.CreateInvoice(icount.CreateDocParams{
		ClientID:      clientID,
		CurrencyID:    currencyID,
		Rate:          reservation.CurrencyRate,
		PaymentMethod: transaction.ToCCPayment(true),
		Items:         items,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create invoice in iCount: %w", err)
	}

	return parseBillingResponse(invoiceRes)
}
