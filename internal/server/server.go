package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"
)

func Start(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	e := echo.New()
	err := routes(e)
	if err != nil {
		return err
	}

	cfg := do.MustInvoke[*config.Config](bootstrap.Injector)
	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", cfg.Server.Port),
		GracefulTimeout: 30 * time.Second,
	}

	err = sc.Start(ctx, e)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln("failed to start server:", err)
	}

	return nil
}
