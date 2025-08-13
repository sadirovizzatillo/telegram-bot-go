package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ensureYtDlp() {
	_, err := exec.LookPath("yt-dlp")
	if err != nil {
		log.Println("yt-dlp not found, downloading...")
		cmd := exec.Command("sh", "-c", `
			curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
			-o /tmp/yt-dlp && chmod +x /tmp/yt-dlp && mv /tmp/yt-dlp /usr/local/bin/yt-dlp
		`)
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Fatalf("failed to install yt-dlp: %v\n%s", err, string(out))
		}
	} else {
		log.Println("yt-dlp already installed")
	}
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL not set")
	}

	// Ensure yt-dlp exists
	ensureYtDlp()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Remove old webhook if exists
	bot.Request(tgbotapi.DeleteWebhookConfig{})

	// Set webhook
	_, err = bot.Request(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Fatal(err)
	}

	info, _ := bot.GetWebhookInfo()
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	log.Printf("Bot started in webhook mode at %s", webhookURL)

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":8080", nil)

	instaRegex := regexp.MustCompile(`https?://(www\.)?instagram\.com/[^\s]+`)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		if instaRegex.MatchString(text) {
			link := instaRegex.FindString(text)
			log.Printf("📌 Instagram link found: %s", link)

			// Send "Yuklanmoqda..."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Yuklanmoqda...")
			sent, _ := bot.Send(msg)

			// Download video
			outputFile := "/tmp/video.mp4"
			cmd := exec.Command("yt-dlp", "-o", outputFile, link)
			if out, err := cmd.CombinedOutput(); err != nil {
				log.Printf("❌ Download failed: %v: %s", err, string(out))
				bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, sent.MessageID, "❌ Video yuklab bo‘lmadi"))
				continue
			}

			// Send video
			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(outputFile))
			bot.Send(video)

			// Remove "Yuklanmoqda..."
			bot.Request(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sent.MessageID))
		}
	}
}
