package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"quree/config"
	"quree/internal/nats"
	"quree/internal/telegram"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/time/rate"
)

var ctx = context.Background()
var tg *telegram.TelegramClient = telegram.Init(ctx)
var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

var msgStreamConfig = jetstream.StreamConfig{
	Name:      config.NATS_MESSAGES_STREAM,
	Subjects:  []string{config.NATS_MESSAGES_SUBJECT},
	Retention: jetstream.WorkQueuePolicy,
	Storage:   jetstream.FileStorage,
}

var msgConsumerConfig = jetstream.ConsumerConfig{
	Durable:       config.NATS_MESSAGES_CONSUMER,
	AckWait:       2 * time.Second,
	MaxAckPending: 1,
	MemoryStorage: true,
}

type tgUser struct {
	chatID        string
	rateLimiter   *rate.Limiter
	latestMsgSent *timestamp.Timestamp
}

var wg sync.WaitGroup

func main() {

	// create tgUser type, which will store user chat_id, individual user rateLimitier from package rate,
	// that would limit messages at 1 per sec for every user
	// and timestamp, which store info about last message sent time

	tgUsersMap := make(map[string]tgUser)

	// create loop that would clean tgUsersMap every 10 secs from users that haven't sent messages for 10 secs

	go func() {
		for {
			time.Sleep(10 * time.Second)
			// loop thgour tgUsersMap and delete users that haven't sent messages for 10 secs
			for k, v := range tgUsersMap {
				// use time.Since
				if time.Since(time.Unix(v.latestMsgSent.Seconds, 0)) > 10*time.Second {
					delete(tgUsersMap, k)
				}
			}

			fmt.Println("Clean users")
			fmt.Println(tgUsersMap)
		}
	}()

	nc.CreateStream(msgStreamConfig)

	cons := nc.CreateConsumer(msgStreamConfig.Name, msgConsumerConfig)

	mh := createConsumeHandler(&tgUsersMap)

	wg.Add(1)
	cons.Consume(mh)
	wg.Wait()

	defer nc.NC.Close()

	// write graceful shutdown logic
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}

func createConsumeHandler(tgUsersMap *map[string]tgUser) jetstream.MessageHandler {

	return func(msg jetstream.Msg) {
		var msgJSON telegram.Message
		json.Unmarshal(msg.Data(), &msgJSON)

		// create following logic:
		// if user is not in tgUsersMap, add him to tgUsersMap with rateLimiter and timestamp time.Now()
		// if user exists in tgUsersMap, check if he can send message, if he can't, do msg.NackWithWait() for a sec
		// if he can, send message and update timestamp in tgUsersMap

		if _, ok := (*tgUsersMap)[msgJSON.ChatID]; !ok {
			limit := rate.Every(time.Second / time.Duration(config.RATE_LIMIT_PER_USER))
			(*tgUsersMap)[msgJSON.ChatID] = tgUser{
				chatID:        msgJSON.ChatID,
				rateLimiter:   rate.NewLimiter(limit, config.RATE_LIMIT_BURST_PER_USER),
				latestMsgSent: &timestamp.Timestamp{Seconds: 0},
			}
		}

		if !(*tgUsersMap)[msgJSON.ChatID].rateLimiter.Allow() {
			msg.NakWithDelay(900 * time.Millisecond)
			fmt.Println("Nacked user: ", msgJSON.ChatID)
			return
		}

		// print out time of sent with seconds and milliseconds only
		fmt.Println("Sent message at ", time.Now().Format("15:04:05.000"))
		tg.SendMessage(msgJSON)

		latestMsgSent := timestamp.Timestamp{Seconds: time.Now().Unix()}
		(*tgUsersMap)[msgJSON.ChatID] = tgUser{
			chatID:        msgJSON.ChatID,
			rateLimiter:   (*tgUsersMap)[msgJSON.ChatID].rateLimiter,
			latestMsgSent: &latestMsgSent,
		}

		msg.Ack()

	}

}
