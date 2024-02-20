package bot

import (
	"log"
	"quree/config"
	"time"

	tele "gopkg.in/telebot.v3"
)

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
		"/start": startHandler,
	},

	middlewaresMap: &[]tele.MiddlewareFunc{
		// ensureLoginMiddleware(),
		miniLogger(),
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
	},

	middlewaresMap: &[]tele.MiddlewareFunc{
		// ensureLoginMiddleware(),
		miniLogger(),
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
	return c.Send("Start!")
}

func ensureLoginMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().IsBot {
				return nil
			}
			return c.Send("Start!")
		}
	}
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
