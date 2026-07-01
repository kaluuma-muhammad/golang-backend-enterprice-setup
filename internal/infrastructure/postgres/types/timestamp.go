package types

import "time"

func ToPGTimestamp(t *time.Time) time.Time {

	if t == nil {
		return time.Time{}
	}

	return *t
}

func FromPGTimestamp(t time.Time) *time.Time {

	if t.IsZero() {
		return nil
	}

	return &t
}
