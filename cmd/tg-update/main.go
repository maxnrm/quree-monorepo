package main

import (
	"fmt"
	"quree/internal/bot/adminbot"
	"quree/internal/bot/userbot"
	ws "quree/internal/webserver"
	"sync"
)

var wg sync.WaitGroup
var userBot = userbot.Bot
var adminBot = adminbot.Bot

func main() {

	wg.Add(5)

	go ws.Start()

	go userBot.Start()
	defer userBot.Stop()
	go userbot.SendLimiter.RemoveOldUserRateLimitersCache()

	go adminBot.Start()
	defer adminBot.Stop()

	go adminbot.SendLimiter.RemoveOldUserRateLimitersCache()

	fmt.Println("Bots started...")

	wg.Wait()

}
