package adminbot

import (
	"fmt"
	"log"
	"quree/config"
	"quree/internal/models"
	"quree/internal/models/enums"

	"quree/internal/pg"

	tele "gopkg.in/telebot.v3"
)

var db = pg.DB

func idHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("%d", c.Chat().ID))
}

func registerHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)

	user := db.GetUserByChatID(chatID)
	if user != nil {
		sm := models.CreateSendableMessage(SendLimiter, &models.Message{
			Content: "Вы уже зарегистрированы!",
		}, nil)

		return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
	}

	err := db.CreateUser(&models.User{
		ChatID:      fmt.Sprint(c.Chat().ID),
		PhoneNumber: "test",
		Role:        enums.ADMIN,
	})

	if err != nil {
		return c.Send(err.Error())
	}

	sm := models.CreateSendableMessage(SendLimiter, &models.Message{
		Content: "Вы зарегистрированы как Админ!",
	}, nil)

	return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
}

func CheckAuthorize() tele.MiddlewareFunc {
	l := log.Default()

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Message().Text == config.ADMIN_AUTH_CODE {
				return next(c)
			}

			chatID := fmt.Sprint(c.Chat().ID)
			user := db.GetUserByChatID(chatID)
			if user == nil || user.Role != enums.ADMIN {
				sm := models.CreateSendableMessage(SendLimiter, &models.Message{
					Content: "Вы не авторизованы! Для доступа к приложению введите код, полученный у куратора.",
				}, nil)

				l.Println("Админ", chatID, "не авторизован")
				return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
			}

			l.Println("Админ", chatID, "авторизован")
			return next(c)
		}
	}
}
