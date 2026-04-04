package configs

import (
	"fmt"
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
	AuthConfig
	DbConfig
}

func NewConfig() *Config {
	_ = godotenv.Load()
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password,
		name, port)
	return &Config{
		AuthConfig: AuthConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
		DbConfig: DbConfig{
			DSN: dsn,
		},
	}
}
