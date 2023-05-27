package notification

import (
	"greenfield-deploy/bot/config"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Messenger interface {
	Notify(User, Notification) error
}

type tgbot struct {
	bot *tgbotapi.BotAPI
}

func NewTgMessanger() Messenger {
	cfg := config.Load("bot/config/config.yaml")
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	// todo config bot
	bot.Debug = false
	return &tgbot{
		bot: bot,
	}
}

func (b *tgbot) Notify(u User, nt Notification) error {
	_, err := b.bot.Send(
		tgbotapi.NewMessage(
			int64(u.ChatID),
			nt.String(),
		))
	return err
}
