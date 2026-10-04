package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/consumer"
	"github.com/samber/do/v2"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

const (
	defaultBatchSize    = 100
	defaultBatchTimeout = 500 * time.Millisecond
	maxAttempts         = 3
)

func init() {
	do.Provide(bootstrap.Injector, NewBroker)
}

type IBroker interface {
	Start(ctx context.Context, routes []consumer.Route) error
}

type Broker struct {
	brokers      []string
	groupId      string
	batchSize    int
	batchTimeout time.Duration
}

var _ IBroker = (*Broker)(nil)

func NewBroker(i do.Injector) (*Broker, error) {
	cfg := config.Get().Broker.Kafka

	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("BROKER_KAFKA_BROKERS is required")
	}

	if cfg.GroupId == "" {
		return nil, fmt.Errorf("BROKER_KAFKA_GROUP_ID is required")
	}

	batchSize := cfg.BatchSize
	if batchSize == 0 {
		batchSize = defaultBatchSize
	}

	batchTimeout := cfg.BatchTimeout
	if batchTimeout == 0 {
		batchTimeout = defaultBatchTimeout
	}

	return &Broker{
		brokers:      cfg.Brokers,
		groupId:      cfg.GroupId,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
	}, nil
}

func (b *Broker) Start(ctx context.Context, routes []consumer.Route) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, route := range routes {
		reader := b.newReader(route.Topic)

		g.Go(func() error {
			defer reader.Close()
			return b.consume(ctx, reader, route)
		})
	}

	return g.Wait()
}

func (b *Broker) newReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     b.brokers,
		GroupID:     b.groupId,
		Topic:       topic,
		StartOffset: kafka.FirstOffset,
	})
}

func (b *Broker) consume(ctx context.Context, reader *kafka.Reader, route consumer.Route) error {
	slog.Info("consuming topic", slog.String("topic", route.Topic))

	fetches := make(chan fetchResult)
	go b.fetch(ctx, reader, fetches)

	for {
		var batch []kafka.Message
		select {
		case <-ctx.Done():
			return nil
		case result, ok := <-fetches:
			if !ok {
				return nil
			}
			if result.err != nil {
				return result.err
			}
			batch = append(batch, result.msg)
		}

		timer := time.NewTimer(b.batchTimeout)
	collect:
		for len(batch) < b.batchSize {
			select {
			case <-timer.C:
				break collect
			case result, ok := <-fetches:
				if !ok {
					break collect
				}
				if result.err != nil {
					timer.Stop()
					return result.err
				}
				batch = append(batch, result.msg)
			}
		}
		timer.Stop()

		b.handle(ctx, reader, route, batch)
	}
}

func (b *Broker) fetch(ctx context.Context, reader *kafka.Reader, fetches chan<- fetchResult) {
	defer close(fetches)

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, io.EOF) {
				return
			}

			select {
			case fetches <- fetchResult{err: err}:
			case <-ctx.Done():
			}
			return
		}

		select {
		case fetches <- fetchResult{msg: msg}:
		case <-ctx.Done():
			return
		}
	}
}

func (b *Broker) handle(ctx context.Context, reader *kafka.Reader, route consumer.Route, batch []kafka.Message) {
	messages := make([]*consumer.Message, len(batch))
	for i, msg := range batch {
		messages[i] = toMessage(msg)
	}

	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = route.Consumer.Consume(ctx, messages)
		if err == nil {
			break
		}

		slog.Error("failed to consume messages",
			slog.String("topic", route.Topic),
			slog.Int("attempt", attempt),
			slog.Any("error", err),
		)

		if attempt == maxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}

	if err != nil {
		slog.Error("giving up, the batch will be redelivered after a restart",
			slog.String("topic", route.Topic),
			slog.Any("error", err),
		)
		return
	}

	err = reader.CommitMessages(context.WithoutCancel(ctx), batch...)
	if err != nil {
		slog.Error("failed to commit messages",
			slog.String("topic", route.Topic),
			slog.Any("error", err),
		)
	}
}

func toMessage(kmsg kafka.Message) *consumer.Message {
	headers := make(map[string][]byte, len(kmsg.Headers))
	for _, header := range kmsg.Headers {
		headers[header.Key] = header.Value
	}

	return &consumer.Message{
		Topic:     kmsg.Topic,
		Key:       kmsg.Key,
		Value:     kmsg.Value,
		Headers:   headers,
		Timestamp: kmsg.Time,
	}
}

type fetchResult struct {
	msg kafka.Message
	err error
}
