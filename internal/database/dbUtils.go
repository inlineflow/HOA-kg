package database

import "github.com/jackc/pgx/v5/pgtype"

func PgStringToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}

	return t.String
}
