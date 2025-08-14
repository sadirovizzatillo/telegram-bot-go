package main

import (
	"log"
	"net/http"

	"telegram-video-bot/internal/bot"
	"telegram-video-bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.Load()

	// Init bot
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("❌ Failed to create bot: %v", err)
	}

	if cfg.WebhookURL == "" {
		// Local polling mode
		log.Println("⚡ Running in POLLING mode")
		bot.Run(api)
	} else {
		// Webhook mode (Railway, production)
		log.Println("🌍 Running in WEBHOOK mode on port", cfg.Port)

		wh, err := tgbotapi.NewWebhook(cfg.WebhookURL)
		if err != nil {
			log.Fatalf("❌ Failed to create webhook config: %v", err)
		}

		_, err = api.Request(wh)
		if err != nil {
			log.Fatalf("❌ Failed to set webhook: %v", err)
		}
		if err != nil {
			log.Fatalf("❌ Failed to set webhook: %v", err)
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			update, err := api.HandleUpdate(r)
			if err != nil {
				log.Println("Webhook error:", err)
				return
			}
			if update.Message != nil {
				// Reuse your existing handler (polling logic)
				go func() {
					// Here, we just simulate "Run" but for one update
					msg := update.Message.Text
					if msg != "" {
						// you can refactor bot.Run to also accept single update if needed
						log.Println("Got message via webhook:", msg)
					}
				}()
			}
		})

		log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
	}
}
