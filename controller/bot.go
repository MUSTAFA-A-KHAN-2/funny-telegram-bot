package controller

import (
	"fmt"
	"log"
	"sync"

	"github.com/MUSTAFA-A-KHAN/funny-telegram-bot/model"
	"github.com/MUSTAFA-A-KHAN/funny-telegram-bot/view"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var userJokes = struct {
	sync.RWMutex
	data map[int64]model.Joke
}{data: make(map[int64]model.Joke)}

// StartBot initializes and starts the bot
func StartBot(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		fmt.Println("update.Message.Command():::::::::::::::::::::::::::;", update.Message.Command())

		switch update.Message.Command() {
		case "start":
			view.SendMessage(bot, update.Message.Chat.ID, "Hi! Use /joke to get a joke setup and /punch to get the punchline.")
		case "help":
			view.SendMessage(bot, update.Message.Chat.ID, "Use /joke to get a joke setup and /punch to get the punchline.")
		case "joke":
			joke, err := model.GetJoke()
			if err != nil {
				view.SendMessage(bot, update.Message.Chat.ID, "Failed to get a joke.")
				continue
			}
			userJokes.Lock()
			userJokes.data[update.Message.Chat.ID] = joke
			userJokes.Unlock()
			view.SendMessage(bot, update.Message.Chat.ID, joke.Setup)
		case "punch":
			userJokes.RLock()
			joke, exists := userJokes.data[update.Message.Chat.ID]
			userJokes.RUnlock()
			if !exists {
				view.SendMessage(bot, update.Message.Chat.ID, "Please request a joke setup first using /joke.")
				continue
			}
			view.SendMessage(bot, update.Message.Chat.ID, joke.Punchline)
			userJokes.Lock()
			delete(userJokes.data, update.Message.Chat.ID)
			userJokes.Unlock()
		default:
			view.SendMessage(bot, update.Message.Chat.ID, "uh-OH")
		}
	}

	return nil
}
