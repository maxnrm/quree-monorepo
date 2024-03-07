package main

import (
	"context"
	"encoding/json"
	"fmt"
	"quree/config"
	"quree/internal/models"
	"quree/internal/nats"
	"quree/internal/sendlimiter"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	tele "gopkg.in/telebot.v3"
)

var wg sync.WaitGroup
var ctx = context.Background()
var userSl = sendlimiter.Init(ctx)
var adminSl = sendlimiter.Init(ctx)

var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

var userBotSender, _ = tele.NewBot(tele.Settings{
	Token:  config.USER_BOT_TOKEN,
	Poller: &tele.LongPoller{Timeout: 10 * time.Second},
})

var adminBotSender, _ = tele.NewBot(tele.Settings{
	Token:  config.ADMIN_BOT_TOKEN,
	Poller: &tele.LongPoller{Timeout: 10 * time.Second},
})

var streamConfig = jetstream.StreamConfig{
	Name:      config.NATS_MESSAGES_STREAM,
	Subjects:  []string{config.NATS_RECEIVER_MESSAGES_SUBJECT},
	Retention: jetstream.WorkQueuePolicy,
	Storage:   jetstream.FileStorage,
}

var userConsumerConfig = jetstream.ConsumerConfig{
	Name:          config.NATS_USER_MESSAGES_CONSUMER,
	Durable:       config.NATS_USER_MESSAGES_CONSUMER,
	FilterSubject: config.NATS_USER_MESSAGES_SUBJECT,
	AckWait:       2 * time.Second,
	MaxAckPending: 60,
	MemoryStorage: true,
}

var adminConsumerConfig = jetstream.ConsumerConfig{
	Name:          config.NATS_ADMIN_MESSAGES_CONSUMER,
	Durable:       config.NATS_ADMIN_MESSAGES_CONSUMER,
	FilterSubject: config.NATS_ADMIN_MESSAGES_SUBJECT,
	AckWait:       2 * time.Second,
	MaxAckPending: 60,
	MemoryStorage: true,
}

func main() {
	go userSl.RemoveOldUserRateLimitersCache()
	go adminSl.RemoveOldUserRateLimitersCache()

	nc.CreateStream(streamConfig)

	userCons := nc.CreateConsumer(streamConfig.Name, userConsumerConfig)
	userMessageHandler := createConsumeHandler(ctx, userBotSender, userSl)

	adminCons := nc.CreateConsumer(streamConfig.Name, adminConsumerConfig)
	adminMessageHandler := createConsumeHandler(ctx, adminBotSender, adminSl)

	wg.Add(2)

	userCons.Consume(userMessageHandler)
	adminCons.Consume(adminMessageHandler)

	fmt.Println("Consuming...")

	wg.Wait()
}

func createConsumeHandler(ctx context.Context, bot *tele.Bot, limiter *sendlimiter.SendLimiter) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		var sendableMessage models.SendableMessage

		err := json.Unmarshal(msg.Data(), &sendableMessage)
		if err != nil {
			fmt.Println("Error while unmarshalling sendableMessage from json:", err)
			return
		}

		msg.DoubleAck(ctx)

		go sendableMessage.Send(bot, limiter)
	}
}
