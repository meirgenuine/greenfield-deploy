package main

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/sync/semaphore"
)

func worker(sm *semaphore.Weighted, wg *sync.WaitGroup, bot *Bot, update tgbotapi.Update) {
	defer func() {
		sm.Release(1)
		wg.Done()
		if r := recover(); r != nil {
			log.Println("recovered", r)
		}
	}()

	user := getUser(update)
	if !bot.CheckUsername(user.Name) {
		log.Println("unknown person", user.Name)
		return
	}

	cmd, args := getCommand(update.Message.Text)

	var resp string
	switch cmd {
	case "/start":
		resp = bot.StartHandler()
	case "/deploy":
		resp = bot.DeployHandler(user, args...)
	default:
		resp = getListCommands()
	}

	log.Println("user", user.Name, "response", resp)
	msg := tgbotapi.NewMessage(
		user.ChatID,
		resp,
	)
	msg.ParseMode = tgbotapi.ModeMarkdown
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("send", err)
	}
}
