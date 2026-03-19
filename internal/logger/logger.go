package logger

import (
	"log/slog"
	"os"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewLogger)
}

type ILogger interface {
	Log() *slog.Logger
}

type Logger struct {
	logger *slog.Logger
}

var _ ILogger = (*Logger)(nil)

func NewLogger(i do.Injector) (*Logger, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Logger{
		logger: logger,
	}, nil
}

func (l *Logger) Log() *slog.Logger {
	return l.logger
}
