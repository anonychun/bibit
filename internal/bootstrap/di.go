package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

var Injector = do.New()

func RunCommand(ctx context.Context, cmd *cli.Command) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return cmd.Run(ctx, os.Args)
	})

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	Injector.ShutdownWithContext(shutdownCtx)
	return g.Wait()
}
