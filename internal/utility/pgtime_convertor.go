package utility

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ToPGTimestamp — конвертация time.Time в pgtype.Timestamp.
func ToPGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// FromPGTimestamp — безопасная конвертация pgtype.Timestamp в time.Time.
func FromPGTimestamp(ts pgtype.Timestamp) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}
