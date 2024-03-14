package main

import (
	"context"
	"fmt"
	"log"
	"quree/config"
	"quree/internal/models"
	"quree/internal/pg"
	"quree/internal/sendlimiter"

	tele "gopkg.in/telebot.v3"
)

var ctx = context.Background()
var limiter = sendlimiter.Init(ctx, config.RATE_LIMIT_GLOBAL, config.RATE_LIMIT_BURST_GLOBAL)

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
	b.Handle("/stats", statsHandler)

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

var db = pg.DB

func startHandler(c tele.Context) error {

	text := "Добро пожаловать! Мой функционал будет доступен в дни фестиваля \"Действуй\" 😌"

	message := &models.SendableMessage{
		Text: &text,
		Recipient: &models.Recipient{
			ChatID: fmt.Sprint(c.Chat().ID),
		},
	}

	return message.Send(c.Bot(), limiter)
}

func statsHandler(c tele.Context) error {
	var text = ""

	numberOfusers := db.CountUsers()
	numberOfAdmins := db.CountAdmins()
	numberOfVisits := db.CountVisits()

	text += "Всего юзеров: *" + fmt.Sprint(numberOfusers) + "*\n\n"
	text += "Всего админов: *" + fmt.Sprint(numberOfAdmins) + "*\n\n"
	text += "Всего визитов: *" + fmt.Sprint(numberOfVisits) + "*\n\n"

	text += "Юзеров, посетившых 4 и более событий: *" + fmt.Sprint(db.CountUsersWithMoreThanFourVisits()) + "*\n\n"
	text += "Юзеров, посетивщих викторину: *" + fmt.Sprint(db.CountUsersWithQuiz()) + "*\n\n"
	text += "Юзеров, посетивших 4 и более событий и викторину: *" + fmt.Sprint(db.CountUsersWithMoreThanFourVisitsAndQuiz()) + "*\n\n"

	message := &models.SendableMessage{
		Text: &text,
		Recipient: &models.Recipient{
			ChatID: fmt.Sprint(c.Chat().ID),
		},
		SendOptions: &tele.SendOptions{
			ParseMode: tele.ModeMarkdownV2,
		},
	}

	return message.Send(c.Bot(), limiter)
}
