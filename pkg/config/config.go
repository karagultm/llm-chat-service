package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	ApiKey      string
}

// godotenv uyumlu değil bu
func Load() *Config {
	_ = godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	cfg := &Config{
		Env:         getEnv("APP_ENV", "dev"),
		Port:        getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ApiKey:      getEnv("OPENAI_API_KEY", ""),
	}
	if cfg.ApiKey == "" {
		log.Println("Warning: OPENAI_API_KEY is not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}
	return value
}
