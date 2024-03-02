package adminbot

import (
	"quree/config"
	"time"

	"quree/internal/bot"

	tele "gopkg.in/telebot.v3"
)

var BotConfig = &bot.BotConfig{
	Settings: &tele.Settings{
		Token:  config.ADMIN_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	},

	CommandHandlersMap: map[string]tele.HandlerFunc{
		"/start": startHandler,
		"/id":    idHandler,
	},

	MiddlewaresMap: &[]tele.MiddlewareFunc{
		bot.MiniLogger(),
	},

	MenuButton: &tele.MenuButton{
		Type: tele.MenuButtonWebApp,
		Text: "Сканер QR",
		WebApp: &tele.WebApp{
			URL: config.ADMIN_WEBAPP_URL,
		},
	},
}
