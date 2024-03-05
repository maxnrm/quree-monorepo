package models

import (
	"quree/internal/sendlimiter"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Recipient struct {
	ChatID string
}

func (r *Recipient) Recipient() string {
	return r.ChatID
}

type SendableMessage struct {
	Text        *string           `json:"text"`
	Caption     *string           `json:"caption"`
	Photo       *tele.Photo       `json:"photo"`
	SendOptions *tele.SendOptions `json:"send_options"`
	Variant     int               `json:"variant"`
	Recipient   *Recipient        `json:"recipient"`
}

func (sm *SendableMessage) createWhat() interface{} {
	var what interface{}

	if sm.Text != nil {
		what = sm.Text
	} else {
		what = sm.Photo
	}

	return what
}

func (sm *SendableMessage) sendWithLimit(bot *tele.Bot, limiter *sendlimiter.SendLimiter) error {
	chatID := sm.Recipient.Recipient()

	userRateLimiter := limiter.GetUserRateLimiter(chatID)
	if userRateLimiter == nil {
		limiter.AddUserRateLimiter(chatID)
		userRateLimiter = limiter.GetUserRateLimiter(chatID)
	}

	err := userRateLimiter.RateLimiter.Wait(limiter.Ctx)
	if err != nil {
		return err
	}

	err = limiter.GlobalRateLimiter.Wait(limiter.Ctx)
	if err != nil {
		return err
	}

	what := sm.createWhat()

	_, err = bot.Send(sm.Recipient, what, sm.SendOptions)
	if err != nil {
		return err
	}

	userRateLimiter.LastMsgSent = time.Now()

	return nil

}

func (sm *SendableMessage) Send(bot *tele.Bot, limiter *sendlimiter.SendLimiter) error {
	return sm.sendWithLimit(bot, limiter)
}
