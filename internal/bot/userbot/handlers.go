package userbot

import (
	"fmt"
	"quree/config"
	"quree/internal/models"
	"quree/internal/models/enums"
	"sort"

	"quree/internal/pg"

	tele "gopkg.in/telebot.v3"
)

var db = pg.DB

func startHandler(c tele.Context) error {

	msgs := db.GetMessagesByType(enums.START)

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Sort < msgs[j].Sort
	})

	var what interface{}

	for _, m := range msgs {

		if m.Content != "" && m.Image != "" {
			file := db.GetFileRecordByID(m.Image)
			what = &tele.Photo{File: tele.FromURL(config.IMGPROXY_PUBLIC_URL + "/" + file.Filename), Caption: m.Content}
		} else if m.Content != "" {
			what = m.Content
		} else if m.Image != "" {
			file := db.GetFileRecordByID(m.Image)
			what = &tele.Photo{File: tele.FromURL(config.IMGPROXY_PUBLIC_URL + "/" + file.Filename)}
		}

		SendLimiter.LimitSend(c, fmt.Sprint(c.Chat().ID), what)
		c.Bot().Send(c.Chat(), what)
	}

	return nil
}

func idHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("%d", c.Chat().ID))
}

func registerHandler(c tele.Context) error {

	err := db.CreateUser(&models.User{
		ChatID:      fmt.Sprint(c.Chat().ID),
		PhoneNumber: "test",
		Role:        enums.USER,
		QRCode:      "53c0d5b2-3b92-4630-ba4a-58721f0df1f5",
	})

	if err != nil {
		return c.Send(err.Error())
	}

	return c.Send("Register!")
}
