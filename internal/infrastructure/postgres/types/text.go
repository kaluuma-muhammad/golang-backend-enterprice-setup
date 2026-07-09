package types

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToPGText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

func ToNullableUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{}
	}

	return pgtype.UUID{
		Bytes: *u,
		Valid: true,
	}
}

func ToNullablePGText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

func FromNullableUUID(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}

	res := uuid.UUID(u.Bytes)
	return &res
}

func FromPGText(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}

	return t.String
}

func FromNullablePGText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}

	value := t.String

	return &value
}
