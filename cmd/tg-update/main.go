package main

import (
	"quree/internal/bot/adminbot"
	"quree/internal/bot/userbot"
	"sync"
)

var wg sync.WaitGroup
var userBot = userbot.Bot
var adminBot = adminbot.Bot

func main() {
	wg.Add(3)

	go userBot.Start()

	go adminBot.Start()

	wg.Wait()
}
