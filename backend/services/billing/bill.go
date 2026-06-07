package billing

import (
	"context"
	"fmt"

	"encore.app/internal/api_errors"
	emailevents "encore.app/internal/email_events"
	"encore.app/internal/icount"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	contact "encore.app/services/accounts/handlers/contact"
	"encore.app/services/reservation"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

var (
	ErrExactlyOneOfOfficeIDOrOrgIDRequired = api_errors.NewValidationError("exactly one of office_id or organization_id must be provided")
	ErrInvalidReservationID                = api_errors.NewValidationError("one or more provided IDs do not belong to the specified billing entity")
	ErrMismatchedCurrencies                = api_errors.NewValidationError("all selected reservations must have the same currency")
)

type BillParams struct {
	IDs            []int64 `json:"ids" validate:"required,min=1"`
	TotalPaid      float64 `json:"total_paid" validate:"required,gt=0"`
	TransferDate   string  `json:"transfer_date" validate:"required,datetime=2006-01-02"`
	OfficeID       *int64  `json:"office_id" encore:"optional"`
	OrganizationID *int64  `json:"organization_id" encore:"optional"`
}

func (r BillParams) Validate() error {
	if err := validation.ValidateStruct(r); err != nil {
		return err
	}

	if r.OfficeID == nil && r.OrganizationID == nil || r.OfficeID != nil && r.OrganizationID != nil {
		return ErrExactlyOneOfOfficeIDOrOrgIDRequired
	}

	return nil
}

type BillResponse struct {
	DocNum string `json:"docNum"`
}

// encore:api auth method=POST path=/bill tag:accountant
func Bill(ctx context.Context, p BillParams) (*BillResponse, error) {
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

	if err := validateIDsBelongToBillingEntity(p.IDs, reservationSet); err != nil {
		rlog.Error("validation failed for billing request", "error", err, "invalid_ids", p.IDs)
		return nil, err
	}

	currency := reservationSet[p.IDs[0]].CurrencyCode
	if err := validateSelectedIDsShareCurrency(currency, p.IDs, reservationSet); err != nil {
		rlog.Error("validation failed for billing request", "error", err, "invalid_ids", p.IDs)
		return nil, err
	}

	invoiceItems := buildInvoiceItems(p.IDs, reservationSet)
	rlog.Info("build items", "items", invoiceItems)
	ic := icount.NewIcount(cfg.Icount.CID(), cfg.Icount.User())
	res, err := ic.CreateInvoice(icount.CreateInvoiceParams{
		ClientID:   int(icountClientRes.ClientID),
		CurrencyID: icount.CurrencyIDsMap[currency],
		Sum:        p.TotalPaid,
		Date:       p.TransferDate,
		AccountID:  cfg.Icount.AccountID(),
		Items:      invoiceItems,
	})

	resp, err := parseBillingResponse(res)
	if err != nil {
		rlog.Error("failed to create invoice in iCount", "error", err, "client_id", icountClientRes.ClientID, "currency", currency, "total_paid", p.TotalPaid)
		return nil, err
	}

	err = reservation.ResolveReservations(ctx, reservation.ResolveReservationsParams{
		IDs: p.IDs,
	})
	if err != nil {
		rlog.Error("failed to resolve reservations after successful billing", "error", err, "reservation_ids", p.IDs)
		emailevents.PublishEmailEvent(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
			Subject: "Failed to resolve reservations after billing",
			Message: fmt.Sprintf("failed to resolve reservations after successful billing, reservation_ids: %v, error: %v", p.IDs, err),
		})
	}

	return resp, nil
}

// validateIDsBelongToBillingEntity checks if all provided reservation IDs belong to the specified billing entity (office or organization) by verifying their presence in the reservationSet. If any ID does not belong to the billing entity, it returns a validation error.
func validateIDsBelongToBillingEntity(ids []int64, reservationsSet reservationSet) error {
	for _, id := range ids {
		if _, exists := reservationsSet[id]; !exists {
			return ErrInvalidReservationID
		}
	}

	return nil
}

// buildInvoiceItems constructs a list of ICountInvoiceItem based on the provided reservation IDs and their corresponding reservations in the reservationSet.
// For each reservation ID, it creates two invoice items: one for the car purchase price (free of tax) and another for the profit + optionally ERP selling price (both requires tax).
func buildInvoiceItems(ids []int64, reservationsSet reservationSet) []icount.ICountInvoiceItem {
	invoiceItems := make([]icount.ICountInvoiceItem, 0, len(ids)*2)
	for _, id := range ids {
		reservation := reservationsSet[id]
		m := 1.0
		if reservation.PaymentStatus == "refund_pending" {
			m = -1.0
		}

		invoiceItems = append(invoiceItems, icount.ICountInvoiceItem{
			Description: cfg.Invoice.PurchaseItemDescription(),
			UnitPrice:   floatPtr(reservation.CarPurchasePrice * m),
			Quantity:    1,
			IsTaxExempt: true,
			SKU:         fmt.Sprintf("%d", id),
		})

		profitDesc := cfg.Invoice.ProfitItemDescription()
		if reservation.ERPSellingPrice > 0 {
			profitDesc = cfg.Invoice.ProfitAndErpItemDescription()
		}
		invoiceItems = append(invoiceItems, icount.ICountInvoiceItem{
			Description:     profitDesc,
			UnitPriceIncvat: floatPtr((reservation.ProfitOnCar + reservation.ERPSellingPrice) * m),
			Quantity:        1,
			IsTaxExempt:     false,
			SKU:             fmt.Sprintf("%d", id),
		})
	}

	return invoiceItems
}

func floatPtr(f float64) *float64 {
	return &f
}

// reservationSet is a helper type for efficient lookup of reservations by ID when building invoice items.
type reservationSet map[int64]reservation.BillingReservation

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

// validateSelectedIDsShareCurrency checks if all selected reservation IDs correspond to reservations that share the same currency.
func validateSelectedIDsShareCurrency(currency string, ids []int64, reservationSet reservationSet) error {
	for _, id := range ids[1:] {
		if reservationSet[id].CurrencyCode != currency {
			return ErrMismatchedCurrencies
		}
	}

	return nil
}

// parseBillingResponse converts the response from the iCount service into a BillResponse. If the iCount response indicates failure, it logs the error details and returns a generic error with the combined error messages from iCount.
func parseBillingResponse(result *icount.ICountCreateDocResponse) (*BillResponse, error) {
	if result.Status {
		return &BillResponse{
			DocNum: result.DocNum,
		}, nil
	} else {
		rlog.Error("icount respond with an error", "reason", result.Reason, "error_description", result.ErrorDescription)
		var errMsg string
		for _, detail := range result.ErrorDetails {
			errMsg += detail
		}

		return nil, api_errors.NewErrorWithDetail(errs.Unknown, errMsg, api_errors.EmptyDetails)
	}
}

func derefInt64(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
