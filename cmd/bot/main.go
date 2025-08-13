package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	publicURL := os.Getenv("PUBLIC_URL") // example: https://mybot.up.railway.app

	if botToken == "" || publicURL == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN and PUBLIC_URL must be set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Delete old webhook (important when switching from polling)
	_, _ = bot.Request(tgbotapi.DeleteWebhookConfig{})

	// Set new webhook
	webhook := tgbotapi.NewWebhook(publicURL + "/webhook")
	if _, err := bot.Request(webhook); err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	// Start HTTP server to receive webhook updates
	updates := bot.ListenForWebhook("/webhook")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Printf("🤖 Bot started in webhook mode: %s", publicURL)

	// URL pattern
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		if urlRegex.MatchString(text) {
			link := urlRegex.FindString(text)
			log.Printf("🔗 Found link: %s", link)

			// Send "Yuklanmoqda..." message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Yuklanmoqda...")
			sentMsg, _ := bot.Send(msg)

			// Download video using yt-dlp
			filePath := "/tmp/video.mp4"
			cmd := exec.Command("yt-dlp", "-o", filePath, link)
			err := cmd.Run()
			if err != nil {
				log.Printf("❌ Extraction error: %v", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Failed to fetch video"))
				continue
			}

			// Send the video
			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
			video.Caption = "📹 Video yuklandi"
			bot.Send(video)

			// Delete "Yuklanmoqda..." message
			del := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sentMsg.MessageID)
			bot.Request(del)
		}
	}
}
