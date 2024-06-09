package view

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendMessage sends a message to the user
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) error {
    msg := tgbotapi.NewMessage(chatID, text)
    _, err := bot.Send(msg)
    return err
}
