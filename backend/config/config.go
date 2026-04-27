package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	godotenv.Load()
}

func GetDBConnectionString() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "rental"
	}
	if password == "" {
		password = "rental123"
	}
	if dbname == "" {
		dbname = "rental"
	}

	return "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
}

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "development_jwt_secret_key_change_in_production"
	}
	return []byte(secret)
}
