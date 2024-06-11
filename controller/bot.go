package controller

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/MUSTAFA-A-KHAN/funny-telegram-bot/model"
	"github.com/MUSTAFA-A-KHAN/funny-telegram-bot/service"
	"github.com/MUSTAFA-A-KHAN/funny-telegram-bot/view"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var userJokes = struct {
	sync.RWMutex
	data map[int64]model.Joke
}{data: make(map[int64]model.Joke)}

// StartBot initializes and starts the bot
func StartBot(token string) error {
	go startHTTPServer() //start http server with go routine

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
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				buttons := []tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData("🔍 Setup", "setup"),
				}
				view.SendMessageWithButtons(bot, update.Message.Chat.ID, "Click 'Setup' to get a joke setup.", buttons)
			default:
				userJokes.RLock()
				_, exists := userJokes.data[update.Message.Chat.ID]
				userJokes.RUnlock()
				if exists {
					handleGuess(bot, update.Message)
				} else {
					view.SendMessage(bot, update.Message.Chat.ID, "No joke setup found. Click 'Setup' to get a new joke.")
				}
			}
		}

		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			switch callback.Data {
			case "setup":
				userJokes.RLock()
				_, exists := userJokes.data[callback.Message.Chat.ID]
				userJokes.RUnlock()
				if exists {
					bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "You've already got a joke setup. Guess the punchline or press 'Punchline' to reveal it."))
				} else {
					joke, err := model.GetJoke()
					if err != nil {
						view.SendMessage(bot, callback.Message.Chat.ID, "Failed to get a joke.")
						continue
					}

					ans := joke.Punchline
					fmt.Println("answer::::", ans)
					// bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, ans))/**uncomment if you need it for multiusers */
					userJokes.Lock()
					userJokes.data[callback.Message.Chat.ID] = joke
					userJokes.Unlock()
					buttons := []tgbotapi.InlineKeyboardButton{
						tgbotapi.NewInlineKeyboardButtonData("🎭 Punchline", "punchline"),
					}
					view.SendMessageWithButtons(bot, callback.Message.Chat.ID, joke.Setup, buttons)
					fmt.Println("Punchline::::", joke.Punchline)
				}
			case "punchline":
				userJokes.RLock()
				joke, exists := userJokes.data[callback.Message.Chat.ID]
				userJokes.RUnlock()
				if exists {
					buttons := []tgbotapi.InlineKeyboardButton{
						tgbotapi.NewInlineKeyboardButtonData("🔍 Setup", "setup"),
					}
					view.SendMessageWithButtons(bot, callback.Message.Chat.ID, joke.Punchline, buttons)
					userJokes.Lock()
					delete(userJokes.data, callback.Message.Chat.ID)
					userJokes.Unlock()
				} else {
					bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, "No joke setup found. Click 'Setup' to get a new joke."))
				}
			}
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		}
	}

	return nil
}

/** responds to the user inputs*/
func handleGuess(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	userJokes.RLock()
	joke, exists := userJokes.data[msg.Chat.ID]
	userJokes.RUnlock()
	if exists {
		// if strings.EqualFold(msg.Text, joke.Punchline) {
		if service.NormalizeAndCompare(msg.Text, joke.Punchline) {
			buttons := []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("🔍 Setup", "setup"),
			}
			view.SendMessageWithButtons(bot, msg.Chat.ID, "😁😁 XDXD! You guessed it right!🏆", buttons)
			userJokes.Lock()
			delete(userJokes.data, msg.Chat.ID)
			userJokes.Unlock()
		} else {
			buttons := []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("🎭 Punchline", "punchline"),
			}
			view.SendMessageWithButtons(bot, msg.Chat.ID, "Hushhh. Try again or click 'Punchline' to reveal the punchline.", buttons)
		}
	}
}

// startHTTPServer starts a simple HTTP server for health checks
func startHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bot is running!")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
