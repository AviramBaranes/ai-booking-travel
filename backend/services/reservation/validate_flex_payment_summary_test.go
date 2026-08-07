package reservation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/handlers/supplier_payments"
	"github.com/xuri/excelize/v2"
)

// Flex booking ids are a two-character prefix followed by six digits.
const (
	bookingMatch         = "4H283419"
	bookingWrongPrice    = "4H491052"
	bookingWrongCurrency = "4P683211"
	bookingUSD           = "4P114872"
	bookingAlreadyPaid   = "4H770313"
	bookingGhost         = "4P905544"
)

// --- Helpers ---

// summaryLine is one row of a Flex payment-required summary sheet.
type summaryLine struct {
	bookingID string
	balance   any // any, so tests can write formatted cells such as "1,234.56"
}

// buildFlexSummary writes a workbook shaped like a real Flex statement: account metadata, the
// header row at 14, a currency row, a "CARS" section label, the booking lines from row 17 with
// the booking id in column B and the balance in column L, a "Total" row, and trailing notes.
func buildFlexSummary(t *testing.T, sheets map[string][]summaryLine) []byte {
	t.Helper()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("failed to close workbook: %v", err)
		}
	}()

	set := func(sheetName, cell string, value any) {
		if err := f.SetCellValue(sheetName, cell, value); err != nil {
			t.Fatalf("failed to write %s!%s: %v", sheetName, cell, err)
		}
	}

	for sheetName, lines := range sheets {
		if _, err := f.NewSheet(sheetName); err != nil {
			t.Fatalf("failed to create sheet %q: %v", sheetName, err)
		}

		// Metadata above the list; the parser must skip it.
		set(sheetName, "J1", "STATEMENT")
		set(sheetName, "I5", "AI Booking Travel LLC")
		set(sheetName, "H12", "Total Payable")

		set(sheetName, "B14", "Booking ID")
		set(sheetName, "L14", "Balance")
		set(sheetName, "I15", sheetName)
		set(sheetName, "L15", sheetName)
		set(sheetName, "B16", "CARS")

		var total float64
		for i, line := range lines {
			row := 17 + i
			set(sheetName, fmt.Sprintf("B%d", row), line.bookingID)
			set(sheetName, fmt.Sprintf("L%d", row), line.balance)
			if balance, ok := line.balance.(float64); ok {
				total += balance
			}
		}

		// The closing total row, then the free-text notes Flex appends below it. Both sit in
		// columns the parser reads, so it has to stop at the total.
		totalRow := 17 + len(lines)
		set(sheetName, fmt.Sprintf("K%d", totalRow), "Total")
		set(sheetName, fmt.Sprintf("L%d", totalRow), total)
		set(sheetName, fmt.Sprintf("B%d", totalRow+2), "Please Note:-")
		set(sheetName, fmt.Sprintf("B%d", totalRow+3), "     • If there is a date/time present against a booking, it will be cancelled.")
	}

	if err := f.DeleteSheet("Sheet1"); err != nil {
		t.Fatalf("failed to delete default sheet: %v", err)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatalf("failed to write workbook to buffer: %v", err)
	}
	return buf.Bytes()
}

// postFlexSummary uploads the given file bytes to the validation endpoint and returns the recorded
// response. Encore does not allow in-app calls to raw endpoints, so the handler is driven directly
// over httptest — which still covers the multipart extraction and the JSON encoding.
func postFlexSummary(t *testing.T, s *Service, fileName string, content []byte) *httptest.ResponseRecorder {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write form file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/supplier-payments/flex/validate", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()
	s.ValidateFlexPaymentSummary(rec, req)
	return rec
}

func decodeValidationResponse(t *testing.T, rec *httptest.ResponseRecorder) *ValidateFlexPaymentSummaryResponse {
	t.Helper()

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var decoded ValidateFlexPaymentSummaryResponse
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return &decoded
}

func approvedByBookingID(approved []ApprovedReservation) map[string]ApprovedReservation {
	byID := make(map[string]ApprovedReservation, len(approved))
	for _, a := range approved {
		byID[a.BrokerReservationID] = a
	}
	return byID
}

func rejectedByBookingID(rejected []RejectedReservation) map[string]RejectedReservation {
	byID := make(map[string]RejectedReservation, len(rejected))
	for _, r := range rejected {
		byID[r.BrokerReservationID] = r
	}
	return byID
}

// --- Tests ---

func TestValidateFlexPaymentSummary(t *testing.T) {
	const userID int64 = 90002
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	// Cost owed to Flex is purchase_price + broker_erp_price, before markup and discount.
	const purchasePrice, brokerErpPrice = 200.0, 30.0
	const owed = purchasePrice + brokerErpPrice

	seedFlex := func(bookingID, currencyCode string) int64 {
		return seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
			p.BrokerReservationID = bookingID
			p.Broker = "flex"
			p.CurrencyCode = currencyCode
			p.PurchasePrice = purchasePrice
			p.BrokerErpPrice = brokerErpPrice
		})
	}

	matchingID := seedFlex(bookingMatch, "EUR")
	wrongPriceID := seedFlex(bookingWrongPrice, "EUR")
	wrongCurrencyID := seedFlex(bookingWrongCurrency, "EUR")
	usdID := seedFlex(bookingUSD, "USD")

	alreadyPaidID := seedFlex(bookingAlreadyPaid, "EUR")
	expenseID := "EXPENSE-FLEXPAY"
	if err := s.query.MarkReservationsPaidToSupplier(ctx, db.MarkReservationsPaidToSupplierParams{
		Ids:               []int64{alreadyPaidID},
		SupplierExpenseID: &expenseID,
	}); err != nil {
		t.Fatalf("failed to mark reservation as paid to supplier: %v", err)
	}

	file := buildFlexSummary(t, map[string][]summaryLine{
		"EUR": {
			{bookingID: bookingMatch, balance: owed},
			{bookingID: bookingWrongPrice, balance: owed + 25},
			{bookingID: bookingAlreadyPaid, balance: owed},
			{bookingID: bookingGhost, balance: owed},
		},
		"USD": {
			{bookingID: bookingUSD, balance: owed},
			// Booked in EUR but billed on the USD sheet — the amounts are not comparable.
			{bookingID: bookingWrongCurrency, balance: owed},
		},
	})

	resp := decodeValidationResponse(t, postFlexSummary(t, s, "flex-summary.xlsx", file))
	approved := approvedByBookingID(resp.Approved)
	rejected := rejectedByBookingID(resp.Rejected)

	t.Run("approves lines matching an outstanding reservation", func(t *testing.T) {
		for bookingID, wantID := range map[string]int64{
			bookingMatch: matchingID,
			bookingUSD:   usdID,
		} {
			got, ok := approved[bookingID]
			if !ok {
				t.Errorf("expected %s to be approved, got rejection %+v", bookingID, rejected[bookingID])
				continue
			}
			if got.ReservationID != wantID {
				t.Errorf("%s approved with reservation id %d, want %d", bookingID, got.ReservationID, wantID)
			}
			if got.Amount != owed {
				t.Errorf("%s approved for %.2f, want %.2f", bookingID, got.Amount, owed)
			}
		}
	})

	t.Run("rejects a line whose balance does not match what we owe", func(t *testing.T) {
		got, ok := rejected[bookingWrongPrice]
		if !ok {
			t.Fatalf("expected %s to be rejected", bookingWrongPrice)
		}
		if got.Reason != supplier_payments.ReasonInvalidPrice {
			t.Errorf("reason is %q, want %q", got.Reason, supplier_payments.ReasonInvalidPrice)
		}
		if got.ReservationID == nil || *got.ReservationID != wrongPriceID {
			t.Errorf("reservation id is %v, want %d", got.ReservationID, wrongPriceID)
		}
		if got.ExpectedAmount == nil || *got.ExpectedAmount != owed {
			t.Errorf("expected amount is %v, want %.2f", got.ExpectedAmount, owed)
		}
		if got.Balance != owed+25 {
			t.Errorf("balance is %.2f, want %.2f", got.Balance, owed+25)
		}
	})

	t.Run("rejects a line billed in the wrong currency", func(t *testing.T) {
		got, ok := rejected[bookingWrongCurrency]
		if !ok {
			t.Fatalf("expected %s to be rejected", bookingWrongCurrency)
		}
		if got.Reason != supplier_payments.ReasonInvalidPrice {
			t.Errorf("reason is %q, want %q", got.Reason, supplier_payments.ReasonInvalidPrice)
		}
		if got.ReservationID == nil || *got.ReservationID != wrongCurrencyID {
			t.Errorf("reservation id is %v, want %d", got.ReservationID, wrongCurrencyID)
		}
		if got.CurrencyCode != "USD" || got.ExpectedCurrencyCode != "EUR" {
			t.Errorf("currencies are billed=%q expected=%q, want billed=USD expected=EUR", got.CurrencyCode, got.ExpectedCurrencyCode)
		}
	})

	t.Run("rejects unknown and already-settled reservations", func(t *testing.T) {
		for _, bookingID := range []string{bookingGhost, bookingAlreadyPaid} {
			got, ok := rejected[bookingID]
			if !ok {
				t.Errorf("expected %s to be rejected", bookingID)
				continue
			}
			if got.Reason != supplier_payments.ReasonNotFound {
				t.Errorf("%s rejected with reason %q, want %q", bookingID, got.Reason, supplier_payments.ReasonNotFound)
			}
			if got.ReservationID != nil {
				t.Errorf("%s should carry no reservation id, got %d", bookingID, *got.ReservationID)
			}
		}
	})

	t.Run("reads balances written as formatted text", func(t *testing.T) {
		formatted := buildFlexSummary(t, map[string][]summaryLine{
			"EUR": {{bookingID: bookingMatch, balance: "230.00"}},
		})
		got := decodeValidationResponse(t, postFlexSummary(t, s, "flex-summary.xlsx", formatted))
		if len(got.Approved) != 1 || got.Approved[0].ReservationID != matchingID {
			t.Errorf("expected the reservation to be approved, got approved=%+v rejected=%+v", got.Approved, got.Rejected)
		}
	})

	t.Run("skips blank rows between lines", func(t *testing.T) {
		withGap := buildFlexSummary(t, map[string][]summaryLine{
			"EUR": {
				{bookingID: bookingMatch, balance: owed},
				{bookingID: "", balance: nil},
				{bookingID: bookingWrongPrice, balance: owed},
			},
		})
		got := decodeValidationResponse(t, postFlexSummary(t, s, "flex-summary.xlsx", withGap))
		if len(got.Approved) != 2 {
			t.Errorf("expected both reservations to be approved across the blank row, got %+v", got.Approved)
		}
	})
}

func TestValidateFlexPaymentSummary_InvalidFile(t *testing.T) {
	s := &Service{query: testQuerier()}

	t.Run("rejects a file that is not a workbook", func(t *testing.T) {
		t.Parallel()
		rec := postFlexSummary(t, s, "notes.txt", []byte("this is not an xlsx"))

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
		assertErrorCode(t, rec, api_errors.CodeInvalidPaymentSummaryFile)
	})

	t.Run("rejects a workbook without the currency sheets", func(t *testing.T) {
		t.Parallel()
		file := buildFlexSummary(t, map[string][]summaryLine{
			"Summary": {{bookingID: bookingGhost, balance: 100.0}},
		})
		rec := postFlexSummary(t, s, "flex-summary.xlsx", file)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
		assertErrorCode(t, rec, api_errors.CodeInvalidPaymentSummaryFile)
	})

	t.Run("rejects a workbook whose columns have moved", func(t *testing.T) {
		t.Parallel()
		file := buildFlexSummary(t, map[string][]summaryLine{
			"EUR": {{bookingID: bookingGhost, balance: 100.0}},
		})
		shifted := shiftBalanceHeader(t, file)

		rec := postFlexSummary(t, s, "flex-summary.xlsx", shifted)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
		assertErrorCode(t, rec, api_errors.CodeInvalidPaymentSummaryFile)
	})

	t.Run("rejects a request that is not multipart", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/supplier-payments/flex/validate", bytes.NewReader(nil))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		s.ValidateFlexPaymentSummary(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rec.Code)
		}
	})
}

// shiftBalanceHeader renames the Balance header so it no longer sits in column L, standing in for
// Flex changing the statement layout under us.
func shiftBalanceHeader(t *testing.T, content []byte) []byte {
	t.Helper()

	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		t.Fatalf("failed to open workbook: %v", err)
	}
	defer f.Close()

	if err := f.SetCellValue("EUR", "L14", "Something Else"); err != nil {
		t.Fatalf("failed to move the balance header: %v", err)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatalf("failed to write workbook to buffer: %v", err)
	}
	return buf.Bytes()
}

func assertErrorCode(t *testing.T, rec *httptest.ResponseRecorder, wantCode string) {
	t.Helper()

	var body struct {
		Details api_errors.ErrorDetails `json:"details"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if body.Details.Code != wantCode {
		t.Errorf("error code is %q, want %q", body.Details.Code, wantCode)
	}
}
