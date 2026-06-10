package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	middlewareAuth "github.com/anonychun/bibit/internal/middleware/auth"
	middlewareLogger "github.com/anonychun/bibit/internal/middleware/logger"
	"github.com/anonychun/bibit/internal/observability"
	usecaseApiV1AppAuth "github.com/anonychun/bibit/internal/usecase/api/v1/app/auth"
	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewHttpServer)
}

type IHttpServer interface {
	Start(ctx context.Context) error
}

type HttpServer struct {
	echo          *echo.Echo
	server        *http.Server
	observability observability.IObservability

	authMiddleware   middlewareAuth.IMiddleware
	loggerMiddleware middlewareLogger.IMiddleware

	apiV1AppAuthHttpHandler usecaseApiV1AppAuth.IHttpHandler
}

var _ IHttpServer = (*HttpServer)(nil)

func NewHttpServer(i do.Injector) (*HttpServer, error) {
	cfg := do.MustInvoke[*config.Config](i)
	o11y := do.MustInvoke[*observability.Observability](i)

	e := echo.NewWithConfig(echo.Config{
		Logger:           o11y.Logger(),
		HTTPErrorHandler: api.HttpErrorHandler,
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Http.Port),
		Handler: e,
	}

	return &HttpServer{
		echo:          e,
		server:        srv,
		observability: o11y,

		authMiddleware:   do.MustInvoke[*middlewareAuth.Middleware](i),
		loggerMiddleware: do.MustInvoke[*middlewareLogger.Middleware](i),

		apiV1AppAuthHttpHandler: do.MustInvoke[*usecaseApiV1AppAuth.HttpHandler](i),
	}, nil
}

func (s *HttpServer) Start(ctx context.Context) error {
	err := s.routes()
	if err != nil {
		return err
	}

	s.observability.Logger().Info("starting http server", slog.String("addr", s.server.Addr))
	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	s.observability.Logger().Info("shutting down http server")
	return s.server.Shutdown(ctx)
}
