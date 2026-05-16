package dbadapters

import (
	"errors"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Numeric = pgtype.Numeric

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

// UuidToString converts a pgtype.UUID to its string representation, returning an empty string if the UUID is null.
func UuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return u.String()
}

// StringToUuid parses a string into a pgtype.UUID. Returns an invalid (null) value if the string is not a valid UUID.
func StringToUuid(s string) pgtype.UUID {
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return pgtype.UUID{}
	}
	return u
}

// TimestamptzToString formats a pgtype.Timestamptz as an RFC3339 string, returning an empty string if the value is null.
func TimestamptzToString(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.Format(time.RFC3339)
}

// DateFromString parses a "2006-01-02" string into a pgtype.Date.
// The input is expected to be pre-validated; on parse error it returns a zero pgtype.Date with Valid=false.
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

// CombineDateTime merges a pgtype.Date and a time string (e.g. "15:04") into a single time.Time.
func CombineDateTime(d pgtype.Date, timeStr string) (time.Time, error) {
	if !d.Valid {
		return time.Time{}, errors.New("invalid date in pgtype.Date")
	}
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	date := d.Time
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC), nil
}

// DBTime converts a time.Time to a pgtype.Timestamptz.
func DBTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// TimeFromDB converts a pgtype.Timestamptz to a time.Time.
func TimeFromDB(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
