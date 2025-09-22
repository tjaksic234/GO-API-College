package main

import (
	"College/config"
	"College/routes"
)

func main() {
	config.InitDB()

	router := routes.SetupRoutes()

	router.Run()
}
