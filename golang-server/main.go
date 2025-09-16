package main

import (
	"github.com/example/golang-postgres-crud/config"
	"github.com/example/golang-postgres-crud/db"
	"github.com/example/golang-postgres-crud/routes"
	"github.com/joho/godotenv"
)

// @title Golang Postgres CRUD API
// @version 1.0
// @description This is a sample server for a Go CRUD application with JWT authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	godotenv.Load()
	config.LoadConfig()
	db.ConnectDatabase()
	router := routes.SetupRouter()
	router.Run()
}
