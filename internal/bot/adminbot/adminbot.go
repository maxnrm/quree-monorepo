package adminbot

import (
	"context"
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

	"github.com/google/uuid"
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

	bot.Handle("/id", idHandler)
	bot.Handle("/start", registerHandler)
	bot.Handle(config.ADMIN_AUTH_CODE, registerHandler)

	// handle buttons
	bot.Handle(replyButtons["start"], registerHandler)
	bot.Handle(replyButtons["help"], helpHandler)

	nc.UsePublishSubject(config.NATS_ADMIN_MESSAGES_SUBJECT)

	return bot
}

func idHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("%d", c.Chat().ID))
}

func helpHandler(c tele.Context) error {
	return c.Send("HELP_ADMIN1")
}

func registerHandler(c tele.Context) error {

	chatID := fmt.Sprint(c.Chat().ID)
	user := db.GetAdminByChatID(chatID)
	if user != nil {
		var message = &models.SendableMessage{
			Text: &textAlreadyRegistered,
			Recipient: &models.Recipient{
				ChatID: chatID,
			},
			SendOptions: &tele.SendOptions{
				ReplyMarkup: menuAuthorized,
			},
		}

		nc.Publish(message)

		c.Bot().SetMenuButton(c.Sender(), tele.MenuButtonDefault)

		return nil
	}

	err := db.CreateAdmin(&dbmodels.Admin{
		ID:          uuid.New().String(),
		DateCreated: time.Now(),
		ChatID:      chatID,
	})
	if err != nil {
		return err
	}

	var message = &models.SendableMessage{
		Text: &textRegistered,
		Recipient: &models.Recipient{
			ChatID: chatID,
		},
	}

	nc.Publish(message)

	c.Bot().SetMenuButton(c.Sender(), tele.MenuButtonDefault)

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
				var message = &models.SendableMessage{
					Text: &textUnauthorized,
					Recipient: &models.Recipient{
						ChatID: chatID,
					},
					SendOptions: &tele.SendOptions{
						ReplyMarkup: menuUnauthorized,
					},
				}

				nc.Publish(message)

				return nil
			}

			l.Println("Админ", chatID, "авторизован")
			return next(c)
		}
	}
}

var webApp = &tele.WebApp{
	URL: config.ADMIN_WEBAPP_URL,
}

var textUnauthorized = "Вы не авторизованы! Для доступа к приложению введите код, полученный у куратора"
var textRegistered = "Вы авторизированы как админ! Нажмите Сканер QR для сканирования"
var textAlreadyRegistered = "Вы уже авторизированы! Нажмите Сканер QR для сканирования"

var menuAuthorized = &tele.ReplyMarkup{
	ResizeKeyboard: true,
	ReplyKeyboard:  replyKeyboardAuthorized,
}

var menuUnauthorized = &tele.ReplyMarkup{
	ResizeKeyboard: true,
	ReplyKeyboard:  replyKeyboardUnauthorized,
}

var replyKeyboardAuthorized = [][]tele.ReplyButton{
	{*replyButtons["help"]},
	{*replyButtons["scanner"]},
}

var replyKeyboardUnauthorized = [][]tele.ReplyButton{
	{*replyButtons["start"]},
}

var replyButtons = map[string]*tele.ReplyButton{
	"start":   {Text: "Начать"},
	"help":    {Text: "Как это работает?"},
	"scanner": {Text: "Сканер QR", WebApp: webApp},
}
