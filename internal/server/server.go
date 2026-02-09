package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	middlewareAuth "github.com/anonychun/bibit/internal/middleware/auth"
	usecaseApiV1AdminAuth "github.com/anonychun/bibit/internal/usecase/api/v1/admin/auth"
	usecaseApiV1AppAuth "github.com/anonychun/bibit/internal/usecase/api/v1/app/auth"
	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewServer)
}

type IServer interface {
	Start(ctx context.Context) error
}

type Server struct {
	echo   *echo.Echo
	server *http.Server

	authMiddleware middlewareAuth.IMiddleware

	apiV1AdminAuthHandler usecaseApiV1AdminAuth.IHandler
	apiV1AppAuthHandler   usecaseApiV1AppAuth.IHandler
}

var _ IServer = (*Server)(nil)

func NewServer(i do.Injector) (*Server, error) {
	cfg := do.MustInvoke[*config.Config](i)

	e := echo.New()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: e,
	}

	return &Server{
		echo:   e,
		server: srv,

		authMiddleware: do.MustInvoke[*middlewareAuth.Middleware](i),

		apiV1AdminAuthHandler: do.MustInvoke[*usecaseApiV1AdminAuth.Handler](i),
		apiV1AppAuthHandler:   do.MustInvoke[*usecaseApiV1AppAuth.Handler](i),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	err := s.routes()
	if err != nil {
		return err
	}

	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
