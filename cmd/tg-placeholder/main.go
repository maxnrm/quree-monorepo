package main

import (
	"context"
	"fmt"
	"log"
	"quree/config"
	"quree/internal/models"
	"quree/internal/sendlimiter"

	tele "gopkg.in/telebot.v3"
)

var ctx = context.Background()
var limiter = sendlimiter.Init(ctx)

func main() {
	b, err := tele.NewBot(tele.Settings{
		Token:  config.USER_PLACEHOLDER_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10},
	})
	if err != nil {
		panic(err)
	}

	b.Use(MiniLogger())
	b.Handle("/start", startHandler)

	fmt.Println("Bot token:", b.Token)
	b.Start()
	defer b.Stop()
}

func MiniLogger() tele.MiddlewareFunc {
	l := log.Default()

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			chatID := c.Chat().ID
			text := c.Message().Text
			l.Println(chatID, text, "ok")
			return next(c)
		}
	}
}

func startHandler(c tele.Context) error {

	text := "Добро пожаловать! Мой функционал будет доступен в дни фестиваля \"Действуй\" 😌"

	message := &models.SendableMessage{
		Text: &text,
	}

	return message.Send(c.Bot(), limiter)
}
