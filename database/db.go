// database/database.go
package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mhandyalf/go-passmanager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_POSTGRES")
	fmt.Printf(">>> Connecting to database...\n")

	// Konfigurasi GORM yang aman
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:                              false, // Disable prepared statements untuk menghindari konflik
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL database: %v", err)
	}

	// Set connection pool
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	DB = db
	log.Println("Successfully connected to PostgreSQL database!")

	// Migrasi dengan error handling yang lebih baik
	migrateDatabase()
}

func migrateDatabase() {
	log.Println("Migrating database...")

	err := DB.AutoMigrate(&models.User{}, &models.Password{})
	if err != nil {
		// Jika error tapi hanya karena tabel sudah ada, abaikan
		if strings.Contains(err.Error(), "already exists") {
			log.Println("Tables already exist, migration skipped")
			return
		}
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed!")
}
