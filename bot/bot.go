package main

import (
	"greenfield-deploy/bot/config"
	"net/http"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	mu     sync.Mutex
	bot    *tgbotapi.BotAPI
	client *http.Client
	conf   *config.Config
}

func NewBot(conf *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	return &Bot{
		bot:    bot,
		conf:   conf,
		client: &http.Client{},
	}, nil
}

func (b *Bot) Channel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) UserName() string {
	return b.bot.Self.UserName
}

func (b *Bot) Send(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// todo send is thread safe?
	return b.bot.Send(msg)
}

func (b *Bot) CheckUsername(username string) bool {
	_, ok := b.conf.Users[username]
	return ok
}

func getCommand(text string) (string, []string) {
	args := strings.Split(text, " ")
	if len(args) < 1 {
		return text, nil
	}
	return args[0], args
}

func getListCommands() string {
	// todo rename args
	return "Available commands:\n\t/deploy <project> <version> <cluster> <namespace> <env>\n"
}

func getUser(update tgbotapi.Update) User {
	return User{
		Name:   update.Message.Chat.UserName,
		ChatID: update.Message.Chat.ID,
	}
}
