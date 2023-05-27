package main

import (
	"log"
	"strings"

	"greenfield-deploy/bot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	conf := config.Load("config/config.yaml")
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	handler := NewHandler(conf)

	log.Printf("authorized on account %s", bot.Self.UserName)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		user := getUser(update)
		if !checkUsername(user.Name, conf) {
			log.Println("unknown person", user.Name)
			continue
		}

		cmd, args := getCommand(update.Message.Text)

		var resp string
		switch cmd {
		case "/start":
			resp = handler.Start()
		case "/deploy":
			resp = handler.Deploy(user, args...)
		default:
			resp = getListCommands()
		}

		_, err = bot.Send(
			tgbotapi.NewMessage(
				user.ChatID,
				resp,
			))
		if err != nil {
			log.Println("send", err)
		}
	}
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

func checkUsername(username string, conf *config.Config) bool {
	_, ok := conf.Users[username]
	return ok
}

func getUser(update tgbotapi.Update) User {
	return User{
		Name:   update.Message.Chat.UserName,
		ChatID: update.Message.Chat.ID,
	}
}
