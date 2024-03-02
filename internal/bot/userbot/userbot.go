package userbot

import (
	"quree/config"
	"quree/internal/bot"
	"time"

	tele "gopkg.in/telebot.v3"
)

var BotConfig = &bot.BotConfig{
	Settings: &tele.Settings{
		Token:  config.USER_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	},

	CommandHandlersMap: map[string]tele.HandlerFunc{
		"/start":    startHandler,
		"/id":       idHandler,
		"/register": registerHandler,
	},

	MiddlewaresMap: &[]tele.MiddlewareFunc{
		bot.MiniLogger(),
	},

	MenuButton: &tele.MenuButton{
		Type: tele.MenuButtonWebApp,
		Text: "Профиль",
		WebApp: &tele.WebApp{
			URL: config.USER_WEBAPP_URL,
		},
	},
}
