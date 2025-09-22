package routes

import (
	"College/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	users := v1.Group("/users")

	users.GET("", handlers.GetUsers)
	users.GET(":id", handlers.GetUserByID)
	users.POST("", handlers.CreateUser)
	users.PATCH(":id", handlers.UpdateUser)
	users.DELETE(":id", handlers.DeleteUser)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "User API is running"})
	})

	return router
}
