package consumer

import (
	"context"
	"time"
)

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string][]byte
	Timestamp time.Time
}

type IConsumer interface {
	Consume(ctx context.Context, messages []*Message) error
}

type Route struct {
	Topic    string
	Consumer IConsumer
}
