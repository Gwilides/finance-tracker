package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AuthConfig struct {
	Secret string
}

type DbConfig struct {
	DSN string
}

type Config struct {
	Auth AuthConfig
	Db   DbConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading .env file failure")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password,
		name, port)
	return &Config{
		Auth: AuthConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
		Db: DbConfig{
			DSN: dsn,
		},
	}
}
