package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex"`
	UserName  string
	Password  string
	CreatedAt time.Time
}

type Password struct {
	gorm.Model
	Title             string `gorm:"not null"`
	Username          string
	EncryptedPassword string `gorm:"not null"` // Password yang sudah dienkripsi
	Tags              string // Pisahkan dengan koma, misalnya "kerja,sosial"
	UserID            uint   // Foreign key untuk menghubungkan ke user
}
