package main

import (
	"context"
	"encoding/json"
	"quree/config"
	"quree/internal/bot"
	"quree/internal/models"
	"quree/internal/nats"
	"quree/internal/pg"
	"quree/internal/sendlimiter"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	tele "gopkg.in/telebot.v3"
)

var wg sync.WaitGroup
var db = pg.DB
var ctx = context.Background()
var sl = sendlimiter.Init(ctx)
var botSender = bot.Init(&bot.BotConfig{
	Settings: &tele.Settings{
		Token:  config.USER_BOT_TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	},
})

var streamConfig = jetstream.StreamConfig{
	Name:      config.NATS_MESSAGES_STREAM,
	Subjects:  []string{config.NATS_MESSAGES_SUBJECT},
	Retention: jetstream.WorkQueuePolicy,
	Storage:   jetstream.FileStorage,
}

var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

var consumerConfig = jetstream.ConsumerConfig{
	Durable:       config.NATS_MESSAGES_CONSUMER,
	AckWait:       2 * time.Second,
	MaxAckPending: 60,
	MemoryStorage: true,
}

func main() {
	wg.Add(1)
	go sl.RemoveOldUserRateLimitersCache()

	nc.CreateStream(streamConfig)

	cons := nc.CreateConsumer(streamConfig.Name, consumerConfig)

	messageHandler := createConsumeHandler(botSender, sl)

	cons.Consume(messageHandler)

	wg.Wait()
}

func createConsumeHandler(bot *tele.Bot, sl *sendlimiter.SendLimiter) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		var msgJSON models.MessageWithRecipient
		json.Unmarshal(msg.Data(), &msgJSON)

		var sm *models.SendableMessage

		if msgJSON.Image != "" {
			file := db.GetFileRecordByID(msgJSON.Image)
			sm = models.CreateSendableMessage(sl, &msgJSON.Message, file)
		} else {
			sm = models.CreateSendableMessage(sl, &msgJSON.Message, nil)
		}

		sm.Send(bot, msgJSON, &tele.SendOptions{})
	}
}
