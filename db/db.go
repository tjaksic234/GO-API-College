package db

import (
	"College/pkg"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	pkg.LoadEnv()
	dsn := os.Getenv("DB_CONNECTION")

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	} else {
		log.Printf("GORM connected to database")
	}
}
