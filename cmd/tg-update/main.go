package main

import (
	"fmt"
	"quree/config"
	"quree/internal/bot/adminbot"
	"quree/internal/bot/userbot"
	"sync"
)

var wg sync.WaitGroup
var userBot = userbot.Bot
var adminBot = adminbot.Bot

func main() {

	fmt.Println("test")
	fmt.Println(config.ADMIN_AUTH_CODE)
	wg.Add(4)

	go userBot.Start()
	go userbot.SendLimiter.RemoveOldUserRateLimitersCache()

	go adminBot.Start()
	go adminbot.SendLimiter.RemoveOldUserRateLimitersCache()

	wg.Wait()
}
