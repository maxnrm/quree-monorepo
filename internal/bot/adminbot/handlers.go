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

func startHandler(c tele.Context) error {
	return c.Send("Нажмите SCANNER для сканирования")
}

func registerHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)
	user := db.GetUserByChatIDAndRole(chatID, enums.ADMIN)

	if user != nil {
		sm := models.CreateSendableMessage(SendLimiter, &models.Message{
			Content: "Вы уже зарегистрированы!",
		}, nil)

		c.Bot().SetMenuButton(c.Sender(), &tele.MenuButton{
			Type: tele.MenuButtonWebApp,
			Text: "SCANNER",
			WebApp: &tele.WebApp{
				URL: config.ADMIN_WEBAPP_URL,
			},
		})

		return sm.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
	}

	err := db.CreateUser(&models.User{
		ChatID:      fmt.Sprint(c.Chat().ID),
		PhoneNumber: "",
		Role:        enums.ADMIN,
	})

	if err != nil {
		return c.Send(err.Error())
	}

	c.Bot().SetMenuButton(c.Sender(), &tele.MenuButton{
		Type: tele.MenuButtonWebApp,
		Text: "SCANNER",
		WebApp: &tele.WebApp{
			URL: config.ADMIN_WEBAPP_URL,
		},
	})

	sm := models.CreateSendableMessage(SendLimiter, &models.Message{
		Content: "Вы зарегистрированы как Админ! Нажмите SCANNER для сканирования",
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
			user := db.GetUserByChatIDAndRole(chatID, enums.ADMIN)

			if user == nil {
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
