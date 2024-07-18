package util

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgtypeUUID() (pgtype.UUID, error) {
	generatedUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Ошибка при создании UUID: %v", err)          //TODO log.Fatal может быть только в main
		return pgtype.UUID{}, fmt.Errorf("uuid.NewUUID:%w", err) //TODO w посмотреть, так обрабатывать ошибки
	}
	//TODO sqlc генерирует pgtype
	// Приведение сгенерированного UUID к типу pgtype.UUID
	dbUUID := pgtype.UUID{}
	copy(dbUUID.Bytes[:], generatedUUID[:]) // Устанавливаем bytes field
	dbUUID.Valid = true                     // Устанавливаем valid field

	id, err := dbUUID.UUIDValue() //TODO не пропускать ошибки
	if err != nil {
		log.Fatalf("Ошибка при создании UUID: %v", err) //TODO log.Fatal может быть только в main
		//
	}
	return id, err
}
