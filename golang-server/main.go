package main

import (
	"github.com/example/golang-postgres-crud/config"
	"github.com/example/golang-postgres-crud/db"
	"github.com/example/golang-postgres-crud/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.LoadConfig()
	db.ConnectDatabase()
	r := routes.SetupRouter()
	r.Run()
}
