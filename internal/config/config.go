package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	WebhookURL string
	Port       string
}

func Load() *Config {
	// Load .env only if not in Railway (local dev)
	if _, ok := os.LookupEnv("RAILWAY_ENVIRONMENT"); !ok {
		_ = godotenv.Load()
	}

	cfg := &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		WebhookURL: os.Getenv("WEBHOOK_URL"),
		Port:       getEnv("PORT", "8080"),
	}

	if cfg.BotToken == "" {
		log.Fatal("❌ BOT_TOKEN not set")
	}
	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
