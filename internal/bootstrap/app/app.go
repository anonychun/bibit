package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/o11y"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func Setup() error {
	err := config.Setup()
	if err != nil {
		return err
	}

	return o11y.Setup()
}

func RunCommand(ctx context.Context, cmd *cli.Command) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer stop()
		return cmd.Run(ctx, os.Args)
	})

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bootstrap.Injector.ShutdownWithContext(shutdownCtx)
	return g.Wait()
}
