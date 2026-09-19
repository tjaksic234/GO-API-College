package db

import (
	"College/pkg"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	pkg.LoadEnv()
	dsn := os.Getenv("DB_CONNECTION_DEV")

	// Retry with backoff instead of a single attempt: the official MySQL
	// image restarts itself once internally on first-ever init of a fresh
	// volume (temp server -> shutdown -> real server), during which
	// connections are refused for a brief window even after Docker's own
	// healthcheck has started reporting healthy. A single attempt can lose
	// this race and crash the app instantly - observed in practice under
	// concurrent multi-stack cloud runs, while Spring/Hikari on the other
	// side retries this transparently. Without this, that timing fluke
	// would unfairly fail Go's measurements for a reason having nothing to
	// do with actual startup/runtime performance.
	var err error
	for i := 0; i < 30; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("DB connection attempt %d/30 failed: %v", i+1, err)
		time.Sleep(time.Second)
	}

	if err != nil {
		log.Fatalf("Error connecting to database after retries: %v", err)
	}
	log.Printf("GORM connected to database")
}
