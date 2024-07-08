package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgtypeCurrentTime() pgtype.Timestamptz {
	currentTime := time.Now()

	// Создаем экземпляр pgtype.Timestamptz и задаем ему значение
	timestamptz := pgtype.Timestamptz{}
	timestamptz.Time = currentTime
	timestamptz.Valid = true

	return timestamptz
}
