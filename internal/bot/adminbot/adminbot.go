package adminbot

import (
	"context"
	"quree/config"
	"time"

	"quree/internal/bot"
	"quree/internal/sendlimiter"

	tele "gopkg.in/telebot.v3"
)

var ctx = context.Background()
var Bot *tele.Bot = bot.Init(BotConfig)
var SendLimiter = sendlimiter.Init(ctx)

var BotConfig = &bot.BotConfig{
	Settings: &tele.Settings{
		Token:  config.ADMIN_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	},

	CommandHandlersMap: map[string]tele.HandlerFunc{
		config.ADMIN_AUTH_CODE: registerHandler,
		"/id":                  idHandler,
	},

	MiddlewaresMap: &[]tele.MiddlewareFunc{
		bot.MiniLogger(),
		CheckAuthorize(),
	},
}
