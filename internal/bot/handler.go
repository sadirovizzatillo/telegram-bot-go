package bot

import (
	"log"
	"regexp"

	"telegram-video-bot/internal/extractor"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var linkRegex = regexp.MustCompile(`https?://[^\s]+`)

func HandleUpdate(api *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	text := update.Message.Text
	if linkRegex.MatchString(text) {
		link := linkRegex.FindString(text)
		log.Println("🔗 Found link:", link)

		// Send loading message
		loadingMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Yuklanmoqda...")
		sentMsg, _ := api.Send(loadingMsg)

		videoURL, err := extractor.GetVideoURL(link)
		if err != nil {
			log.Println("❌ Extraction error:", err)
			edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, sentMsg.MessageID, "❌ Video yuklab bo‘lmadi.")
			api.Send(edit)
			return
		}

		// Send video
		video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileURL(videoURL))
		api.Send(video)

		// Delete loading message
		api.Request(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sentMsg.MessageID))
	}
}
