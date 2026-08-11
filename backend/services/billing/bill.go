package billing

import (
	"context"
	"fmt"

	"encore.app/internal/api_errors"
	"encore.app/internal/icount"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	contact "encore.app/services/accounts/handlers/contact"
	"encore.app/services/accounts/handlers/office"
	"encore.app/services/accounts/handlers/organization"
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation"
	"encore.dev/beta/errs"
	"encore.dev/config"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

type invoiceConfig struct {
	PurchaseItemDescription            config.String
	ProfitItemDescription              config.String
	ProfitAndErpItemDescription        config.String
	CancellationPenaltyItemDescription config.String
	NoShowPenaltyItemDescription       config.String
	DiscountItemDescription            config.String
}

var (
	ErrExactlyOneOfOfficeIDOrOrgIDRequired = api_errors.NewValidationError("exactly one of office_id or organization_id must be provided")
	ErrInvalidReservationID                = api_errors.NewValidationError("one or more provided IDs do not belong to the specified billing entity")
	ErrMismatchedCurrencies                = api_errors.NewValidationError("all selected reservations must have the same currency")
	ErrNothingToBill                       = api_errors.NewValidationError("at least one reservation or penalty must be provided")
	emailRequestedTopic                    = pubsub.TopicRef[pubsub.Publisher[*emailevents.EmailEvent]](emailevents.EmailRequestedTopic)
	emailPublisher                         = emailevents.NewPublisher(emailRequestedTopic)
)

type BillParams struct {
	IDs                 []int64 `json:"ids"`
	PenaltyIDs          []int64 `json:"penalty_ids"`
	SkipInvoiceCreation bool    `json:"skip_invoice_creation" encore:"optional"`
	TotalPaid           float64 `json:"total_paid" validate:"required,gt=0"`
	TransferDate        string  `json:"transfer_date" validate:"required,datetime=2006-01-02"`
	OfficeID            *int64  `json:"office_id" encore:"optional"`
	OrganizationID      *int64  `json:"organization_id" encore:"optional"`
	// Deduction (tax withheld at source) and Discount (stated VAT-inclusive) are given in the
	// billed currency and appear on the invoice only — TotalPaid is already net of both.
	Deduction float64 `json:"deduction" encore:"optional" validate:"omitempty,gte=0"`
	Discount  float64 `json:"discount" encore:"optional" validate:"omitempty,gte=0"`
}

func (r BillParams) Validate() error {
	if err := validation.ValidateStruct(r); err != nil {
		return err
	}

	if r.OfficeID == nil && r.OrganizationID == nil || r.OfficeID != nil && r.OrganizationID != nil {
		return ErrExactlyOneOfOfficeIDOrOrgIDRequired
	}

	// A bill may cover reservations, penalties, or both — but not nothing.
	if len(r.IDs) == 0 && len(r.PenaltyIDs) == 0 {
		return ErrNothingToBill
	}

	return nil
}

type BillResponse struct {
	InvoiceDocNum string `json:"invoiceDocNum"`
	ReceiptDocNum string `json:"receiptDocNum"`
}

// encore:api auth method=POST path=/bill tag:accountant
func (s *Service) Bill(ctx context.Context, p BillParams) (*BillResponse, error) {
	icountClientRes, err := accounts.GetIcountClientID(ctx, contact.GetIcountClientIDParams{
		OfficeID:       p.OfficeID,
		OrganizationID: p.OrganizationID,
	})
	if err != nil {
		rlog.Error("failed to get iCount client ID for billing entity", "error", err, "office_id", p.OfficeID, "org_id", p.OrganizationID)
		return nil, api_errors.ErrInternalError
	}

	openReservations, err := reservation.ListOpenReservationsByBillingEntity(ctx, &reservation.ListOpenReservationsByBillingEntityParams{
		OfficeID: derefInt64(p.OfficeID),
		OrgID:    derefInt64(p.OrganizationID),
	})

	if err != nil {
		rlog.Error("failed to list open reservations for billing entity", "error", err, "office_id", p.OfficeID, "org_id", p.OrganizationID)
		return nil, err
	}

	reservationSet := createReservationSet(openReservations.CurrencyGroups)
	penaltySet := createPenaltySet(openReservations.CurrencyGroups)

	if err := validateIDsBelongToBillingEntity(p.IDs, reservationSet); err != nil {
		rlog.Error("validation failed for billing request", "error", err, "invalid_ids", p.IDs)
		return nil, err
	}

	if err := validateIDsBelongToBillingEntity(p.PenaltyIDs, penaltySet); err != nil {
		rlog.Error("validation failed for billing request", "error", err, "invalid_penalty_ids", p.PenaltyIDs)
		return nil, err
	}

	currency, err := resolveBillCurrency(p, reservationSet, penaltySet)
	if err != nil {
		rlog.Error("validation failed for billing request", "error", err, "invalid_ids", p.IDs, "invalid_penalty_ids", p.PenaltyIDs)
		return nil, err
	}

	total := calculateTotalAmount(reservationSet, penaltySet, p.IDs, p.PenaltyIDs)
	if p.SkipInvoiceCreation {
		if err := resolveBilledItems(ctx, p, total); err != nil {
			return nil, err
		}

		return &BillResponse{
			InvoiceDocNum: "N/A",
			ReceiptDocNum: "N/A",
		}, nil
	}

	invcDocNum, err := createInvoiceDoc(p, reservationSet, penaltySet, icountClientRes.ClientID, currency, cfg.Icount.AccountID())
	if err != nil {
		rlog.Error("failed to create invoice in iCount", "error", err, "client_id", icountClientRes.ClientID, "currency", currency, "total_paid", p.TotalPaid)
		return nil, err
	}

	recptDocNum, err := createReceiptDoc(p, reservationSet, penaltySet, icountClientRes.ClientID, currency, cfg.Icount.AccountID())
	if err != nil {
		rlog.Error("failed to create receipt in iCount", "error", err, "client_id", icountClientRes.ClientID, "currency", currency, "total_paid", p.TotalPaid)
		if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
			Subject: "Failed to resolve reservations after invoice creation",
			Message: fmt.Sprintf("failed to create receipt and resolve reservations after successful invoice creation, reservation_ids: %v, invoiceDocNum: %v, error: %v", p.IDs, invcDocNum, err),
		}); publishErr != nil {
			rlog.Error("failed to publish critical error email event", "reservation_ids", p.IDs, "invoice doc num", invcDocNum, "error", publishErr)
		}
		return nil, err
	}

	if err := resolveBilledItems(ctx, p, total); err != nil {
		return nil, err
	}

	return &BillResponse{
		InvoiceDocNum: invcDocNum,
		ReceiptDocNum: recptDocNum,
	}, nil
}

// validateIDsBelongToBillingEntity checks if all provided IDs belong to the specified billing entity (office or organization) by verifying their presence in the set of items open for that entity. If any ID does not belong to the billing entity, it returns a validation error.
func validateIDsBelongToBillingEntity[T any](ids []int64, set map[int64]T) error {
	for _, id := range ids {
		if _, exists := set[id]; !exists {
			return ErrInvalidReservationID
		}
	}

	return nil
}

// buildInvoiceItems constructs a list of ICountInvoiceItem based on the provided reservation and penalty IDs and their corresponding entries in the sets.
// For each reservation ID, it creates two invoice items: one for the car purchase price (free of tax) and another for the profit + optionally ERP selling price (both requires tax).
// Each penalty ID adds a single tax-exempt item, as the supplier's fee is passed on to the customer with no profit.
func buildInvoiceItems(ids, penaltyIDs []int64, reservationsSet reservationSet, penaltiesSet penaltySet, discount float64) []icount.ICountInvoiceItem {
	invoiceItems := make([]icount.ICountInvoiceItem, 0, len(ids)*2+len(penaltyIDs)+1)
	for _, id := range ids {
		reservation := reservationsSet[id]
		rItems := buildReservationInvoiceItems(reservation, 0, 0) //currencyID and currencyRate omitted
		invoiceItems = append(invoiceItems, rItems...)
	}

	for _, id := range penaltyIDs {
		invoiceItems = append(invoiceItems, buildPenaltyInvoiceItem(penaltiesSet[id], 0, 0)) //currencyID and currencyRate omitted
	}

	if discount > 0 {
		invoiceItems = append(invoiceItems, discountInvoiceItem(discount))
	}

	return invoiceItems
}

// discountInvoiceItem states the discount as a negative line priced the way the profit it comes out
// of is priced — VAT-inclusive and taxable. iCount's document-level discount field cannot be used:
// it spreads the discount pro-rata across the exempt and taxable lines, so most of it lands on the
// tax-free purchase price and the invoice ends up above what the receipt collects.
func discountInvoiceItem(discount float64) icount.ICountInvoiceItem {
	return icount.ICountInvoiceItem{
		Description:     cfg.Invoice.DiscountItemDescription(),
		Quantity:        1,
		IsTaxExempt:     false,
		UnitPriceIncvat: floatPtr(-discount),
	}
}

// buildReceiptItems constructs a list of ICountInvoiceItem based on the provided reservation and penalty IDs and their corresponding entries in the sets.
func buildReceiptItems(ids, penaltyIDs []int64, reservationsSet reservationSet, penaltiesSet penaltySet) []icount.ICountInvoiceItem {
	invoiceItems := make([]icount.ICountInvoiceItem, 0, len(ids)+len(penaltyIDs))
	for _, id := range ids {
		reservation := reservationsSet[id]
		m := 1.0
		if reservation.PaymentStatus == "refund_pending" {
			m = -1.0
		}
		currencyID, _ := icount.CurrencyIDsMap[reservation.CurrencyCode]
		item := icount.ICountInvoiceItem{
			Description:  reservation.BrokerReservationID,
			Quantity:     1,
			IsTaxExempt:  true,
			SKU:          fmt.Sprintf("%d", reservation.ID),
			UnitPrice:    floatPtr(reservation.TotalPrice * m),
			CurrencyID:   currencyID,
			CurrencyRate: reservation.CurrencyRate,
		}
		invoiceItems = append(invoiceItems, item)
	}

	for _, id := range penaltyIDs {
		penalty := penaltiesSet[id]
		currencyID, _ := icount.CurrencyIDsMap[penalty.CurrencyCode]
		invoiceItems = append(invoiceItems, buildPenaltyInvoiceItem(penalty, currencyID, penalty.CurrencyRate))
	}

	return invoiceItems
}

// buildPenaltyInvoiceItem constructs the single line a penalty contributes to a document. The fee is
// a pass-through of the supplier's charge, so it carries no profit and is fully tax exempt — the same
// treatment as a reservation's car purchase price.
func buildPenaltyInvoiceItem(penalty reservation.BillingPenalty, currencyID int, currencyRate float64) icount.ICountInvoiceItem {
	return icount.ICountInvoiceItem{
		Description:  fmt.Sprintf("%s %s", penaltyItemDescription(penalty.Type), penalty.BrokerReservationID),
		Quantity:     1,
		IsTaxExempt:  true,
		SKU:          penaltySKU(penalty.ID),
		UnitPrice:    floatPtr(penalty.Amount),
		CurrencyID:   currencyID,
		CurrencyRate: currencyRate,
	}
}

// penaltyItemDescription returns the configured document description for a penalty type.
func penaltyItemDescription(penaltyType string) string {
	if penaltyType == reservation.PenaltyTypeNoShow {
		return cfg.Invoice.NoShowPenaltyItemDescription()
	}

	return cfg.Invoice.CancellationPenaltyItemDescription()
}

// penaltySKU namespaces a penalty id, since reservations and penalties are numbered independently
// and may otherwise collide on the same document.
func penaltySKU(id int64) string {
	return fmt.Sprintf("P%d", id)
}

func buildReservationInvoiceItems(reservation reservation.BillingReservation, currencyID int, currencyRate float64) []icount.ICountInvoiceItem {
	invoiceItems := make([]icount.ICountInvoiceItem, 0, 2)
	m := 1.0
	if reservation.PaymentStatus == "refund_pending" {
		m = -1.0
	}

	invoiceItems = append(invoiceItems, icount.ICountInvoiceItem{
		Description:  cfg.Invoice.PurchaseItemDescription(),
		UnitPrice:    floatPtr(reservation.CarPurchasePrice * m),
		Quantity:     1,
		IsTaxExempt:  true,
		SKU:          fmt.Sprintf("%d", reservation.ID),
		CurrencyID:   currencyID,
		CurrencyRate: currencyRate,
	})

	profitDesc := cfg.Invoice.ProfitItemDescription()
	if reservation.ERPSellingPrice > 0 {
		profitDesc = cfg.Invoice.ProfitAndErpItemDescription()
	}
	invoiceItems = append(invoiceItems, icount.ICountInvoiceItem{
		Description:     profitDesc,
		UnitPriceIncvat: floatPtr((reservation.TotalProfit) * m),
		Quantity:        1,
		IsTaxExempt:     false,
		SKU:             fmt.Sprintf("%d", reservation.ID),
		CurrencyID:      currencyID,
	})

	return invoiceItems
}

func floatPtr(f float64) *float64 {
	return &f
}

// reservationSet is a helper type for efficient lookup of reservations by ID when building invoice items.
type reservationSet map[int64]reservation.BillingReservation

// penaltySet is a helper type for efficient lookup of penalties by ID when building invoice items.
type penaltySet map[int64]reservation.BillingPenalty

// createReservationSet creates a set of reservations indexed by their ID for efficient lookup.
func createReservationSet(currencyGroups []reservation.CurrencyGroup) map[int64]reservation.BillingReservation {
	reservationSet := make(map[int64]reservation.BillingReservation)
	for _, group := range currencyGroups {
		for _, r := range group.Reservations {
			reservationSet[r.ID] = r
		}
	}

	return reservationSet
}

// createPenaltySet creates a set of penalties indexed by their ID for efficient lookup.
func createPenaltySet(currencyGroups []reservation.CurrencyGroup) map[int64]reservation.BillingPenalty {
	penaltySet := make(map[int64]reservation.BillingPenalty)
	for _, group := range currencyGroups {
		for _, p := range group.Penalties {
			penaltySet[p.ID] = p
		}
	}

	return penaltySet
}

// resolveBillCurrency returns the currency the document is issued in, which every selected
// reservation and penalty must share since a document covers a single currency.
func resolveBillCurrency(p BillParams, reservationSet reservationSet, penaltySet penaltySet) (string, error) {
	currencies := make([]string, 0, len(p.IDs)+len(p.PenaltyIDs))
	for _, id := range p.IDs {
		currencies = append(currencies, reservationSet[id].CurrencyCode)
	}
	for _, id := range p.PenaltyIDs {
		currencies = append(currencies, penaltySet[id].CurrencyCode)
	}

	for _, currency := range currencies[1:] {
		if currency != currencies[0] {
			return "", ErrMismatchedCurrencies
		}
	}

	return currencies[0], nil
}

// parseBillingResponse converts the response from the iCount service into a BillResponse. If the iCount response indicates failure, it logs the error details and returns a generic error with the combined error messages from iCount.
func parseBillingResponse(result *icount.ICountCreateDocResponse) (string, error) {
	if result.Status {
		return result.DocNum, nil
	} else {
		rlog.Error("icount respond with an error", "reason", result.Reason, "error_description", result.ErrorDescription)
		var errMsg string
		for _, detail := range result.ErrorDetails {
			errMsg += detail
		}

		return "", api_errors.NewErrorWithDetail(errs.Unknown, errMsg, api_errors.EmptyDetails)
	}
}

// resolveBilledItems resolves everything the bill covers: first the penalties, then the
// reservations, which also updates the billing entity's balance by the total of both.
func resolveBilledItems(ctx context.Context, p BillParams, total float64) error {
	if err := resolvePenalties(ctx, p); err != nil {
		return err
	}

	return resolveReservations(ctx, p, total)
}

// resolvePenalties marks the billed penalties as collected from the customer. As in
// resolveReservations, a failure that happens after the money already moved is alerted on rather
// than returned, so the caller does not retry a bill that has already been issued.
func resolvePenalties(ctx context.Context, p BillParams) error {
	if len(p.PenaltyIDs) == 0 {
		return nil
	}

	if err := reservation.ResolvePenalties(ctx, reservation.ResolvePenaltiesParams{
		IDs: p.PenaltyIDs,
	}); err != nil {
		if !p.SkipInvoiceCreation {
			rlog.Error("failed to resolve penalties after successful billing", "error", err, "penalty_ids", p.PenaltyIDs)
			if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
				Subject: "Failed to resolve penalties after billing",
				Message: fmt.Sprintf("failed to resolve penalties after successful billing, penalty_ids: %v, error: %v", p.PenaltyIDs, err),
			}); publishErr != nil {
				rlog.Error("failed to publish critical error email event", "penalty_ids", p.PenaltyIDs, "error", publishErr)
			}
			return nil
		}
		rlog.Error("failed to resolve penalties", "error", err, "penalty_ids", p.PenaltyIDs)
		return api_errors.ErrInternalError
	}

	return nil
}

// resolveReservations attempts to resolve the reservations associated with the billed invoice.
// If resolving fails after a successful billing, it logs the error and sends a critical error email notification.
// It also updates the balance due for the associated office or organization based on the billing entity specified in the BillParams.
func resolveReservations(ctx context.Context, p BillParams, total float64) error {
	var err error
	if len(p.IDs) > 0 {
		err = reservation.ResolveReservations(ctx, reservation.ResolveReservationsParams{
			IDs: p.IDs,
		})
	}
	if err != nil {
		if !p.SkipInvoiceCreation {
			rlog.Error("failed to resolve reservations after successful billing", "error", err, "reservation_ids", p.IDs)
			if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
				Subject: "Failed to resolve reservations after billing",
				Message: fmt.Sprintf("failed to resolve reservations after successful billing, reservation_ids: %v, error: %v", p.IDs, err),
			}); publishErr != nil {
				rlog.Error("failed to publish critical error email event", "reservation_ids", p.IDs, "error", publishErr)
			}
			return nil
		}
		rlog.Error("failed to resolve reservations", "error", err, "reservation_ids", p.IDs)
		return api_errors.ErrInternalError
	}

	if p.OfficeID != nil {
		if err := accounts.UpdateOfficeBalanceDue(ctx, office.UpdateOfficeBalanceDueParams{
			ID:            *p.OfficeID,
			BalanceChange: total * -1,
		}); err != nil {
			rlog.Error("failed to update office balance due after billing", "error", err, "office_id", p.OfficeID, "amount", total)
			return api_errors.ErrInternalError
		}
		return nil
	}

	if err := accounts.UpdateOrganizationBalanceDue(ctx, organization.UpdateOrganizationBalanceDueParams{
		ID:            *p.OrganizationID,
		BalanceChange: total * -1,
	}); err != nil {
		rlog.Error("failed to update organization balance due after billing", "error", err, "organization_id", p.OrganizationID, "amount", total)
		return api_errors.ErrInternalError
	}

	return nil
}

func derefInt64(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// calculateTotalAmount computes the total amount for the selected reservations and penalties by summing each one's price multiplied by its currency rate. It returns the computed total as a float64.
func calculateTotalAmount(reservationSet reservationSet, penaltySet penaltySet, ids, penaltyIDs []int64) float64 {
	var total float64
	for _, id := range ids {
		r := reservationSet[id]
		total += (r.TotalPrice * r.CurrencyRate)
	}

	for _, id := range penaltyIDs {
		p := penaltySet[id]
		total += (p.Amount * p.CurrencyRate)
	}

	return total
}

func createInvoiceDoc(p BillParams, reservationSet reservationSet, penaltySet penaltySet, clientID int32, currency string, accountID int) (string, error) {
	ic := icount.NewIcount()
	invoiceItems := buildInvoiceItems(p.IDs, p.PenaltyIDs, reservationSet, penaltySet, p.Discount)
	res, err := ic.CreateInvoice(icount.CreateDocParams{
		ClientID:      int(clientID),
		CurrencyID:    icount.CurrencyIDsMap[currency],
		PaymentMethod: icount.NewBankTransferPayment(p.TotalPaid, p.TransferDate, accountID),
		Items:         invoiceItems,
	})

	if err != nil {
		rlog.Error("failed to create invoice in iCount", "error", err, "client_id", clientID, "currency", currency, "total_paid", p.TotalPaid)
		return "", api_errors.ErrInternalError
	}

	docNum, err := parseBillingResponse(res)
	if err != nil {
		rlog.Error("failed to create invoice in iCount", "error", err, "client_id", clientID, "currency", currency, "total_paid", p.TotalPaid)
		return "", err
	}

	return docNum, nil
}

// createReceiptDoc issues the receipt for what the billing entity actually settled. iCount requires
// a receipt's payments to cover its items, so the discount is deducted from the items and the
// withholding tax is declared as a deduction — leaving exactly TotalPaid to be covered by the
// transfer.
func createReceiptDoc(p BillParams, reservationSet reservationSet, penaltySet penaltySet, clientID int32, currency string, accountID int) (string, error) {
	ic := icount.NewIcount()
	currencyID := icount.CurrencyIDsMap[currency]
	receiptItems := buildReceiptItems(p.IDs, p.PenaltyIDs, reservationSet, penaltySet)
	if p.Discount > 0 {
		receiptItems = append(receiptItems, discountReceiptItem(p.Discount, currencyID))
	}

	res, err := ic.CreateReceipt(icount.CreateDocParams{
		ClientID:      int(clientID),
		CurrencyID:    currencyID,
		PaymentMethod: icount.NewBankTransferPayment(p.TotalPaid, p.TransferDate, accountID),
		Items:         receiptItems,
		Deductions:    icount.NewWithholdingTaxDeduction(p.Deduction),
	})
	if err != nil {
		rlog.Error("failed to create receipt in iCount", "error", err, "client_id", clientID, "currency", currency, "total_paid", p.TotalPaid)
		return "", api_errors.ErrInternalError
	}

	return parseBillingResponse(res)
}

// discountReceiptItem states the discount as a negative line rather than iCount's document-level
// discount field, which is read in ILS: an item carries the document's own currency, so the
// receipt's items stay in the billed currency and need no conversion to balance against the
// transfer. Like every other receipt line it is tax exempt — a receipt records money, not VAT.
func discountReceiptItem(discount float64, currencyID int) icount.ICountInvoiceItem {
	return icount.ICountInvoiceItem{
		Description: cfg.Invoice.DiscountItemDescription(),
		Quantity:    1,
		IsTaxExempt: true,
		UnitPrice:   floatPtr(-discount),
		CurrencyID:  currencyID,
	}
}
