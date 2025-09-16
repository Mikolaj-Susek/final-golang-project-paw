package routes

import (
	"github.com/example/golang-postgres-crud/handlers"
	"github.com/example/golang-postgres-crud/middleware"
	"github.com/gin-gonic/gin"

	_ "github.com/example/golang-postgres-crud/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/text-readings", handlers.CreateTextReading)
		api.GET("/text-readings", handlers.GetTextReadings)
		api.GET("/text-readings/:id", handlers.GetTextReading)
		api.PUT("/text-readings/:id", handlers.UpdateTextReading)
		api.DELETE("/text-readings/:id", handlers.DeleteTextReading)

		api.POST("/ocr", handlers.PerformOcr)
	}

	return router
}
