package main

import (
	"quree/internal/bot"
	"quree/internal/bot/adminbot"
	"quree/internal/bot/userbot"
	"sync"
)

var wg sync.WaitGroup
var userBot = bot.Init(userbot.BotConfig)
var adminBot = bot.Init(adminbot.BotConfig)

func main() {
	wg.Add(2)

	go userBot.Start()

	go adminBot.Start()

	wg.Wait()
}
