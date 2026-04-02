package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

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
