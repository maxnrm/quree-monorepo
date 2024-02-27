package telegram

import (
	"context"
	"quree/config"
	"time"

	"github.com/monaco-io/request"
	"golang.org/x/time/rate"
)

type TelegramClient struct {
	httpClient request.Client
	limiter    *rate.Limiter
	ctx        context.Context
}

type Message struct {
	ChatID   string `json:"chat_id"`
	Text     string `json:"text"`
	BotToken string `json:"bot_token"`
}

func Init(ctx context.Context) *TelegramClient {

	var tg TelegramClient

	tgAPIBaseURL := config.TG_API_BASE_URL
	botToken := config.USER_BOT_TOKEN
	sendMessageMethodName := "sendMessage"
	sendMessageURL := tgAPIBaseURL + "/bot" + botToken + "/" + sendMessageMethodName

	tg.httpClient = request.Client{
		URL:    sendMessageURL,
		Header: map[string]string{"Content-type": "application/json"},
		Method: "POST",
	}

	limit := rate.Every(time.Second / time.Duration(config.RATE_LIMIT_GLOBAL))
	tg.limiter = rate.NewLimiter(limit, config.RATE_LIMIT_BURST_GLOBAL)

	tg.ctx = ctx

	return &tg

}

func (tg *TelegramClient) SendMessage(msg Message) error {
	err := (*tg.limiter).Wait(tg.ctx)

	if err != nil {
		return err
	}

	tg.httpClient.JSON = msg
	tg.httpClient.Send()

	return nil
}
