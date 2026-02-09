package main

import (
	"context"
	"errors"
	"log"

	"github.com/anonychun/bibit/cmd/generate/internal"
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "generate",
		Usage: "Generate project components",
	}

	cmd.Commands = []*cli.Command{
		{
			Name:  "migration",
			Usage: "Generate a new database migration",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "name",
				},
				&cli.StringArg{
					Name:  "type",
					Value: "sql",
				},
			},
			Action: func(_ context.Context, c *cli.Command) error {
				name := c.StringArg("name")
				if name == "" {
					return errors.New("missing migration name")
				}

				return internal.GenerateMigration(name, c.StringArg("type"))
			},
		},
		{
			Name:  "usecase",
			Usage: "Generate a new usecase",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "name",
				},
			},
			Action: func(_ context.Context, c *cli.Command) error {
				name := c.StringArg("name")
				if name == "" {
					return errors.New("missing usecase name")
				}

				return internal.GenerateUsecase(name)
			},
		},
		{
			Name:  "repository",
			Usage: "Generate a new repository",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "name",
				},
			},
			Action: func(_ context.Context, c *cli.Command) error {
				name := c.StringArg("name")
				if name == "" {
					return errors.New("missing repository name")
				}

				return internal.GenerateRepository(name)
			},
		},
		{
			Name:  "entity",
			Usage: "Generate a new entity",
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "name",
				},
			},
			Action: func(_ context.Context, c *cli.Command) error {
				name := c.StringArg("name")
				if name == "" {
					return errors.New("missing entity name")
				}

				return internal.GenerateEntity(name)
			},
		},
	}

	err := bootstrap.RunCommand(context.Background(), cmd)
	if err != nil {
		log.Fatalln("Failed to run command:", err)
	}
}
