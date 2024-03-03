package userbot

import (
	"context"
	"quree/config"
	"quree/internal/bot"
	"quree/internal/sendlimiter"
	"time"

	tele "gopkg.in/telebot.v3"
)

var ctx = context.Background()
var Bot *tele.Bot = bot.Init(BotConfig)
var SendLimiter = sendlimiter.Init(ctx)

var BotConfig = &bot.BotConfig{
	Settings: &tele.Settings{
		Token:  config.USER_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	},

	CommandHandlersMap: map[string]tele.HandlerFunc{
		"/start":    startHandler,
		"/id":       idHandler,
		"/register": registerHandler,
		"/qr":       qrHandler,
	},

	MiddlewaresMap: &[]tele.MiddlewareFunc{
		bot.MiniLogger(),
		CheckAuthorize(),
	},

	MenuButton: &tele.MenuButton{
		Type: tele.MenuButtonWebApp,
		Text: "Профиль",
		WebApp: &tele.WebApp{
			URL: config.USER_WEBAPP_URL,
		},
	},
}
