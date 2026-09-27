# Bibit

Bibit is a Go project template for backend services. One command gives you a working HTTP server, gRPC server, background worker, PostgreSQL setup, migrations, code generators, tests, and observability, so you can start on business logic right away.

The code comes with a small sign-up/sign-in feature (`user` and `user_session`). That feature exists to show how each layer is written. Read it, copy its patterns, then delete it. See [Starter kit](#starter-kit).

## Contents

- [Requirements](#requirements)
- [Create a project](#create-a-project)
- [Commands](#commands)
- [Project layout](#project-layout)
- [How the pieces connect](#how-the-pieces-connect)
- [Building a feature, layer by layer](#building-a-feature-layer-by-layer)
  - [1. Entity](#1-entity)
  - [2. Migration](#2-migration)
  - [3. Repository](#3-repository)
  - [4. Usecase](#4-usecase)
  - [5. HTTP handler and routes](#5-http-handler-and-routes)
  - [6. gRPC handler](#6-grpc-handler)
  - [7. Middleware](#7-middleware)
  - [8. Background jobs](#8-background-jobs)
  - [9. File storage](#9-file-storage)
- [Shared packages](#shared-packages)
- [Testing](#testing)
- [Observability](#observability)
- [Deployment](#deployment)
- [Starter kit](#starter-kit)

## Requirements

- Go 1.27+
- PostgreSQL 18+ (the schema uses `uuidv7()`)
- [goreman](https://github.com/mattn/goreman) for `./bin/dev`: `go install github.com/mattn/goreman@latest`
- `protoc`, only if you write gRPC services

The scripts in `bin/` install the other tools (air, mockery, nilaway, protoc plugins) the first time you run them.

## Create a project

```bash
wget -qO- https://raw.githubusercontent.com/anonychun/bibit/refs/heads/main/new.sh | bash -s github.com/anonychun/billing-api
```

The script clones Bibit, replaces the module path `github.com/anonychun/bibit` with the name you pass, copies `.env.sample` to `.env`, and writes the result to a directory named after the last path segment (`billing-api` in the example).

Then:

```bash
cd billing-api
# fill in .env
./bin/db setup
./bin/dev
```

## Commands

Every command lives in `cmd/` and has a wrapper script in `bin/`.

| Command | What it does |
| --- | --- |
| `./bin/dev` | Runs `air`, which rebuilds on file changes and starts the processes in `Procfile` (server and worker) with goreman |
| `./bin/server start` | Starts the HTTP and gRPC servers |
| `./bin/worker start` | Starts the background job worker |
| `./bin/db create` / `drop` | Creates or drops the database named in `DB_SQL_NAME` |
| `./bin/db migrate` / `rollback` | Applies all pending migrations, or reverts the last one |
| `./bin/db seed` | Runs the seeder |
| `./bin/db setup` | `create`, then `migrate`, then `seed` |
| `./bin/db reset` | `drop`, then `setup` |
| `./bin/generate <kind> <name>` | Generates code. See the layer sections below |
| `./bin/test [args]` | Recreates the test database and runs `go test` |
| `./bin/mockery` | Regenerates `mock.go` files for every interface under `internal/` |
| `./bin/protoc` | Compiles every `.proto` file in `proto/` into `pkg/pb/` |
| `./bin/nilaway [args]` | Runs the nilaway nil-safety checker on `./...` |

## Project layout

```
cmd/
  db/             database CLI
  generate/       code generators
  server/         HTTP + gRPC servers
  worker/         background job worker
internal/
  api/            response envelope, API errors, validation errors
  bootstrap/      dependency injector and process lifecycle
  client/         clients for other systems (River job queue)
  config/         environment configuration
  consts/         shared constants and API errors
  current/        request-scoped values on context.Context
  db/             database connections, migrator, seeder
  dto/            response shapes shared across usecases
  entity/         database models
  job/            background jobs, one package per job
  lib/            small helpers
  middleware/     HTTP middleware, one package per middleware
  o11y/           logging, metrics, tracing
  repository/     data access, one package per table
  server/         HTTP server, gRPC server, routes
  storage/        file storage (S3)
  usecase/        business logic and handlers, one package per feature
  validation/     struct validation
  worker/         worker process and job registration
migrations/       SQL and Go migrations (embedded in the binary)
proto/            protobuf definitions
pkg/pb/           generated protobuf code
public/           static files served at /
```

## How the pieces connect

Bibit uses [samber/do](https://github.com/samber/do) for dependency injection. Every component follows the same shape:

```go
package user

func init() {
	do.Provide(bootstrap.Injector, NewRepository)
}

type IRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

type Repository struct {
	sqlDB dbSql.IDB
}

var _ IRepository = (*Repository)(nil)

func NewRepository(i do.Injector) (*Repository, error) {
	return &Repository{
		sqlDB: do.MustInvoke[*dbSql.PostgresDB](i),
	}, nil
}
```

The rules behind that shape:

- `init()` registers the constructor with the global `bootstrap.Injector`. The injector builds a component the first time something asks for it, so a command only opens the connections it uses.
- Each component has an `I`-prefixed interface and a struct that implements it. The `var _ I... = (*T)(nil)` line makes the compiler check that.
- A constructor fetches its dependencies with `do.MustInvoke[*Concrete](i)` and stores them as interfaces. Tests can then swap in the generated mocks.
- A package's `init()` only runs if something imports the package. The import chain starts at `cmd/`: the server imports handlers, handlers import usecases, usecases import repositories. A new component becomes available once it sits on that chain.
- Any component with a `Shutdown(ctx context.Context) error` method is shut down when the process receives SIGINT or SIGTERM. `bootstrap.RunCommand` gives shutdown 30 seconds.
- Packages that share a name get an import alias made of the layer and the path: `repositoryUser`, `repositoryUserSession`, `usecaseApiV1AppAuth`, `middlewareAuth`, `jobHello`, `dbSql`.

A request passes through the layers in this order:

```
HTTP/gRPC request
  -> middleware          (internal/middleware)
  -> handler             (internal/usecase/<feature>/http_handler.go)
  -> usecase             (internal/usecase/<feature>/usecase.go)
  -> repository          (internal/repository/<table>)
  -> PostgreSQL
```

Handlers deal with transport: binding input, cookies, status codes. Usecases hold the business rules and know nothing about HTTP. Repositories hold the queries and know nothing about business rules.

## Building a feature, layer by layer

The steps below build a hypothetical `product` feature from the bottom up.

### 1. Entity

Entities are [bun](https://bun.uptrace.dev/) models in `internal/entity`. Generate an empty file with:

```bash
./bin/generate entity product
```

Embed `Base` to get a UUIDv7 `Id`, `CreatedAt`, and `UpdatedAt`. `Base` also sets `UpdatedAt` before every bun update query. Bun maps `Product` to the `products` table and `PriceCents` to `price_cents`, so most fields need no tags.

```go
package entity

type Product struct {
	Base

	Name       string
	PriceCents int64
	OwnerId    uuid.UUID
	Owner      *User `bun:"rel:belongs-to,join:owner_id=id"`
}
```

Put behavior that belongs to one record on the entity. The starter kit does this with `User.HashPassword`, `User.ComparePassword`, and `UserSession.GenerateToken`.

### 2. Migration

Migrations use [goose](https://github.com/pressly/goose) and live in `migrations/`. They are embedded in the binary, so deployments need no extra files.

```bash
./bin/generate migration create_products        # SQL (default)
./bin/generate migration backfill_prices go     # Go
```

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	name TEXT NOT NULL,
	price_cents BIGINT NOT NULL,
	owner_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
```

Apply it with `./bin/db migrate`. `20010114000001_river.go` is an example of a Go migration; it creates River's job tables.

### 3. Repository

A repository package holds the queries for one table. Generate one with:

```bash
./bin/generate repository product
```

This creates `internal/repository/product/repository.go`. For a table with an underscore, keep the underscore in the package name (`user_session`).

Always run queries through `r.sqlDB.DB(ctx)`. It returns the transaction stored in the context if there is one, and the connection pool otherwise. That is how repository methods join a transaction without taking a transaction argument.

```go
type IRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
}

func (r *Repository) FindById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	product := &entity.Product{}
	err := r.sqlDB.DB(ctx).NewSelect().Model(product).Where("id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) Create(ctx context.Context, product *entity.Product) error {
	_, err := r.sqlDB.DB(ctx).NewInsert().Model(product).Exec(ctx)
	return err
}
```

Return database errors as they are. A missing row comes back as `sql.ErrNoRows`, and the usecase decides what that means.

#### Transactions

The root `internal/repository` package has `Repository.Transaction`. It opens a transaction, stores it in the context it passes to your function, and commits when the function returns `nil`. Any error rolls it back. Every repository call made with that inner `ctx` runs inside the transaction.

```go
err = u.repository.Transaction(ctx, func(ctx context.Context) error {
	err := u.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return u.userSessionRepository.Create(ctx, userSession)
})
```

Inject it like any other repository: `repository: do.MustInvoke[*repository.Repository](i)`.

### 4. Usecase

A usecase package holds one feature: its business logic, its request and response types, and its handlers. The package path mirrors the route, so `internal/usecase/api/v1/app/auth` serves `/api/v1/app/auth/*`.

```bash
./bin/generate usecase api/v1/app/product
```

This creates `usecase.go`, `http_handler.go`, and `dto.go`, with `Usecase` and `HttpHandler` already registered in the injector. If the feature also serves gRPC, add `grpc_handler.go` by hand as described in [6. gRPC handler](#6-grpc-handler). Delete `http_handler.go` if the feature is gRPC-only.

#### dto.go

Request and response structs for this feature. Tags control JSON binding and validation:

```go
type CreateRequest struct {
	OwnerId    uuid.UUID `json:"-"`
	Name       string    `json:"name" validate:"required" field:"name" label:"Name"`
	PriceCents int64     `json:"priceCents" validate:"required|min:1" field:"priceCents" label:"Price"`
}

type CreateResponse struct {
	Id uuid.UUID `json:"id"`
}
```

- `validate` holds [gookit/validate](https://github.com/gookit/validate) rules.
- `field` is the key the error appears under in the response.
- `label` is the name used in error messages.
- `json:"-"` marks values the handler fills in itself (IP address, user agent, current user) so clients cannot set them.

#### usecase.go

```go
type IUsecase interface {
	Create(ctx context.Context, req CreateRequest) (*CreateResponse, error)
}

type Usecase struct {
	validator         validation.IValidator
	productRepository repositoryProduct.IRepository
}

func NewUsecase(i do.Injector) (*Usecase, error) {
	return &Usecase{
		validator:         do.MustInvoke[*validation.Validator](i),
		productRepository: do.MustInvoke[*repositoryProduct.Repository](i),
	}, nil
}

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	validationErr := u.validator.Struct(&req)
	if validationErr.IsFail() {
		return nil, validationErr
	}

	product := &entity.Product{
		Name:       req.Name,
		PriceCents: req.PriceCents,
		OwnerId:    current.User(ctx).Id,
	}

	err := u.productRepository.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{Id: product.Id}, nil
}
```

`validator.Struct` returns an `api.ValidationError`, which is a `map[string][]string`. You can add your own checks to it before returning, as `auth.SignUp` does for a taken email address:

```go
if isEmailAddressExists {
	validationErr.AddError("emailAddress", consts.ErrEmailAddressAlreadyRegistered)
}
```

#### Errors

Usecases return plain Go errors. The HTTP error handler (`api.HttpErrorHandler`) turns them into responses:

| Error returned | Status | `errors` field |
| --- | --- | --- |
| `*api.Error` | its `Status` | `{"message": "..."}` when `Errors` is a string, otherwise `Errors` as given |
| `api.ValidationError` | 422 | `{"params": {"field": ["message"]}}` |
| `*echo.HTTPError` | its code | `{"message": "..."}` |
| anything else | 500 | `{"message": "Something went wrong"}` |

Declare the API errors your features share in `internal/consts/error.go`:

```go
var ErrProductNotFound = &api.Error{Status: http.StatusNotFound, Errors: "Product not found"}
```

Unknown errors never leak their message to the client.

### 5. HTTP handler and routes

HTTP runs on [Echo v5](https://echo.labstack.com/). A handler reads the request, calls the usecase, and writes the response. It holds no business logic.

```go
type IHttpHandler interface {
	Create(c *echo.Context) error
}

type HttpHandler struct {
	usecase IUsecase
}

func NewHttpHandler(i do.Injector) (*HttpHandler, error) {
	return &HttpHandler{
		usecase: do.MustInvoke[*Usecase](i),
	}, nil
}

func (h *HttpHandler) Create(c *echo.Context) error {
	req := CreateRequest{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	res, err := h.usecase.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return api.NewResponse(c).SetStatus(http.StatusCreated).SetData(res).Send()
}
```

Always pass `c.Request().Context()` to the usecase. Middleware stores the current user and the trace in that context.

#### Responses

`api.NewResponse(c)` builds the JSON envelope every endpoint returns:

```json
{ "ok": true, "meta": null, "data": { "id": "..." }, "errors": null }
```

`ok` is true for 2xx statuses. On success `errors` is null, and on failure `data` is null.

| Method | Use |
| --- | --- |
| `SetStatus(code)` | Status code, 200 by default |
| `SetData(v)` | Response body data |
| `SetMeta(v)` | Extra data such as pagination |
| `SetErrors(err)` | Maps an error as in the table above |
| `Send()` | Writes the response |
| `SendMessage(msg)` | Sends `{"message": msg}` as data |
| `SendOk()` | Sends `{"message": "ok"}` with status 200 |

#### Registering the route

1. Add the handler to `HttpServer` in `internal/server/http_server.go`:

   ```go
   type HttpServer struct {
   	...
   	apiV1AppProductHttpHandler usecaseApiV1AppProduct.IHttpHandler
   }

   func NewHttpServer(i do.Injector) (*HttpServer, error) {
   	...
   	return &HttpServer{
   		...
   		apiV1AppProductHttpHandler: do.MustInvoke[*usecaseApiV1AppProduct.HttpHandler](i),
   	}, nil
   }
   ```

2. Add the route in `internal/server/route.go`. The `namespace` helper groups routes by path prefix, and middleware added inside a namespace applies only to that group:

   ```go
   namespace(apiRouter, "/v1", func(e *echo.Group) {
   	namespace(e, "/app", func(e *echo.Group) {
   		e.Use(s.authMiddleware.AuthenticateUser)

   		e.POST("/products", s.apiV1AppProductHttpHandler.Create)
   	})
   })
   ```

`route.go` also sets up the global middleware: panic recovery, request IDs, request logging, OpenTelemetry, and CORS. The CORS allow list contains local frontend origins (`localhost:3000`, `localhost:5173`); edit it for your environments.

Built-in endpoints:

| Path | Purpose |
| --- | --- |
| `GET /up` | Health check |
| `GET /metrics` | Prometheus metrics |
| `/` | Static files from `public/`, embedded in the binary |

### 6. gRPC handler

A usecase can serve gRPC next to, or instead of, HTTP. The `health` usecase does both and is the reference.

1. Write the service in `proto/<name>/service.proto` and run `./bin/protoc`. Generated code goes to `pkg/pb/<name>`.
2. Add `grpc_handler.go` to the usecase. Embed the generated `Unimplemented...Server` and call the usecase:

   ```go
   type IGrpcHandler interface {
   	pb.ServiceServer
   }

   type GrpcHandler struct {
   	pb.UnimplementedServiceServer
   	usecase IUsecase
   }

   func (h *GrpcHandler) Up(ctx context.Context, req *pb.UpRequest) (*pb.UpResponse, error) {
   	res, err := h.usecase.Up(ctx)
   	if err != nil {
   		return nil, err
   	}

   	return &pb.UpResponse{Status: res.Status}, nil
   }
   ```

3. Register it in `registerGrpcHandlers` in `internal/server/grpc_server.go`:

   ```go
   pbProduct.RegisterServiceServer(srv, do.MustInvoke[*usecaseApiV1AppProduct.GrpcHandler](i))
   ```

Server reflection is on, so `grpcurl` and similar tools can list services without the proto files.

### 7. Middleware

Each middleware is a package in `internal/middleware` with the usual interface, struct, and constructor. A middleware method has Echo's signature:

```go
type IMiddleware interface {
	RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc
}

func (m *Middleware) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		user := current.User(c.Request().Context())
		if user == nil || !user.IsAdmin {
			return consts.ErrUnauthorized
		}

		return next(c)
	}
}
```

Add it to `HttpServer` the same way as a handler, then call `e.Use(s.adminMiddleware.RequireAdmin)` inside the namespace it should guard.

To pass data to later layers, put it in the request context with the `current` package. The auth middleware does this with the signed-in user:

```go
ctx := current.SetUser(c.Request().Context(), user)
c.SetRequest(c.Request().WithContext(ctx))
```

### 8. Background jobs

Jobs run on [River](https://riverqueue.com/), a job queue stored in PostgreSQL. `./bin/worker start` runs them.

```bash
./bin/generate job send_receipt
```

This creates `internal/job/send_receipt/job.go`. `Args` is the job payload and is stored as JSON. `Kind()` must return a unique name. `Work` does the job, and returning an error makes River retry it.

```go
type Args struct {
	OrderId uuid.UUID
}

func (Args) Kind() string {
	return "send_receipt"
}

func (j *Job) Work(ctx context.Context, job *river.Job[Args]) error {
	order, err := j.orderRepository.FindById(ctx, job.Args.OrderId)
	...
}
```

Register the job in `NewWorker` in `internal/worker/worker.go`. `addWorkers` is generic over one `Args` type, so call it once per job:

```go
err := addWorkers(riverClient.Workers(), do.MustInvoke[*jobHello.Job](i))
if err != nil {
	return nil, err
}

err = addWorkers(riverClient.Workers(), do.MustInvoke[*jobSendReceipt.Job](i))
if err != nil {
	return nil, err
}
```

To enqueue a job, inject the River client and insert the job's `Args`:

```go
riverClient: do.MustInvoke[*clientRiver.Client](i),
```

```go
_, err = u.riverClient.Client().Insert(ctx, jobSendReceipt.Args{OrderId: order.Id}, nil)
```

`Insert` uses its own connection, so it does not join a `Repository.Transaction`.

### 9. File storage

`internal/storage/s3` is a client for S3-compatible object storage such as AWS S3, Cloudflare R2, or MinIO. Set `STORAGE_S3_ENDPOINT`, `STORAGE_S3_BUCKET`, `STORAGE_S3_ACCESS_KEY_ID`, and `STORAGE_S3_SECRET_ACCESS_KEY` in `.env`. The injector builds the client on first use, so these can stay empty until a feature needs storage.

Inject it into a usecase and store it as the interface:

```go
type Usecase struct {
	s3Storage            storageS3.IStorage
	attachmentRepository repositoryAttachment.IRepository
}

func NewUsecase(i do.Injector) (*Usecase, error) {
	return &Usecase{
		s3Storage:            do.MustInvoke[*storageS3.Storage](i),
		attachmentRepository: do.MustInvoke[*repositoryAttachment.Repository](i),
	}, nil
}
```

| Method | Use |
| --- | --- |
| `PutObject` | Upload a file |
| `GetObject` | Download a file |
| `PresignGetObject` | Create a temporary download URL. The AWS SDK default lifetime is 15 minutes; pass `s3.WithPresignExpires(d)` to change it |

Each method takes the AWS SDK input struct. When `Bucket` is empty, the client fills in `STORAGE_S3_BUCKET`, so you usually only set `Key`.

#### Uploading a file

The template already has an `attachments` table and an `entity.Attachment` for file metadata. The starter kit has no attachment repository; create one with `./bin/generate repository attachment`.

The handler reads the multipart file and passes it to the usecase:

```go
func (h *HttpHandler) Upload(c *echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	res, err := h.usecase.Upload(c.Request().Context(), UploadRequest{File: fileHeader})
	if err != nil {
		return err
	}

	return api.NewResponse(c).SetStatus(http.StatusCreated).SetData(res).Send()
}
```

The usecase uploads the object, then saves its metadata:

```go
func (u *Usecase) Upload(ctx context.Context, req UploadRequest) (*UploadResponse, error) {
	file, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	attachment := entity.NewAttachmentFromFileHeader(req.File)
	_, err = u.s3Storage.PutObject(ctx, &s3.PutObjectInput{
		Key:  aws.String(attachment.ObjectName),
		Body: file,
	})
	if err != nil {
		return nil, err
	}

	err = u.attachmentRepository.Create(ctx, attachment)
	if err != nil {
		return nil, err
	}

	blueprint, err := dto.NewAttachmentBlueprint(ctx, attachment)
	if err != nil {
		return nil, err
	}

	return &UploadResponse{Attachment: blueprint}, nil
}
```

`NewAttachmentFromFileHeader` gives each file a unique `ObjectName` (a ULID plus the original extension), so uploads never overwrite each other. `NewAttachmentFromFile` does the same for an `*os.File`.

To return a file to clients, use `dto.NewAttachmentBlueprint`. It returns `{id, fileName, url}`, where `url` is a presigned download link. It returns `nil` for a `nil` attachment, so optional files need no extra check.

## Shared packages

| Package | What it gives you |
| --- | --- |
| `current` | Getters and setters for request-scoped values on `context.Context`: `Tx` (set by `Repository.Transaction`) and `User` (set by the auth middleware). Add a key and a getter/setter pair for new values |
| `validation` | `Validator.Struct(&req)` returns an `api.ValidationError` with every failing field, not only the first |
| `api` | Response envelope, `api.Error`, `api.ValidationError` |
| `consts` | Shared constants such as cookie names and API errors |
| `dto` | Response shapes several usecases return, such as `AttachmentBlueprint` (see [File storage](#9-file-storage)) |
| `lib` | `GetModuleName` and `ExtractPackageName`, used by the generators and o11y |

## Testing

```bash
./bin/test                                  # all packages
./bin/test ./internal/repository/...        # one tree
./bin/test -run TestUsecase_SignUp ./internal/usecase/api/v1/app/auth
```

`./bin/test` reads `.env`, points `DB_SQL_NAME` at `<name>_test`, drops and recreates that database, runs migrations, and then runs `go test` with any arguments you pass.

Each layer is tested differently:

- **Repositories** run against the real test database. Build the repository with the injected connection and use unique values (for example `uuid.NewString() + "@example.com"`) so tests don't collide:

  ```go
  repository := &Repository{sqlDB: do.MustInvoke[*dbSql.PostgresDB](bootstrap.Injector)}
  ```

- **Usecases** use mocks for their dependencies:

  ```go
  userRepository := repositoryUser.NewMockIRepository(t)
  usecase := &Usecase{userRepository: userRepository}

  userRepository.EXPECT().ExistsByEmailAddress(ctx, "ada@example.com").Return(false, nil).Once()
  ```

- **HTTP handlers** use `httptest` with a mocked usecase and assert on the status, JSON body, and cookies.

Mocks are generated by [mockery](https://github.com/vektra/mockery) into a `mock.go` next to each interface. Run `./bin/mockery` after adding or changing an interface. Never edit `mock.go` by hand.

## Observability

Every command calls `o11y.Setup` at startup:

- **Logs** go to stdout as JSON through `log/slog`. Use `slog.Info(...)` and friends anywhere.
- **Metrics** use OpenTelemetry with a Prometheus exporter, served at `GET /metrics`. The Echo OpenTelemetry middleware records HTTP metrics.
- **Traces** use OpenTelemetry, are sampled at 100%, and go to `OTLP_ENDPOINT` over gRPC. The service name is the Go module path.

## Deployment

The `Dockerfile` builds the `db`, `server`, and `worker` binaries into a Debian slim image. On start, `bin/docker-entrypoint` runs `db create`, `db migrate`, and `db seed`, then starts both processes from `Procfile` with goreman. To run one process per container, override the command:

```bash
docker run --env-file .env my-image /app/bin/server start
docker run --env-file .env my-image /app/bin/worker start
```

## Starter kit

These files exist to demonstrate the patterns. Delete them, or rewrite them for your domain, once you have your own feature in place:

- `internal/entity/user.go`, `internal/entity/user_session.go`
- `internal/repository/user/`, `internal/repository/user_session/`
- `internal/usecase/api/v1/app/auth/`
- `internal/middleware/auth/`
- `internal/job/hello/`, and its registration in `internal/worker/worker.go`
- `internal/consts/cookie.go` and the auth errors in `internal/consts/error.go`
- The auth handler field in `internal/server/http_server.go` and the `/api/v1/app` routes in `internal/server/route.go`
- The `users` and `user_sessions` tables in `migrations/20010114000000_init.sql`
- `current.User` in `internal/current/current.go`, if you drop the auth middleware

Keep the rest: `health` (backs `/up` and the gRPC health service), the root `repository` package (transactions), `attachment`, `storage`, and the infrastructure packages.
