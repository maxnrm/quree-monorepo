package main

import bot "quree/internal/bot"

func main() {

	tgBot := bot.Init()

	tgBot.Start()

}
