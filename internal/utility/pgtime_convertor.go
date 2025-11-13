package utility

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func FromPGTimestamp(ts pgtype.Timestamp) time.Time {
	if ts.Valid {
		return ts.Time
	}

	return time.Time{}
}
