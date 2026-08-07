package broker

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Flex sends its payment-required statement as an xlsx workbook with one sheet per currency.
// Each sheet opens with account metadata, a header row, a currency row and a "CARS" section
// label; the booking lines run from flexSummaryFirstRow down to a "Total" row, followed by a
// block of free-text notes.
const (
	flexSummaryHeaderRow = 14
	flexSummaryFirstRow  = 17

	// Column B holds the booking id, K the row label ("Total" on the closing row) and L the
	// balance still owed. Column A is empty.
	flexSummaryBookingIDCol = 1
	flexSummaryLabelCol     = 10
	flexSummaryBalanceCol   = 11

	flexSummaryBookingIDHeader = "Booking ID"
	flexSummaryBalanceHeader   = "Balance"
	flexSummaryTotalLabel      = "Total"
)

// flexSummaryCurrencies are the sheet names Flex uses, one per settlement currency.
var flexSummaryCurrencies = []string{"EUR", "USD"}

// ReadPaymentSummary parses a Flex payment-required statement into one group per currency sheet.
// Sheets absent from the workbook are skipped; it is an error for none of them to be present.
func (f *Flex) ReadPaymentSummary(r io.Reader) ([]PaymentSummaryGroup, error) {
	file, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("flex payment summary: open workbook: %w", err)
	}
	defer file.Close()

	groups := make([]PaymentSummaryGroup, 0, len(flexSummaryCurrencies))
	for _, currencyCode := range flexSummaryCurrencies {
		index, err := file.GetSheetIndex(currencyCode)
		if err != nil {
			return nil, fmt.Errorf("flex payment summary: look up sheet %q: %w", currencyCode, err)
		}
		if index == -1 {
			continue
		}

		lines, err := readFlexSummarySheet(file, currencyCode)
		if err != nil {
			return nil, err
		}

		groups = append(groups, PaymentSummaryGroup{CurrencyCode: currencyCode, Lines: lines})
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("flex payment summary: workbook has none of the expected sheets %v", flexSummaryCurrencies)
	}

	return groups, nil
}

func readFlexSummarySheet(file *excelize.File, sheetName string) ([]PaymentSummaryLine, error) {
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("flex payment summary: read sheet %q: %w", sheetName, err)
	}

	if err := assertFlexSummaryLayout(rows, sheetName); err != nil {
		return nil, err
	}

	lines := make([]PaymentSummaryLine, 0)
	for i := flexSummaryFirstRow - 1; i < len(rows); i++ {
		row := rows[i]

		// The "Total" row closes the list; everything below it is free-text notes that would
		// otherwise be mistaken for booking lines.
		if strings.EqualFold(strings.TrimSpace(cellAt(row, flexSummaryLabelCol)), flexSummaryTotalLabel) {
			break
		}

		bookingID := strings.TrimSpace(cellAt(row, flexSummaryBookingIDCol))
		balanceCell := strings.TrimSpace(cellAt(row, flexSummaryBalanceCol))
		// Blank rows and section labels ("CARS", "MOTORHOMES") carry no balance.
		if bookingID == "" || balanceCell == "" {
			continue
		}

		balance, err := parseAmount(balanceCell)
		if err != nil {
			// Row numbers are 1-based in the workbook, so report i+1 to match what the user sees.
			return nil, fmt.Errorf("flex payment summary: sheet %q row %d: %w", sheetName, i+1, err)
		}

		lines = append(lines, PaymentSummaryLine{BookingID: bookingID, Balance: balance})
	}

	return lines, nil
}

// assertFlexSummaryLayout fails loudly if the columns we read by position are not the ones we
// expect, rather than silently reconciling money against the wrong numbers.
func assertFlexSummaryLayout(rows [][]string, sheetName string) error {
	if len(rows) < flexSummaryHeaderRow {
		return fmt.Errorf("flex payment summary: sheet %q has no header row %d", sheetName, flexSummaryHeaderRow)
	}

	header := rows[flexSummaryHeaderRow-1]
	for col, want := range map[int]string{
		flexSummaryBookingIDCol: flexSummaryBookingIDHeader,
		flexSummaryBalanceCol:   flexSummaryBalanceHeader,
	} {
		got := strings.TrimSpace(cellAt(header, col))
		if !strings.EqualFold(got, want) {
			colName, _ := excelize.ColumnNumberToName(col + 1)
			return fmt.Errorf("flex payment summary: sheet %q column %s is %q, expected %q", sheetName, colName, got, want)
		}
	}

	return nil
}

func cellAt(row []string, col int) string {
	// GetRows trims trailing empty cells, so short rows are expected.
	if col >= len(row) {
		return ""
	}
	return row[col]
}

// parseAmount reads a spreadsheet money cell, tolerating thousands separators, currency
// symbols and parenthesised negatives ("(1,234.56)" meaning -1234.56).
func parseAmount(cell string) (float64, error) {
	s := strings.TrimSpace(cell)
	if s == "" {
		return 0, fmt.Errorf("balance cell is empty")
	}

	negative := strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")")
	if negative {
		s = strings.TrimSuffix(strings.TrimPrefix(s, "("), ")")
	}

	s = strings.NewReplacer(",", "", " ", "", " ", "", "€", "", "$", "").Replace(s)

	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("balance %q is not a number", cell)
	}

	if negative {
		amount = -amount
	}
	return amount, nil
}
