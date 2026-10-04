package worker

import (
	"testing"

	"github.com/anonychun/bibit/internal/consumer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddRoutes(t *testing.T) {
	handler := consumer.NewMockIConsumer(t)

	t.Run("appends routes", func(t *testing.T) {
		routes := []consumer.Route{}

		err := addRoutes(&routes,
			consumer.Route{Topic: "hello", Consumer: handler},
			consumer.Route{Topic: "goodbye", Consumer: handler},
		)

		require.NoError(t, err)
		require.Len(t, routes, 2)
		assert.Equal(t, "hello", routes[0].Topic)
		assert.Equal(t, "goodbye", routes[1].Topic)
	})

	t.Run("rejects a route without handler", func(t *testing.T) {
		routes := []consumer.Route{}

		err := addRoutes(&routes, consumer.Route{Topic: "hello"})

		assert.ErrorContains(t, err, `"hello"`)
	})

	t.Run("rejects a duplicate topic", func(t *testing.T) {
		routes := []consumer.Route{{Topic: "hello", Consumer: handler}}

		err := addRoutes(&routes, consumer.Route{Topic: "hello", Consumer: handler})

		assert.ErrorContains(t, err, `"hello"`)
	})
}
