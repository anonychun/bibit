package main

import (
	"context"
	"log"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/server"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

func main() {
	cmd := &cli.Command{
		Name:  "server",
		Usage: "Manage the HTTP and gRPC servers",
	}

	cmd.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the servers",
			Action: func(ctx context.Context, c *cli.Command) error {
				httpSrv := do.MustInvoke[*server.HttpServer](bootstrap.Injector)
				grpcSrv := do.MustInvoke[*server.GrpcServer](bootstrap.Injector)

				g, ctx := errgroup.WithContext(ctx)
				g.Go(func() error { return httpSrv.Start(ctx) })
				g.Go(func() error { return grpcSrv.Start(ctx) })
				return g.Wait()
			},
		},
	}

	err := bootstrap.RunCommand(context.Background(), cmd)
	if err != nil {
		log.Fatalln("Failed to run command:", err)
	}
}
