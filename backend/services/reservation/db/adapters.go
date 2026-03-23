package db

import (
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
