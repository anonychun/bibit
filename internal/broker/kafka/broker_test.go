package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anonychun/bibit/internal/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestToMessage(t *testing.T) {
	t.Run("maps kafka message fields", func(t *testing.T) {
		timestamp := time.Date(2026, 10, 4, 9, 0, 0, 0, time.UTC)
		kmsg := kafka.Message{
			Topic:   "hello",
			Key:     []byte("user-1"),
			Value:   []byte("hello world"),
			Headers: []kafka.Header{{Key: "trace_id", Value: []byte("abc")}},
			Time:    timestamp,
		}

		message := toMessage(kmsg)

		assert.Equal(t, "hello", message.Topic)
		assert.Equal(t, []byte("user-1"), message.Key)
		assert.Equal(t, []byte("hello world"), message.Value)
		assert.Equal(t, map[string][]byte{"trace_id": []byte("abc")}, message.Headers)
		assert.Equal(t, timestamp, message.Timestamp)
	})

	t.Run("maps a message without headers and key", func(t *testing.T) {
		message := toMessage(kafka.Message{Topic: "hello", Value: []byte("hello world")})

		assert.Empty(t, message.Key)
		assert.Empty(t, message.Headers)
	})
}

func TestHandle(t *testing.T) {
	b := &Broker{backoff: time.Millisecond}

	t.Run("retries until max attempts, then stops the consumer", func(t *testing.T) {
		mc := consumer.NewMockIConsumer(t)
		mc.On("Consume", mock.Anything, mock.Anything).Return(errors.New("boom")).Times(maxAttempts)

		route := consumer.Route{Topic: "hello", Consumer: mc}
		batch := []kafka.Message{{Topic: "hello", Value: []byte("hello world")}}

		err := b.handle(context.Background(), nil, route, batch)

		require.Error(t, err)
		assert.ErrorContains(t, err, "boom")
		assert.ErrorContains(t, err, `"hello"`)
		mc.AssertNumberOfCalls(t, "Consume", maxAttempts)
	})

	t.Run("a successful batch fails to commit on a reader without a consumer group", func(t *testing.T) {
		mc := consumer.NewMockIConsumer(t)
		mc.On("Consume", mock.Anything, mock.Anything).Return(nil).Once()

		reader := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{"127.0.0.1:9092"}, Topic: "hello"})

		route := consumer.Route{Topic: "hello", Consumer: mc}
		batch := []kafka.Message{{Topic: "hello", Value: []byte("hello world")}}

		err := b.handle(context.Background(), reader, route, batch)

		require.Error(t, err)
		assert.ErrorContains(t, err, "commit")
		mc.AssertNumberOfCalls(t, "Consume", 1)
	})
}
