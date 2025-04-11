package database

import (
	"fmt"
	"os"

	"github.com/hidiyitis/portal-pegawai/internal/core/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}

	// Ambil nilai dari .env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	println("Connected to database")
	_ = db.AutoMigrate(&domain.Department{}, &domain.User{}, &domain.Agenda{}, &domain.Participant{})

	return db
}
