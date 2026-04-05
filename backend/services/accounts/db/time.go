package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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

// TimePtrFromDB converts a pgtype.Timestamptz to a *time.Time, returning nil if not valid.
func TimePtrFromDB(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
