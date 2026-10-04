package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/consumer"
	"github.com/rabbitmq/amqp091-go"
	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"
)

const (
	defaultExchange     = "bibit"
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
	conn         *amqp091.Connection
	exchange     string
	groupId      string
	batchSize    int
	batchTimeout time.Duration
}

var _ IBroker = (*Broker)(nil)

func NewBroker(i do.Injector) (*Broker, error) {
	cfg := config.Get().Broker.Rabbitmq

	if cfg.Url == "" {
		return nil, fmt.Errorf("BROKER_RABBITMQ_URL is required")
	}

	if cfg.GroupId == "" {
		return nil, fmt.Errorf("BROKER_RABBITMQ_GROUP_ID is required")
	}

	exchange := cfg.Exchange
	if exchange == "" {
		exchange = defaultExchange
	}

	batchSize := cfg.BatchSize
	if batchSize == 0 {
		batchSize = defaultBatchSize
	}

	batchTimeout := cfg.BatchTimeout
	if batchTimeout == 0 {
		batchTimeout = defaultBatchTimeout
	}

	conn, err := amqp091.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}

	return &Broker{
		conn:         conn,
		exchange:     exchange,
		groupId:      cfg.GroupId,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
	}, nil
}

func (b *Broker) Start(ctx context.Context, routes []consumer.Route) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, route := range routes {
		g.Go(func() error {
			return b.consume(ctx, route)
		})
	}

	return g.Wait()
}

func (b *Broker) Shutdown(ctx context.Context) error {
	slog.Info("shutting down rabbitmq backend")
	return b.conn.Close()
}

func (b *Broker) consume(ctx context.Context, route consumer.Route) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue := b.queueName(route.Topic)

	err = ch.ExchangeDeclare(b.exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind(queue, route.Topic, b.exchange, false, nil)
	if err != nil {
		return err
	}

	err = ch.Qos(b.batchSize, 0, false)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	slog.Info("consuming queue", slog.String("queue", queue))

	for {
		var batch []amqp091.Delivery
		select {
		case <-ctx.Done():
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				return nil
			}
			batch = append(batch, delivery)
		}

		timer := time.NewTimer(b.batchTimeout)
	collect:
		for len(batch) < b.batchSize {
			select {
			case <-timer.C:
				break collect
			case delivery, ok := <-deliveries:
				if !ok {
					break collect
				}
				batch = append(batch, delivery)
			}
		}
		timer.Stop()

		b.handle(ctx, route, batch)
		if ctx.Err() != nil {
			return nil
		}
	}
}

func (b *Broker) handle(ctx context.Context, route consumer.Route, batch []amqp091.Delivery) {
	queue := b.queueName(route.Topic)

	messages := make([]*consumer.Message, len(batch))
	for i, delivery := range batch {
		messages[i] = toMessage(delivery)
	}

	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = route.Consumer.Consume(ctx, messages)
		if err == nil {
			break
		}

		slog.Error("failed to consume messages",
			slog.String("queue", queue),
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
		slog.Error("giving up, discarding the batch",
			slog.String("queue", queue),
			slog.Any("error", err),
		)

		for _, delivery := range batch {
			nackErr := delivery.Nack(false, false)
			if nackErr != nil {
				slog.Error("failed to nack message", slog.String("queue", queue), slog.Any("error", nackErr))
			}
		}

		return
	}

	for _, delivery := range batch {
		ackErr := delivery.Ack(false)
		if ackErr != nil {
			slog.Error("failed to ack message", slog.String("queue", queue), slog.Any("error", ackErr))
		}
	}
}

func (b *Broker) queueName(topic string) string {
	return fmt.Sprintf("%s.%s", b.groupId, topic)
}

func toMessage(delivery amqp091.Delivery) *consumer.Message {
	headers := make(map[string][]byte, len(delivery.Headers))
	for key, value := range delivery.Headers {
		if bytes, ok := value.([]byte); ok {
			headers[key] = bytes
		}
	}

	return &consumer.Message{
		Topic:     delivery.RoutingKey,
		Value:     delivery.Body,
		Headers:   headers,
		Timestamp: delivery.Timestamp,
	}
}
