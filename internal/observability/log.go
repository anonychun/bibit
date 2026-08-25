package observability

import (
	"log/slog"
	"os"
)

func newLogger() (*slog.Logger, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger, nil
}
