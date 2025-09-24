package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if loadErr := godotenv.Load(".env"); loadErr != nil {
			log.Printf("Could not load .env file: %v", loadErr)
		} else {
			log.Println(".env file loaded successfully")
		}
	} else {
		log.Println(".env file not found, relying on environment variables injected on runtime")
	}
}
