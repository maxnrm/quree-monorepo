package bot

import (
	"fmt"
	"log"
	"quree/config"
	"quree/internal/models"
	"quree/internal/models/enums"
	"quree/internal/pg"
	"time"

	tele "gopkg.in/telebot.v3"
)

var db = pg.DB

type botConfig struct {
	settings           *tele.Settings
	commandHandlersMap map[string]tele.HandlerFunc
	middlewaresMap     *[]tele.MiddlewareFunc
	menuButton         *tele.MenuButton
}

// var setProgramBtns [][]tele.ReplyButton = [][]tele.ReplyButton{
// 	{tele.ReplyButton{Text: "Начать", Contact: true}},
// }

var userBotConfig = &botConfig{
	settings: &tele.Settings{
		Token:  config.USER_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	},

	commandHandlersMap: map[string]tele.HandlerFunc{
		"/start":    startHandler,
		"/id":       idHandler,
		"/register": registerHandler,
		"/me":       getUserHandler,
	},

	middlewaresMap: &[]tele.MiddlewareFunc{
		miniLogger(),
		// ensureLoginMiddleware(),
	},

	menuButton: &tele.MenuButton{
		Type: tele.MenuButtonWebApp,
		Text: "Профиль",
		WebApp: &tele.WebApp{
			URL: config.USER_WEBAPP_URL,
		},
	},
}

var adminBotConfig = &botConfig{
	settings: &tele.Settings{
		Token:  config.ADMIN_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	},

	commandHandlersMap: map[string]tele.HandlerFunc{
		"/start": startHandler,
		"/id":    idHandler,
	},

	middlewaresMap: &[]tele.MiddlewareFunc{
		miniLogger(),
		// ensureLoginMiddleware(),
	},

	menuButton: &tele.MenuButton{
		Type: tele.MenuButtonWebApp,
		Text: "Сканер QR",
		WebApp: &tele.WebApp{
			URL: config.ADMIN_WEBAPP_URL,
		},
	},
}

var UserBot = Init(userBotConfig)
var AdminBot = Init(adminBotConfig)

func Init(c *botConfig) *tele.Bot {

	log.Println("bot token:", c.settings.Token)

	b, err := tele.NewBot(*c.settings)
	if err != nil {
		log.Fatal(err)
	}

	//	b.SetMenuButton(b.Me, c.menuButton)

	for _, m := range *c.middlewaresMap {
		b.Use(m)
	}

	for c, h := range c.commandHandlersMap {
		b.Handle(c, h)
	}

	return b
}

func startHandler(c tele.Context) error {

	msgs := db.GetMessagesByType(enums.START)

	for _, m := range msgs {
		if m.Content != "" {
			c.Send(m.Content)
		}
	}

	return c.Send("Start!")
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

func getUserHandler(c tele.Context) error {
	ci := fmt.Sprint(c.Chat().ID)
	chatID := db.GetUserByChatID(ci).ChatID
	fmt.Println(chatID)
	return c.Send(chatID)
}

func miniLogger() tele.MiddlewareFunc {
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
