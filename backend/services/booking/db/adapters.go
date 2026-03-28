package db

import (
	"math/big"

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

// NullBrokerTranslationStatusFromString converts a status string to a NullBrokerTranslationStatus.
// Returns an invalid (null) value if the string is empty.
func NullBrokerTranslationStatusFromString(s string) NullBrokerTranslationStatus {
	if s == "" {
		return NullBrokerTranslationStatus{}
	}
	return NullBrokerTranslationStatus{
		BrokerTranslationStatus: BrokerTranslationStatus(s),
		Valid:                   true,
	}
}
