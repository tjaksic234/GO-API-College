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
	db.InitDB()
	if err := db.DB.AutoMigrate(&models.User{}, &models.Room{}, &models.Message{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	hub := server.NewHub()
	go hub.Run()
	s := server.NewServer(hub)

	// ReleaseMode + Recovery-only (no gin.Default()'s per-request access
	// logger): keeps stdout I/O overhead comparable to the Spring Boot side
	// under K6 load, so latency/throughput differences reflect the
	// runtime/framework, not incidental logging verbosity.
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	r.GET("/ws", s.HandleConnections)

	// REST
	r.GET("/rooms", s.GetAllRooms)
	r.GET("/rooms/active", s.GetActiveRooms)
	r.GET("/messages/:room", s.GetRoomMessages)
	r.GET("/users", s.GetAllUsers)

	port := 8080
	fmt.Printf("Server running on http://localhost:%d\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
