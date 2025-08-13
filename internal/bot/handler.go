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

			videoURL, err := extractor.GetVideoURL(link)
			if err != nil {
				log.Println("❌ Extraction error:", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Failed to fetch video"))
				continue
			}

			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileURL(videoURL))
			_, err = bot.Send(video)
			if err != nil {
				log.Println("⚠️ Video too big, sending link instead")
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, videoURL))
			}
		}
	}
}
