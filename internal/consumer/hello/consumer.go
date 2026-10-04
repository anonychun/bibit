package hello

import (
	"context"
	"log/slog"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/consumer"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewConsumer)
}

type Consumer struct{}

var _ consumer.IConsumer = (*Consumer)(nil)

func NewConsumer(i do.Injector) (*Consumer, error) {
	return &Consumer{}, nil
}

func (c *Consumer) Consume(ctx context.Context, messages []*consumer.Message) error {
	for _, message := range messages {
		slog.Info("hello",
			slog.String("topic", message.Topic),
			slog.String("value", string(message.Value)),
		)
	}

	return nil
}
