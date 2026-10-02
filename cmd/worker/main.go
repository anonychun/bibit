package main

import (
	"context"
	"log"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/bootstrap/app"
	"github.com/anonychun/bibit/internal/worker"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

func main() {
	err := app.Setup()
	if err != nil {
		log.Fatalln("Failed to setup:", err)
	}

	cmd := &cli.Command{
		Name:  "worker",
		Usage: "Manage the worker process",
	}

	cmd.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the worker",
			Action: func(ctx context.Context, c *cli.Command) error {
				wrk := do.MustInvoke[*worker.Worker](bootstrap.Injector)
				return wrk.Start(ctx)
			},
		},
	}

	err = app.RunCommand(context.Background(), cmd)
	if err != nil {
		log.Fatalln("Failed to run command:", err)
	}
}
