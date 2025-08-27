// database/database.go
package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Buat string koneksi DSN (Data Source Name)
	dsn := os.Getenv("DB_POSTGRES")

	// Buka koneksi ke database PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL database: %v", err)
	}

	// Assign koneksi ke variabel global DB
	DB = db
	log.Println("Successfully connected to PostgreSQL database!")
}
