package broker

import (
	"errors"
	"fmt"
	"strings"

	"encore.dev/rlog"
	"github.com/xuri/excelize/v2"
)

// Avance sends its station list as an xlsx workbook with a single "Station" sheet: one header row
// followed by one row per station.
const (
	avanceStationsSheet    = "Station"
	avanceStationsFirstRow = 2

	// Zero-based indices of the columns we import; the sheet carries more (charge, working hours,
	// phone, coordinates) that broker.Location has no place for.
	avanceStationCodeCol        = 0 // A
	avanceStationDescriptionCol = 1 // B
	avanceStationTypeCol        = 3 // D
	avanceStationCityCol        = 5 // F
	avanceStationCountryCol     = 7 // H

	avanceStationCodeHeader = "Code"
	avanceStationTypeHeader = "Station type"

	// avanceStationIataHeader names the column our Apps Script adds to Avance's sheet. Avance's own
	// sheet has no IATA codes, and neither the station code nor its type reliably gives one: Athens
	// Airport is 022, and Santorini Airport is typed "Office".
	avanceStationIataHeader = "IATA"
)

var (
	ErrAvanceReaderNotInitialized = errors.New("Avance reader is not initialized")
	ErrAvanceInvalidWorkbook      = errors.New("invalid workbook format for Avance locations")
)

// avanceAirportStationTypes are the types Avance uses for a desk at the airport itself, which the
// results page badges as an airport pickup. Shuttle and meet-and-greet pickups at an airport are a
// different arrangement and Avance classifies them separately.
var avanceAirportStationTypes = map[string]struct{}{
	"Airport":   {},
	"Airport 1": {},
}

var countryNamesByISO = map[string]string{
	"GR": "Greece",
}

// GetLocationsPage retrieves a page of locations from the Avance workbook. The cursor is unused
// since all stations are returned in a single page.
func (a *Avance) GetLocationsPage(cursor string) (LocationPage, error) {
	if a.r == nil {
		return LocationPage{}, ErrAvanceReaderNotInitialized
	}

	file, err := excelize.OpenReader(a.r)
	if err != nil {
		return LocationPage{}, fmt.Errorf("avance locations: open workbook: %w", err)
	}
	defer file.Close()

	rows, err := file.GetRows(avanceStationsSheet)
	if err != nil {
		return LocationPage{}, fmt.Errorf("avance locations: read sheet %q: %w", avanceStationsSheet, err)
	}

	iataCol, err := validateAvanceHeader(rows)
	if err != nil {
		return LocationPage{}, err
	}

	locations := make([]Location, 0, len(rows))
	for i, row := range rows[avanceStationsFirstRow-1:] {
		rowNum := i + avanceStationsFirstRow

		code := avanceCell(row, avanceStationCodeCol)
		if code == "" {
			continue
		}

		name := avanceCell(row, avanceStationDescriptionCol)
		countryCode := strings.ToUpper(avanceCell(row, avanceStationCountryCol))
		if name == "" || countryCode == "" {
			// country_code is not nullable and drives markup lookup, so a station missing it is
			// left out until the sheet is corrected rather than guessed at.
			rlog.Warn("skipping Avance station with missing required fields",
				"row", rowNum, "code", code, "name", name, "countryCode", countryCode)
			continue
		}

		locations = append(locations, Location{
			ID:          code,
			Name:        name,
			City:        avanceCell(row, avanceStationCityCol),
			CountryCode: countryCode,
			Country:     avanceCountryName(countryCode),
			Iata:        strings.ToUpper(avanceCell(row, iataCol)),
		})
	}

	if len(locations) == 0 {
		return LocationPage{}, ErrAvanceInvalidWorkbook
	}

	return LocationPage{Locations: locations}, nil
}

// validateAvanceHeader checks the layout before any data is read, so a reordered workbook fails
// loudly instead of importing every station under the wrong name. It returns the IATA column, which
// is found by its header since it is appended to Avance's own columns.
func validateAvanceHeader(rows [][]string) (int, error) {
	if len(rows) < avanceStationsFirstRow {
		return 0, ErrAvanceInvalidWorkbook
	}

	header := rows[0]
	if avanceCell(header, avanceStationCodeCol) != avanceStationCodeHeader ||
		avanceCell(header, avanceStationTypeCol) != avanceStationTypeHeader {
		return 0, ErrAvanceInvalidWorkbook
	}

	for i := range header {
		if strings.EqualFold(avanceCell(header, i), avanceStationIataHeader) {
			return i, nil
		}
	}

	// Without it every airport would import as a new location beside the existing one.
	return 0, fmt.Errorf("%w: no %q column, run the stations Apps Script first", ErrAvanceInvalidWorkbook, avanceStationIataHeader)
}

// avanceCountryName maps the sheet's ISO code to the display name held in locations.country. The
// upserts overwrite country on conflict, so returning the bare code would degrade the name an
// existing canonical location already carries.
func avanceCountryName(countryCode string) string {
	if name, ok := countryNamesByISO[countryCode]; ok {
		return name
	}

	return countryCode
}

// avanceCell tolerates short rows: excelize omits trailing empty cells.
func avanceCell(row []string, col int) string {
	if col >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[col])
}
