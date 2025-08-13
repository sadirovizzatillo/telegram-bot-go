package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = false

	// ✅ Correct handling of NewWebhook return value
	webhook, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatalf("Failed to create webhook: %v", err)
	}

	// Set webhook
	if _, err := bot.Request(webhook); err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatalf("Failed to get webhook info: %v", err)
	}
	if info.URL != webhookURL {
		log.Fatalf("Webhook not set correctly: %v", info)
	}

	log.Printf("Bot started in webhook mode at %s", webhookURL)

	// Regex for Instagram links only
	instaRegex := regexp.MustCompile(`https?://(www\.)?instagram\.com/[^\s]+`)

	updates := bot.ListenForWebhook("/")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		if instaRegex.MatchString(text) {
			link := instaRegex.FindString(text)
			log.Printf("📌 Instagram link found: %s", link)

			// Send "Yuklanmoqda..." message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Yuklanmoqda...")
			sent, _ := bot.Send(msg)

			// Download the video
			videoPath, err := downloadInstagramVideo(link)
			if err != nil {
				log.Printf("❌ Download failed: %v", err)
				edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, sent.MessageID, "❌ Video yuklab bo‘lmadi")
				bot.Send(edit)
				continue
			}
			defer os.Remove(videoPath)

			// Send video
			videoFile := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(videoPath))
			if _, err := bot.Send(videoFile); err != nil {
				log.Printf("❌ Failed to send video: %v", err)
			} else {
				// Replace the "Yuklanmoqda..." with success message
				edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, sent.MessageID, "✅ Video yuklandi")
				bot.Send(edit)
			}
		}
	}
}

// downloadInstagramVideo downloads an Instagram video using yt-dlp
func downloadInstagramVideo(url string) (string, error) {
	output := "video.mp4"
	cmd := exec.Command("yt-dlp", "-o", output, url)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	return output, nil
}
