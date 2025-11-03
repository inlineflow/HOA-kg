package database

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgStringToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}

	return t.String
}

func Text(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

func Now() pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}
}
