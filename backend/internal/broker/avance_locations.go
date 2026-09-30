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
)

var (
	ErrAvanceReaderNotInitialized = errors.New("Avance reader is not initialized")
	ErrAvanceInvalidWorkbook      = errors.New("invalid workbook format for Avance locations")
)

// avanceAirportStationTypes are the types Avance uses for a desk at the airport itself, the only
// stations that carry an IATA code. Shuttle and meet-and-greet pickups at an airport are a
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

	if err := validateAvanceHeader(rows); err != nil {
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
			Iata:        avanceStationIata(code, avanceCell(row, avanceStationTypeCol)),
		})
	}

	if len(locations) == 0 {
		return LocationPage{}, ErrAvanceInvalidWorkbook
	}

	return LocationPage{Locations: locations}, nil
}

// validateAvanceHeader checks the layout before any data is read, so a reordered workbook fails
// loudly instead of importing every station under the wrong name.
func validateAvanceHeader(rows [][]string) error {
	if len(rows) < avanceStationsFirstRow {
		return ErrAvanceInvalidWorkbook
	}

	header := rows[0]
	if avanceCell(header, avanceStationCodeCol) != avanceStationCodeHeader ||
		avanceCell(header, avanceStationTypeCol) != avanceStationTypeHeader {
		return ErrAvanceInvalidWorkbook
	}

	return nil
}

// avanceStationIata returns the IATA code an airport station serves, taken from the first three
// characters of its code (ATHAP -> ATH), or an empty string for any other station type.
func avanceStationIata(code, stationType string) string {
	if _, ok := avanceAirportStationTypes[strings.TrimSpace(stationType)]; !ok {
		return ""
	}

	if len(code) < 3 {
		return ""
	}

	// Some airports have numeric live codes (Athens is 022), whose prefix is not an IATA code.
	prefix := strings.ToUpper(code[:3])
	for _, r := range prefix {
		if r < 'A' || r > 'Z' {
			return ""
		}
	}

	return prefix
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
