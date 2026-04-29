package billing

import (
	"fmt"

	"encore.dev/config"
	"encore.dev/rlog"
	"github.com/xuri/excelize/v2"
)

type billingConfig struct {
	MonthlyReport monthlyReportConfig
}

type monthlyReportConfig struct {
	Headers monthlyReportHeadersConfig
	Styles  monthlyReportStylesConfig
}

type monthlyReportHeadersConfig struct {
	OfficeName           config.String
	AgentName            config.String
	DriverName           config.String
	ReservationCreatedAt config.String
	ReservationID        config.String
	VoucherDate          config.String
	VoucherNumber        config.String
	AgentVoucherNumber   config.String
	PickupDate           config.String
	ReturnDate           config.String
	CountryCode          config.String
	RentalDays           config.String
	Currency             config.String
	NetPrice             config.String
	FullCoverage         config.String
	TotalNetPrice        config.String
}

type monthlyReportStylesConfig struct {
	HeaderBackgroundColor    config.String
	RefundRowBackgroundColor config.String
	TotalRowBackgroundColor  config.String
	BorderColor              config.String
}

var cfg = config.Load[*billingConfig]()

// excelize doesn't support an "outline" border type; borders must be declared
// per side. Style 1 is a thin solid line.
func cellBorders() []excelize.Border {
	color := cfg.MonthlyReport.Styles.BorderColor()
	sides := []string{"left", "right", "top", "bottom"}
	borders := make([]excelize.Border, len(sides))
	for i, s := range sides {
		borders[i] = excelize.Border{Type: s, Color: color, Style: 1}
	}
	return borders
}

var centerAlignment = &excelize.Alignment{Horizontal: "center", Vertical: "center"}

func ptr[T any](v T) *T { return &v }

// Column indexes for the monthly report sheet. Keep in sync with monthlyReportHeaders.
const (
	colOfficeName = iota
	colAgentName
	colDriverName
	colReservationCreatedAt
	colReservationID
	colVoucherDate
	colVoucherNumber
	colAgentVoucherNumber
	colPickupDate
	colReturnDate
	colCountryCode
	colRentalDays
	colCurrency
	colNetPrice
	colFullCoverage
	colTotalNetPrice

	monthlyReportColCount
)

// monthlyReportHeaders builds the localized column headers for the monthly
// report sheet from the loaded service config, indexed by the col* constants.
func monthlyReportHeaders() [monthlyReportColCount]string {
	h := cfg.MonthlyReport.Headers
	return [monthlyReportColCount]string{
		colOfficeName:           h.OfficeName(),
		colAgentName:            h.AgentName(),
		colDriverName:           h.DriverName(),
		colReservationCreatedAt: h.ReservationCreatedAt(),
		colReservationID:        h.ReservationID(),
		colVoucherDate:          h.VoucherDate(),
		colVoucherNumber:        h.VoucherNumber(),
		colAgentVoucherNumber:   h.AgentVoucherNumber(),
		colPickupDate:           h.PickupDate(),
		colReturnDate:           h.ReturnDate(),
		colCountryCode:          h.CountryCode(),
		colRentalDays:           h.RentalDays(),
		colCurrency:             h.Currency(),
		colNetPrice:             h.NetPrice(),
		colFullCoverage:         h.FullCoverage(),
		colTotalNetPrice:        h.TotalNetPrice(),
	}
}

func generateExcelReport(report Report) ([]byte, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			rlog.Error("failed to close excel file", "error", err)
		}
	}()

	sheetName := report.OrganizationName
	if !report.IsOrganic && len(report.TransactionGroups) > 0 && len(report.TransactionGroups[0].Reservations) > 0 {
		officeName := report.TransactionGroups[0].Reservations[0].OfficeName
		sheetName += " " + officeName
	}

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create new sheet in excel file %w", err)
	}

	// Remove the default sheet so the report has only one sheet.
	if defaultSheet := "Sheet1"; defaultSheet != sheetName {
		if err := f.DeleteSheet(defaultSheet); err != nil {
			return nil, fmt.Errorf("failed to delete default sheet %w", err)
		}
	}

	if err := writeMonthlyReportSheet(f, sheetName, report); err != nil {
		return nil, err
	}

	f.SetActiveSheet(index)

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write excel file to buffer %w", err)
	}

	return buffer.Bytes(), nil
}

func writeMonthlyReportSheet(f *excelize.File, sheetName string, report Report) error {
	rtl := true
	if err := f.SetSheetView(sheetName, 0, &excelize.ViewOptions{RightToLeft: &rtl}); err != nil {
		return fmt.Errorf("failed to set sheet view %w", err)
	}

	styles := cfg.MonthlyReport.Styles
	borders := cellBorders()

	defaultStyle, err := f.NewStyle(&excelize.Style{
		Border:       borders,
		Alignment:    centerAlignment,
		CustomNumFmt: ptr("0.##"),
	})
	if err != nil {
		return fmt.Errorf("failed to create default style %w", err)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Border:    borders,
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{styles.HeaderBackgroundColor()}, Pattern: 1},
		Alignment: centerAlignment,
	})
	if err != nil {
		return fmt.Errorf("failed to create header style %w", err)
	}

	refundStyle, err := f.NewStyle(&excelize.Style{
		Border:       borders,
		Fill:         excelize.Fill{Type: "pattern", Color: []string{styles.RefundRowBackgroundColor()}, Pattern: 1},
		Alignment:    centerAlignment,
		CustomNumFmt: ptr("0.##"),
	})
	if err != nil {
		return fmt.Errorf("failed to create refund style %w", err)
	}

	totalStyle, err := f.NewStyle(&excelize.Style{
		Border:       borders,
		Font:         &excelize.Font{Bold: true},
		Fill:         excelize.Fill{Type: "pattern", Color: []string{styles.TotalRowBackgroundColor()}, Pattern: 1},
		Alignment:    centerAlignment,
		CustomNumFmt: ptr("0.##"),
	})
	if err != nil {
		return fmt.Errorf("failed to create total style %w", err)
	}

	lastCol, err := excelize.ColumnNumberToName(monthlyReportColCount)
	if err != nil {
		return fmt.Errorf("failed to resolve last column name %w", err)
	}

	headers := monthlyReportHeaders()
	headerRow := make([]any, monthlyReportColCount)
	for i, h := range headers {
		headerRow[i] = h
	}
	if err := f.SetSheetRow(sheetName, "A1", &headerRow); err != nil {
		return fmt.Errorf("failed to write header row %w", err)
	}
	if err := f.SetCellStyle(sheetName, "A1", lastCol+"1", headerStyle); err != nil {
		return fmt.Errorf("failed to apply header style %w", err)
	}

	// Freeze the header row so it stays visible while scrolling.
	if err := f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
		Selection:   []excelize.Selection{{SQRef: "A2", ActiveCell: "A2", Pane: "bottomLeft"}},
	}); err != nil {
		return fmt.Errorf("failed to freeze header row %w", err)
	}

	rowNum := 2
	for i, tg := range report.TransactionGroups {
		if i > 0 {
			// Blank separator row between currency groups.
			rowNum++
		}

		for _, r := range tg.Reservations {
			row := reservationToReportRow(r)
			startCell, err := excelize.CoordinatesToCellName(1, rowNum)
			if err != nil {
				return fmt.Errorf("failed to resolve row start cell %w", err)
			}
			if err := f.SetSheetRow(sheetName, startCell, &row); err != nil {
				return fmt.Errorf("failed to write reservation row %w", err)
			}
			endCell := lastCol + fmt.Sprintf("%d", rowNum)
			rowStyle := defaultStyle
			if r.TotalPrice < 0 {
				rowStyle = refundStyle
			}
			if err := f.SetCellStyle(sheetName, startCell, endCell, rowStyle); err != nil {
				return fmt.Errorf("failed to apply reservation row style %w", err)
			}
			rowNum++
		}

		totalRowLength := 2
		totalRow := make([]any, totalRowLength)
		totalRow[0] = tg.Currency
		totalRow[1] = tg.TotalAmount
		startCell, err := excelize.CoordinatesToCellName(monthlyReportColCount-totalRowLength+1, rowNum)
		if err != nil {
			return fmt.Errorf("failed to resolve total row start cell %w", err)
		}
		if err := f.SetSheetRow(sheetName, startCell, &totalRow); err != nil {
			return fmt.Errorf("failed to write total row %w", err)
		}
		endCell := lastCol + fmt.Sprintf("%d", rowNum)
		if err := f.SetCellStyle(sheetName, startCell, endCell, totalStyle); err != nil {
			return fmt.Errorf("failed to apply total style %w", err)
		}
		rowNum++
	}

	return nil
}

func reservationToReportRow(r Reservations) []any {
	row := make([]any, monthlyReportColCount)
	row[colOfficeName] = r.OfficeName
	row[colAgentName] = r.AgentName
	row[colDriverName] = r.DriverName
	row[colReservationCreatedAt] = r.ReservationCreationDate
	row[colReservationID] = r.ReservationBrokerID
	row[colVoucherDate] = r.VoucherDate
	row[colVoucherNumber] = r.ReservationID
	row[colAgentVoucherNumber] = r.VoucherNumber
	row[colPickupDate] = r.PickupDate
	row[colReturnDate] = r.ReturnDate
	row[colCountryCode] = r.CountryCode
	row[colRentalDays] = r.RentalDays
	row[colCurrency] = r.Currency
	row[colNetPrice] = r.CarPrice
	row[colFullCoverage] = r.ERPPrice
	row[colTotalNetPrice] = r.TotalPrice
	return row
}
