package bot

import (
	"log"
	"quree/internal/sendlimiter"

	tele "gopkg.in/telebot.v3"
)

type BotConfig struct {
	SendLimiter        *sendlimiter.SendLimiter
	Settings           *tele.Settings
	CommandHandlersMap map[string]tele.HandlerFunc
	MiddlewaresMap     *[]tele.MiddlewareFunc
	MenuButton         *tele.MenuButton
}

func Init(c *BotConfig) *tele.Bot {

	log.Println("bot token:", c.Settings.Token)

	b, err := tele.NewBot(*c.Settings)
	if err != nil {
		log.Fatal(err)
	}

	if c.MiddlewaresMap != nil {
		for _, m := range *c.MiddlewaresMap {
			b.Use(m)
		}
	}

	if c.CommandHandlersMap != nil {
		for c, h := range c.CommandHandlersMap {
			b.Handle(c, h)
		}
	}

	return b
}

func MiniLogger() tele.MiddlewareFunc {
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
