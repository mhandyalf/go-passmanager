package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID
	Email     string `gorm:"uniqueIndex"`
	UserName  string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
}

type Password struct {
	gorm.Model
	Title             string `gorm:"not null"`
	Username          string
	EncryptedPassword string `gorm:"not null"` // Password yang sudah dienkripsi
	Tags              string // Pisahkan dengan koma, misalnya "kerja,sosial"
	UserID            uuid.UUID
}
