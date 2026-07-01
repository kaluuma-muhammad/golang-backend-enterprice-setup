package types

import "github.com/jackc/pgx/v5/pgtype"

func ToPGText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
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
