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
	var dsn string
	profile := os.Getenv("GO_PROFILES_ACTIVE")

	switch profile {
	case "dev":
		dsn = os.Getenv("DB_CONNECTION_DEV")
	case "docker":
		dsn = os.Getenv("DB_CONNECTION_DOCKER")
	default:
		log.Fatalf("Unknown profile %s", profile)
	}
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	} else {
		log.Printf("GORM connected to database")
	}

	DB.AutoMigrate(&models.User{})
	log.Println("Database connected and migrated")
}
