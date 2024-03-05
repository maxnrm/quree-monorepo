package adminbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"quree/config"
	"time"

	"quree/internal/helpers"
	"quree/internal/models"
	"quree/internal/nats"
	"quree/internal/pg"
	"quree/internal/pg/dbmodels"
	"quree/internal/sendlimiter"

	tele "gopkg.in/telebot.v3"
)

var ctx = context.Background()
var SendLimiter = sendlimiter.Init(ctx)
var db = pg.DB
var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

func Init() *tele.Bot {

	token := config.ADMIN_BOT_TOKEN

	log.Println("bot token:", token)

	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Use(helpers.BotMiniLogger())
	bot.Use(CheckAuthorize())

	bot.Handle("/start", startHandler)
	bot.Handle("/id", idHandler)
	bot.Handle("/register", registerHandler)

	return bot
}

func idHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("%d", c.Chat().ID))
}

func startHandler(c tele.Context) error {
	return c.Send("Нажмите SCANNER для сканирования")
}

func registerHandler(c tele.Context) error {

	menuButton := &tele.MenuButton{
		Type: tele.MenuButtonWebApp,
		Text: "SCANNER",
		WebApp: &tele.WebApp{
			URL: config.ADMIN_WEBAPP_URL,
		},
	}

	chatID := fmt.Sprint(c.Chat().ID)
	user := db.GetAdminByChatID(chatID)
	if user != nil {
		text := "Вы уже зарегистрированы!"

		var message = models.SendableMessage{
			Text: &text,
			Recipient: &models.Recipient{
				ChatID: chatID,
			},
		}

		json, err := json.Marshal(message)
		if err != nil {
			return err
		}

		nc.NC.Publish(config.NATS_ADMIN_MESSAGES_SUBJECT+"."+chatID, json)

		c.Bot().SetMenuButton(c.Sender(), menuButton)
	}

	err := db.CreateAdmin(&dbmodels.Admin{
		ChatID: chatID,
	})
	if err != nil {
		return err
	}

	text := "Вы зарегистрированы как админ! Нажмите SCANNER для сканирования"

	var message = models.SendableMessage{
		Text: &text,
		Recipient: &models.Recipient{
			ChatID: chatID,
		},
	}

	json, err := json.Marshal(message)
	if err != nil {
		return err
	}

	nc.NC.Publish(config.NATS_ADMIN_MESSAGES_SUBJECT+"."+chatID, json)

	c.Bot().SetMenuButton(c.Sender(), menuButton)

	return err
}

func CheckAuthorize() tele.MiddlewareFunc {
	l := log.Default()

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Message().Text == config.ADMIN_AUTH_CODE {
				return next(c)
			}

			chatID := fmt.Sprint(c.Chat().ID)
			user := db.GetAdminByChatID(chatID)

			if user == nil {
				text := "Вы не авторизованы! Для доступа к приложению введите код, полученный у куратора."

				var message = models.SendableMessage{
					Text: &text,
					Recipient: &models.Recipient{
						ChatID: chatID,
					},
				}

				json, err := json.Marshal(message)
				if err != nil {
					return err
				}

				nc.NC.Publish(config.NATS_ADMIN_MESSAGES_SUBJECT+"."+chatID, json)
			}

			l.Println("Админ", chatID, "авторизован")
			return next(c)
		}
	}
}
