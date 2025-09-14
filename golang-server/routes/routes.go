package routes

import (
	"github.com/example/golang-postgres-crud/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)
	r.GET("/users/:id", handlers.GetUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)
	r.POST("/ocr", handlers.PerformOcr)
	return r
}
