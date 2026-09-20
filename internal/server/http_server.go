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
	usecaseApiV1AppAuth "github.com/anonychun/bibit/internal/usecase/api/v1/app/auth"
	usecaseHealth "github.com/anonychun/bibit/internal/usecase/health"
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
	echo   *echo.Echo
	server *http.Server

	authMiddleware   middlewareAuth.IMiddleware
	loggerMiddleware middlewareLogger.IMiddleware

	apiV1AppAuthHttpHandler usecaseApiV1AppAuth.IHttpHandler

	healthHttpHandler usecaseHealth.IHttpHandler
}

var _ IHttpServer = (*HttpServer)(nil)

func NewHttpServer(i do.Injector) (*HttpServer, error) {
	cfg := do.MustInvoke[*config.Config](i)

	e := echo.NewWithConfig(echo.Config{
		Logger:           slog.Default(),
		HTTPErrorHandler: api.HttpErrorHandler,
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Http.Port),
		Handler: e,
	}

	return &HttpServer{
		echo:   e,
		server: srv,

		authMiddleware:   do.MustInvoke[*middlewareAuth.Middleware](i),
		loggerMiddleware: do.MustInvoke[*middlewareLogger.Middleware](i),

		apiV1AppAuthHttpHandler: do.MustInvoke[*usecaseApiV1AppAuth.HttpHandler](i),

		healthHttpHandler: do.MustInvoke[*usecaseHealth.HttpHandler](i),
	}, nil
}

func (s *HttpServer) Start(ctx context.Context) error {
	err := s.routes()
	if err != nil {
		return err
	}

	slog.Info("starting http server", slog.String("addr", s.server.Addr))
	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	slog.Info("shutting down http server")
	return s.server.Shutdown(ctx)
}
