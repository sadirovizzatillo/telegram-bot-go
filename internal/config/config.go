package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, reading environment variables")
	}

	cfg := Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
	}

	if cfg.TelegramToken == "" {
		log.Fatal("❌ TELEGRAM_TOKEN is required")
	}

	return cfg
}
