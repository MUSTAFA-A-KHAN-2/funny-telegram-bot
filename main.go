package main

import (
	"log"

	"github.com/MUSTAFA-A-KHAN/funny-telegram-bot/controller"
)

func main() {
	botToken := "7484235929:AAHiFUWLo2nmqXyMe9gby7yc0SBUb8ZysE4"
	err := controller.StartBot(botToken)
	if err != nil {
		log.Panic("error strting bot: ", err)
	}
}
