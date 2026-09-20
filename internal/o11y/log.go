package o11y

import (
	"log/slog"
	"os"
)

func newLogger() (*slog.Logger, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	return logger, nil
}
