package nats

import (
	"context"
	"log"

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
	Ctx context.Context
	NC  *nats.Conn
	JS  jetstream.JetStream
}

func Init(settings NatsSettings) *NatsClient {
	var natsClient NatsClient

	natsClient.Ctx = settings.Ctx

	natsClient.NC, _ = nats.Connect(settings.URL)

	natsClient.JS, _ = jetstream.New(natsClient.NC)

	return &natsClient
}

func (nc *NatsClient) CreateStream(streamConfig jetstream.StreamConfig) *jetstream.Stream {
	s, err := nc.JS.CreateOrUpdateStream(nc.Ctx, streamConfig)
	if err != nil {
		log.Fatal("Error creating stream", err)
	}

	return &s
}

func (nc *NatsClient) CreateConsumer(stream string, consumerConfig jetstream.ConsumerConfig) jetstream.Consumer {
	c, err := nc.JS.CreateOrUpdateConsumer(nc.Ctx, stream, consumerConfig)
	if err != nil {
		panic(err)
	}

	return c
}
