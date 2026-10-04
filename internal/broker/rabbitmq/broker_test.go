package rabbitmq

import (
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestToMessage(t *testing.T) {
	t.Run("maps delivery fields", func(t *testing.T) {
		timestamp := time.Date(2026, 10, 4, 9, 0, 0, 0, time.UTC)
		delivery := amqp091.Delivery{
			RoutingKey: "hello",
			Body:       []byte("hello world"),
			Headers:    amqp091.Table{"trace_id": []byte("abc")},
			Timestamp:  timestamp,
		}

		message := toMessage(delivery)

		assert.Equal(t, "hello", message.Topic)
		assert.Equal(t, []byte("hello world"), message.Value)
		assert.Equal(t, map[string][]byte{"trace_id": []byte("abc")}, message.Headers)
		assert.Equal(t, timestamp, message.Timestamp)
	})

	t.Run("skips headers that are not bytes", func(t *testing.T) {
		delivery := amqp091.Delivery{
			RoutingKey: "hello",
			Body:       []byte("hello world"),
			Headers:    amqp091.Table{"retries": int64(3)},
		}

		message := toMessage(delivery)

		assert.Empty(t, message.Headers)
	})
}

func TestBackend_QueueName(t *testing.T) {
	backend := &Broker{groupId: "bibit"}

	assert.Equal(t, "bibit.hello", backend.queueName("hello"))
}
