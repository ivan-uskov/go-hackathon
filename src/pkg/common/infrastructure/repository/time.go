package repository

import (
	"database/sql"
	"time"
)

func TimePointer(nullableTime sql.NullTime) *time.Time {
	if nullableTime.Valid {
		return &nullableTime.Time
	}

	return nil
}
