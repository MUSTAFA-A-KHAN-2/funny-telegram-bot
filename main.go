package main

import (
    "log"
    "telegram-bot/controller"
)

func main() {
    botToken := "YOUR TOKEN HERE"
    err := controller.StartBot(botToken)
    if err != nil {
        log.Panic(err)
    }
}
