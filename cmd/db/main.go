package main

import (
	"context"
	"log"
	"os"

	"github.com/anonychun/bibit/internal/bootstrap"
	dbManager "github.com/anonychun/bibit/internal/db/manager"
	dbMigrator "github.com/anonychun/bibit/internal/db/migrator"
	dbSeeder "github.com/anonychun/bibit/internal/db/seeder"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{}

	cmd.Commands = []*cli.Command{
		{
			Name:  "migrate",
			Usage: "Apply all pending migrations",
			Action: func(ctx context.Context, c *cli.Command) error {
				migratorDB := do.MustInvoke[*dbMigrator.DB](bootstrap.Injector)
				return migratorDB.Migrate(ctx)
			},
		},
		{
			Name:  "rollback",
			Usage: "Revert the last applied migration",
			Action: func(ctx context.Context, c *cli.Command) error {
				migratorDB := do.MustInvoke[*dbMigrator.DB](bootstrap.Injector)
				return migratorDB.Rollback(ctx)
			},
		},
		{
			Name:  "create",
			Usage: "Create a new database",
			Action: func(ctx context.Context, c *cli.Command) error {
				managerDB := do.MustInvoke[*dbManager.DB](bootstrap.Injector)
				return managerDB.CreateDatabase(ctx)
			},
		},
		{
			Name:  "drop",
			Usage: "Drop the database",
			Action: func(ctx context.Context, c *cli.Command) error {
				managerDB := do.MustInvoke[*dbManager.DB](bootstrap.Injector)
				return managerDB.DropDatabase(ctx)
			},
		},
		{
			Name:  "seed",
			Usage: "Seed the database with initial data",
			Action: func(ctx context.Context, c *cli.Command) error {
				seederDB := do.MustInvoke[*dbSeeder.DB](bootstrap.Injector)
				return seederDB.Seed(ctx)
			},
		},
		{
			Name:  "setup",
			Usage: "Setup the database",
			Action: func(ctx context.Context, c *cli.Command) error {
				managerDB := do.MustInvoke[*dbManager.DB](bootstrap.Injector)
				err := managerDB.CreateDatabase(ctx)
				if err != nil {
					return err
				}

				migratorDB := do.MustInvoke[*dbMigrator.DB](bootstrap.Injector)
				err = migratorDB.Migrate(ctx)
				if err != nil {
					return err
				}

				seederDB := do.MustInvoke[*dbSeeder.DB](bootstrap.Injector)
				return seederDB.Seed(ctx)
			},
		},
		{
			Name:  "reset",
			Usage: "Reset the database",
			Action: func(ctx context.Context, c *cli.Command) error {
				managerDB := do.MustInvoke[*dbManager.DB](bootstrap.Injector)
				err := managerDB.DropDatabase(ctx)
				if err != nil {
					return err
				}

				err = managerDB.CreateDatabase(ctx)
				if err != nil {
					return err
				}

				migratorDB := do.MustInvoke[*dbMigrator.DB](bootstrap.Injector)
				err = migratorDB.Migrate(ctx)
				if err != nil {
					return err
				}

				seederDB := do.MustInvoke[*dbSeeder.DB](bootstrap.Injector)
				return seederDB.Seed(ctx)
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatalln("Failed to run command:", err)
	}
}
