package main

import (
	"log"
	"telegram-video-bot/internal/bot"
	"telegram-video-bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.Load()

	tgBot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("🤖 Authorized on account %s", tgBot.Self.UserName)

	bot.Run(tgBot)
}
