package kafka

import (
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
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
