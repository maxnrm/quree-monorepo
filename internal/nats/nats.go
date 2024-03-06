package nats

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"quree/internal/models"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsSettings struct {
	Ctx context.Context
	URL string
}

type NatsClient struct {
	Ctx context.Context
	NC  *nats.Conn
	JS  jetstream.JetStream
	PS  *string
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

func (nc *NatsClient) UsePublishSubject(subject string) {
	nc.PS = &subject
}

func (nc *NatsClient) Publish(message *models.SendableMessage) {
	if nc.PS == nil {
		err := errors.New("no subject was set for publishing")
		log.Fatal(err)
		return
	}

	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	if message.Photo != nil {
		log.Println("publishing photo to", *nc.PS)
		nc.NC.Publish(*nc.PS, json)
	}

	if message.Text != nil {
		log.Println("publishing message", *message.Text, "to", *nc.PS)
		nc.NC.Publish(*nc.PS, json)
	}
}
