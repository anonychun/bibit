package main

import (
	"context"
	"log"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/server"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "server",
		Usage: "Manage the HTTP server",
	}

	cmd.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the server",
			Action: func(ctx context.Context, c *cli.Command) error {
				srv := do.MustInvoke[*server.Server](bootstrap.Injector)
				return srv.Start(ctx)
			},
		},
	}

	err := bootstrap.RunCommand(context.Background(), cmd)
	if err != nil {
		log.Fatalln("Failed to run command:", err)
	}
}
