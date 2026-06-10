package server

import (
	"github.com/anonychun/bibit/public"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func namespace(e *echo.Group, path string, f func(e *echo.Group)) {
	f(e.Group(path))
}

func (s *HttpServer) routes() error {
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(s.loggerMiddleware.RequestLogger)

	apiRouter := s.echo.Group("/api")
	namespace(apiRouter, "/v1", func(e *echo.Group) {
		namespace(e, "/app", func(e *echo.Group) {
			e.Use(s.authMiddleware.AuthenticateUser)

			e.POST("/auth/signup", s.apiV1AppAuthHttpHandler.SignUp)
			e.POST("/auth/signin", s.apiV1AppAuthHttpHandler.SignIn)
			e.POST("/auth/signout", s.apiV1AppAuthHttpHandler.SignOut)
			e.GET("/auth/me", s.apiV1AppAuthHttpHandler.Me)
		})

		namespace(e, "/landing", func(e *echo.Group) {
		})
	})

	s.echo.StaticFS("/", public.PublicFs)
	s.echo.GET("/up", func(c *echo.Context) error {
		return nil
	})

	return nil
}
