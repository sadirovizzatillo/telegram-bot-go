package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	publicURL := os.Getenv("PUBLIC_URL") // e.g. https://mybot.up.railway.app
	if publicURL == "" {
		log.Fatal("PUBLIC_URL not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Set webhook
	_, err = bot.Request(tgbotapi.NewWebhook(publicURL + "/webhook"))
	if err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	log.Println("Bot started in webhook mode on", publicURL)

	// Handler for Telegram updates
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Println("Error decoding update:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if update.Message != nil && update.Message.Text != "" {
			// Simple reply (replace with your video download logic)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You sent: "+update.Message.Text)
			bot.Send(msg)
		}

		w.WriteHeader(http.StatusOK)
	})

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
