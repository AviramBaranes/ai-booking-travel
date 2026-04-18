package db

import (
	"errors"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// NumericFromFloat64 converts a float64 to a pgtype.Numeric.
func NumericFromFloat64(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.ScanScientific(big.NewFloat(f).Text('f', 6))
	return n
}

// NumericToFloat64 converts a pgtype.Numeric to a float64.
func NumericToFloat64(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

// DateFromString parses a "2006-01-02" string into a pgtype.Date.
func DateFromString(s string) pgtype.Date {
	var d pgtype.Date
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return d
	}
	d.Time = t
	d.Valid = true
	return d
}

// DateToString formats a pgtype.Date as a "2006-01-02" string.
func DateToString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

// TimestamptzToString formats a pgtype.Timestamptz as an RFC3339 string.
func TimestamptzToString(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.Format(time.RFC3339)
}

// CombineDateTime merges a pgtype.Date and a time string (e.g. "15:04:05Z07:00") into a single time.Time.
func CombineDateTime(d pgtype.Date, timeStr string) (time.Time, error) {
	if !d.Valid {
		return time.Time{}, errors.New("invalid date")
	}
	dateStr := d.Time.Format("2006-01-02")
	dateTimeStr := dateStr + "T" + timeStr
	return time.Parse(time.RFC3339, dateTimeStr)
}
