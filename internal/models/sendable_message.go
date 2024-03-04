package models

import (
	"fmt"
	"quree/config"
	"quree/internal/sendlimiter"
	"time"

	tele "gopkg.in/telebot.v3"
)

type SendableMessage struct {
	*Message
	*File
	SendLimiter *sendlimiter.SendLimiter
}

func CreateSendableMessage(sl *sendlimiter.SendLimiter, m *Message, f *File) *SendableMessage {
	return &SendableMessage{
		Message:     m,
		File:        f,
		SendLimiter: sl,
	}
}

func (sm *SendableMessage) createWhat() interface{} {
	var what interface{}
	if sm.Content != "" && sm.File != nil {
		what = &tele.Photo{File: tele.FromURL(config.IMGPROXY_PUBLIC_URL + "/" + sm.File.Filename), Caption: sm.Content}
		fmt.Println(what.(*tele.Photo).File.FileURL)
	} else if sm.Content != "" {
		what = sm.Content
	} else if sm.File != nil {
		what = &tele.Photo{File: tele.FromURL(config.IMGPROXY_PUBLIC_URL + "/" + sm.File.Filename)}
	}

	return what
}

func (sm *SendableMessage) sendWithLimit(b *tele.Bot, r tele.Recipient, opt *tele.SendOptions) error {
	chatID := r.Recipient()

	userRateLimiter := sm.SendLimiter.GetUserRateLimiter(chatID)
	if userRateLimiter == nil {
		sm.SendLimiter.AddUserRateLimiter(chatID)
		userRateLimiter = sm.SendLimiter.GetUserRateLimiter(chatID)
	}

	err := userRateLimiter.RateLimiter.Wait(sm.SendLimiter.Ctx)
	if err != nil {
		return err
	}

	err = sm.SendLimiter.GlobalRateLimiter.Wait(sm.SendLimiter.Ctx)
	if err != nil {
		return err
	}

	what := sm.createWhat()

	_, err = b.Send(r, what, opt)
	if err != nil {
		return err
	}

	userRateLimiter.LastMsgSent = time.Now()

	return nil

}

func (sm *SendableMessage) Send(b *tele.Bot, r tele.Recipient, opt *tele.SendOptions) error {
	return sm.sendWithLimit(b, r, opt)
}
