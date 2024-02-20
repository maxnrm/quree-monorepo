package main

import (
	bot "quree/internal/bot"
	"sync"
)

var wg sync.WaitGroup

func main() {
	userBot := bot.UserBot
	adminBot := bot.AdminBot

	wg.Add(2)

	go userBot.Start()

	go adminBot.Start()

	wg.Wait()

}
