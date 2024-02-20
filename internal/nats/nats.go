package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsSettings struct {
	Ctx      context.Context
	URL      string
	Stream   string
	Consumer string
}

type NatsClient struct {
	Ctx      context.Context
	NC       *nats.Conn
	JS       jetstream.JetStream
	Stream   jetstream.Stream
	Consumer jetstream.Consumer
}

func Init(settings NatsSettings) *NatsClient {
	var natsClient NatsClient

	natsClient.NC, _ = nats.Connect(settings.URL)
	// defer natsClient.NC.Drain()

	natsClient.JS, _ = jetstream.New(natsClient.NC)
	natsClient.Stream, _ = natsClient.JS.Stream(settings.Ctx, settings.Stream)
	natsClient.Consumer, _ = natsClient.Stream.Consumer(settings.Ctx, settings.Consumer)

	return &natsClient
}
