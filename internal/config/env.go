package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	JWT_SECRET   string
	KafkaBrokers string
}

var Env Config

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	Env = Config{
		DatabaseURL:  getEnv("DB_URL", ""),
		JWT_SECRET:   getEnv("JWT_SECRET", ""),
		KafkaBrokers: getEnv("KAFKA_BROKERS", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
