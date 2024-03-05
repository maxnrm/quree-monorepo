package models

import (
	"quree/internal/sendlimiter"
	"time"

	tele "gopkg.in/telebot.v3"
)

type SendableMessage struct {
	*Message
	Recipient tele.Recipient `json:"recipient"`
	Limiter   *sendlimiter.SendLimiter
	Bot       *tele.Bot
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

func (sm *SendableMessage) sendWithLimit() error {
	chatID := sm.Recipient.Recipient()

	userRateLimiter := sm.Limiter.GetUserRateLimiter(chatID)
	if userRateLimiter == nil {
		sm.Limiter.AddUserRateLimiter(chatID)
		userRateLimiter = sm.Limiter.GetUserRateLimiter(chatID)
	}

	err := userRateLimiter.RateLimiter.Wait(sm.Limiter.Ctx)
	if err != nil {
		return err
	}

	err = sm.Limiter.GlobalRateLimiter.Wait(sm.Limiter.Ctx)
	if err != nil {
		return err
	}

	what := sm.createWhat()

	_, err = sm.Bot.Send(sm.Recipient, what, sm.SendOptions)
	if err != nil {
		return err
	}

	userRateLimiter.LastMsgSent = time.Now()

	return nil

}

func (sm *SendableMessage) Send() error {
	return sm.sendWithLimit()
}
