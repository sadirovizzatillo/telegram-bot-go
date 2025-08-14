package config

import (
	"log"
	"os"
)

type Config struct {
	BotToken   string
	WebhookURL string
}

func Load() *Config {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("❌ BOT_TOKEN not set")
	}

	return &Config{
		BotToken:   token,
		WebhookURL: os.Getenv("WEBHOOK_URL"), // optional
	}
}
