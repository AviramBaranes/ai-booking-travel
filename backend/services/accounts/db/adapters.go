package db

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// NumericParam converts a *float64 to a pgtype.Numeric, handling nil values appropriately.
func NumericParam(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	var n pgtype.Numeric
	n.Scan(fmt.Sprintf("%f", *f))
	return n
}

// FloatFromNumeric converts a pgtype.Numeric back to a *float64, returning nil if not valid.
func FloatFromNumeric(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, _ := n.Float64Value()
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// TextParam converts a *string to a pgtype.Text, handling nil values appropriately.
func TextParam(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// StringFromTextParam converts a pgtype.Text back to a *string, returning nil if the Text is not valid.
func StringFromTextParam(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// StringFromTimeParam converts a pgtype.Timestamptz back to a string, returning an empty string if the Timestamptz is not valid.
func StringFromTimeParam(t pgtype.Timestamptz) string {
	if !t.Valid {
		return ""
	}
	return t.Time.String()
}
