package types

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func ToPGTimestampPtr(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{}
	}

	return pgtype.Timestamp{
		Time:  *t,
		Valid: true,
	}
}

func FromPGTimestamp(ts pgtype.Timestamp) time.Time {
	if !ts.Valid {
		return time.Time{}
	}

	return ts.Time
}

func FromPGTimestampPtr(ts pgtype.Timestamp) *time.Time {
	if !ts.Valid {
		return nil
	}

	t := ts.Time
	return &t
}
