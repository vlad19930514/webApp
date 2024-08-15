package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"type:text;unique"`
	Age       uint8
	CreatedAt time.Time
}
