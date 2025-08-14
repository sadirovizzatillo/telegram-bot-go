package main

import (
	"log"
	"net/http"
	"os"

	"github.com/izzatillo/telegram-video-bot/config"
	"github.com/izzatillo/telegram-video-bot/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init bot
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("❌ Failed to init bot: %v", err)
	}

	api.Debug = true
	log.Printf("✅ Authorized on account %s", api.Self.UserName)

	// If webhook URL is provided → run webhook mode (Railway)
	if cfg.WebhookURL != "" {
		wh, err := tgbotapi.NewWebhook(cfg.WebhookURL)
		if err != nil {
			log.Fatalf("❌ Failed to create webhook config: %v", err)
		}

		_, err = api.Request(wh)
		if err != nil {
			log.Fatalf("❌ Failed to set webhook: %v", err)
		}

		info, err := api.GetWebhookInfo()
		if err != nil {
			log.Fatalf("❌ Failed to get webhook info: %v", err)
		}
		if info.URL != cfg.WebhookURL {
			log.Fatalf("❌ Webhook URL mismatch. Expected %s, got %s", cfg.WebhookURL, info.URL)
		}

		log.Printf("🚀 Bot started in WEBHOOK mode at %s", cfg.WebhookURL)

		updates := api.ListenForWebhook("/")
		go http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil)

		for update := range updates {
			bot.HandleUpdate(api, update)
		}
	} else {
		// Otherwise → run polling mode (local dev)
		log.Println("🚀 Bot started in POLLING mode")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := api.GetUpdatesChan(u)

		for update := range updates {
			bot.HandleUpdate(api, update)
		}
	}
}
