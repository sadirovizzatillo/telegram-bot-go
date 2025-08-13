package bot

import (
	"log"
	"regexp"

	"telegram-video-bot/internal/extractor"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var linkRegex = regexp.MustCompile(`https?://[^\s]+`)

func Run(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		if linkRegex.MatchString(text) {
			link := linkRegex.FindString(text)
			log.Println("🔗 Found link:", link)

			// Send "Yuklanmoqda..." immediately
			loadingMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Yuklanmoqda...")
			sentMsg, _ := bot.Send(loadingMsg)

			videoURL, err := extractor.GetVideoURL(link)
			if err != nil {
				log.Println("❌ Extraction error:", err)
				edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, sentMsg.MessageID, "❌ Video yuklab bo‘lmadi.")
				bot.Send(edit)
				continue
			}

			// Replace "Yuklanmoqda..." with the video
			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileURL(videoURL))
			bot.Send(video)

			// Optionally delete the loading message instead of editing
			bot.Request(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sentMsg.MessageID))
		}
	}
}
