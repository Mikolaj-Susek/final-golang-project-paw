package config

import "os"

var DB_URL string

func LoadConfig() {
	DB_URL = os.Getenv("DATABASE_URL")
}
