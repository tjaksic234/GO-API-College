package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"College/db"
	"College/models"
	"College/server"
)

func main() {
	// DB
	db.InitDB()
	if err := db.DB.AutoMigrate(&models.User{}, &models.Room{}, &models.Message{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	// Hub + server
	hub := server.NewHub()
	go hub.Run()
	s := server.NewServer(hub)

	// Gin
	r := gin.Default()
	r.Use(cors.Default())

	// Health
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// WebSocket
	r.GET("/ws", s.HandleConnections)

	// REST
	r.GET("/rooms", s.GetAllRooms) // persistent rooms (from DB)
	r.GET("/rooms/active", s.GetActiveRooms)
	r.GET("/messages/:room", s.GetRoomMessages)
	r.GET("/users", s.GetAllUsers)

	// (optional) serve your static client if you drop index.html in project root:
	// r.StaticFile("/", "./index.html")

	port := 8080
	fmt.Printf("Server running on http://localhost:%d\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
