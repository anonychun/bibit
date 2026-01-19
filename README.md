# Bibit

Bibit is a tool to bootstrap new projects with a predefined structure and configuration. It helps developers quickly set up a new project with some opinionated defaults in mind.

The goal is after you run the initialization script, you will have a fully functional project structure with all the necessary files and folders in place, so you can start developing your application and focusing on business logic right away without worrying about setting up the basic structure.

## Initialize a new project

You can initialize a new project by running the following command in your terminal:

```bash
wget -qO- https://raw.githubusercontent.com/anonychun/bibit/refs/heads/main/new.sh | bash -s <project-name>
```

Replace `<project-name>` with the desired name for your new project. This command will create a new directory with the specified project name and set up the necessary files and folders.

The new directory will be created in the current working directory with the basename of the project name you provided. For example, if you run the command with `github.com/anonychun/verification-api`, a new directory named `verification-api` will be created.

## Project structure

The initialized project will have the following structure:

- **`bin`** - Scripts for various development and deployment tasks.
- **`cmd`** - Main application entry points and CLI commands.
  - **`server`** - HTTP server application.
  - **`db`** - Database management CLI.
  - **`generate`** - Code generation utilities.
- **`migrations`** - Database migration files.
- **`internal`** - Internal application code.
  - **`api`** - HTTP API utilities.
  - **`bootstrap`** - Coordination of application dependencies.
  - **`config`** - Configuration management with environment variable loading.
  - **`consts`** - Application constants.
  - **`current`** - Context utilities for request-scoped data.
  - **`db`** - Database layer.
  - **`dto`** - Data transfer objects.
  - **`entity`** - Database models and business entities.
  - **`middleware`** - HTTP middleware.
  - **`repository`** - Data access layer with database operations.
  - **`scheduler`** - Background job scheduling.
  - **`server`** - HTTP server setup and routing configuration.
  - **`service`** - Application services.
  - **`storage`** - File storage utilities.
  - **`usecase`** - Application layer with business logic (use cases and handlers).
  - **`validation`** - Input validation utilities.
  - **`worker`** - Background worker processes.

## Usage

After initializing a new project, navigate to the project directory:

```bash
cd <project-name>
```

Fill the environment variables in the `.env` file as needed.

### Development

To start the development server with hot-reloading, run:

```bash
./bin/dev
```

This will start the server and watch for file changes, automatically restarting the server when changes are detected.

### Generate code

To generate code you can use the generator command in the `cmd/generate` package.

#### Migration

To create a new database migration, run:

```bash
./bin/generate migration <migration-name> <migration-type>
```

The `<migration-type>` can be either `sql` or `go`, the default is `sql`.

#### Usecase

To generate a new usecase component, run:

```bash
./bin/generate usecase <component-name>
```

E.g., `./bin/generate usecase api/v1/user`.

#### Repository

To generate a new repository, run:

```bash
./bin/generate repository <repository-name>
```

E.g., `./bin/generate repository user`.

### Database

You can manage your database using the provided CLI commands.

#### Create database

To create the database, run:

```bash
./bin/db create
```

#### Drop database

To drop the database, run:

```bash
./bin/db drop
```

#### Migrate database

To migrate the database, run:

```bash
./bin/db migrate
```

#### Rollback database

To rollback the last applied migration, run:

```bash
./bin/db rollback
```

#### Seed database

To seed the database with initial data, run:

```bash
./bin/db seed
```

#### Setup database

To set up the database (create, migrate, and seed), run:

```bash
./bin/db setup
```

#### Reset database

To reset the database (drop, create, migrate, and seed), run:

```bash
./bin/db reset
```

### Server

To start the HTTP server, run:

```bash
./bin/server start
```

### Transaction

To execute a function within a database transaction in the use case layer, you can use the `repository.Transaction` function. Here's an example:

```go
func (u *UsecaseImpl) Delete(ctx context.Context, req DeleteRequest) error {
	exists, err := u.adminRepository.ExistsById(ctx, req.Id)
	if err != nil {
		return err
	}

	if !exists {
		return consts.ErrAdminNotFound
	}

	return repository.Transaction(ctx, func(ctx context.Context) error {
		err := u.adminSessionRepository.DeleteAllByAdminId(ctx, req.Id)
		if err != nil {
			return err
		}

		return u.adminRepository.DeleteById(ctx, req.Id)
	})
}
```

### Current

Current is a package that provides utilities for managing request-scoped data using context. It allows you to set and get values associated with the current request, such as user information or request ID.

Here's an example of how to use the `current` package:

```go
// Setting a value in the context
user := &entity.User{ID: 1, Name: "Achun"}
ctx = current.SetUser(ctx, user)

// Getting a value from the context
user := current.User(ctx)
```

## Starter kit

Bibit comes with default starter kit to help you get started quickly.

- **Admin management**: Includes admin authentication and CRUD operations.

You can explore the generated code, adapt it to your project's requirements, and use it to become familiar with the structure and conventions.
