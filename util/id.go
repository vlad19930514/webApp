package util

import (
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgtypeUUID() pgtype.UUID {
	generatedUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Ошибка при создании UUID: %v", err)
	}
	//TODO sqlc генерирует pgtype
	// Приведение сгенерированного UUID к типу pgtype.UUID
	dbUUID := pgtype.UUID{}
	copy(dbUUID.Bytes[:], generatedUUID[:]) // Устанавливаем bytes field
	dbUUID.Valid = true                     // Устанавливаем valid field

	id, _ := dbUUID.UUIDValue()
	return id
}
