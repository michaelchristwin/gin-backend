package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBDSN string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}
	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBDSN: getEnv("DBDSN", "file:app.db?_foreign_keys=on"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
