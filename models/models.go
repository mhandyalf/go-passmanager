package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Password struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	Title             string `gorm:"type:varchar(255);not null" json:"title"`
	Username          string `gorm:"type:varchar(255);not null" json:"username"`
	EncryptedPassword string `gorm:"type:text;not null" json:"-"` // Hide dari JSON
	URL               string `gorm:"type:varchar(500)" json:"url"`
	Notes             string `gorm:"type:text" json:"notes"`
	Tags              string `gorm:"type:varchar(500)" json:"tags"` // JSON array lebih baik
	Category          string `gorm:"type:varchar(100)" json:"category"`
	IsFavorite        bool   `gorm:"default:false" json:"is_favorite"`

	// Security fields
	LastAccessed *time.Time `json:"last_accessed"`
	AccessCount  int        `gorm:"default:0" json:"access_count"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign Key
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`

	// Virtual field (ignored by GORM, hanya untuk response JSON)
	DecryptedPassword string `gorm:"-" json:"decrypted_password,omitempty"`
}

type UpdatePasswordInput struct {
	Title    string `json:"title"`
	Username string `json:"username"`
	Password string `json:"password"` // plain password dari FE
	Tags     string `json:"tags"`
}
