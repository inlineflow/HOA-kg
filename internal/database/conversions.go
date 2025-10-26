package database

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgTimestampToTime(t pgtype.Timestamp) *time.Time {
	if !t.Valid {
		return nil
	}

	return &t.Time
}
