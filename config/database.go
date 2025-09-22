package config

import (
	"College/internal/models"
	"College/pkg"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	pkg.LoadEnv()
	var err error
	dsn := os.Getenv("DB_CONNECTION_DEV")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	} else {
		log.Printf("GORM connected to database")
	}

	DB.AutoMigrate(&models.User{})
	log.Println("Database connected and migrated")
}
