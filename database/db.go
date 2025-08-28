// database/database.go
package database

import (
	"fmt"
	"log"
	"os"

	"github.com/mhandyalf/go-passmanager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Buat string koneksi DSN (Data Source Name)
	dsn := os.Getenv("DB_POSTGRES")

	// Debug print dulu sebelum connect
	fmt.Printf(">>> DSN raw from env: %q\n", dsn)

	// Buka koneksi ke database PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL database: %v", err)
	}

	fmt.Printf("%q\n", dsn)
	// Assign koneksi ke variabel global DB
	DB = db
	log.Println("Successfully connected to PostgreSQL database!")

	// Migrasi model ke database
	migrateDatabase()
}

// migrateDatabase akan memigrasi semua model ke database
func migrateDatabase() {
	log.Println("Migrating database...")
	// Gunakan AutoMigrate untuk membuat atau memperbarui tabel
	err := DB.AutoMigrate(&models.User{}, &models.Password{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed!")
}
