package main

import (
	"quree/internal/bot/adminbot"
	"quree/internal/bot/userbot"
	ws "quree/internal/webserver"
	"sync"
)

var wg sync.WaitGroup
var userBot = userbot.Init()
var adminBot = adminbot.Init()

func main() {

	wg.Add(3)

	go ws.Start()

	go userBot.Start()
	defer userBot.Stop()

	go adminBot.Start()
	defer adminBot.Stop()

	wg.Wait()

}
