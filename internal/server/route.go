package server

import (
	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/public"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func namespace(e *echo.Group, path string, f func(e *echo.Group)) {
	f(e.Group(path))
}

func (s *Server) routes() error {
	s.echo.HTTPErrorHandler = api.HttpErrorHandler
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.RequestID())
	// s.echo.Use(middleware.RequestLogger())

	apiRouter := s.echo.Group("/api")
	namespace(apiRouter, "/v1", func(e *echo.Group) {
		namespace(e, "/admin", func(e *echo.Group) {
			e.Use(s.authMiddleware.AuthenticateAdmin)

			e.POST("/auth/signin", s.apiV1AdminAuthHandler.SignIn)
			e.POST("/auth/signout", s.apiV1AdminAuthHandler.SignOut)
			e.GET("/auth/me", s.apiV1AdminAuthHandler.Me)
		})

		namespace(e, "/app", func(e *echo.Group) {
			e.Use(s.authMiddleware.AuthenticateUser)

			e.POST("/auth/signup", s.apiV1AppAuthHandler.SignUp)
			e.POST("/auth/signin", s.apiV1AppAuthHandler.SignIn)
			e.POST("/auth/signout", s.apiV1AppAuthHandler.SignOut)
			e.GET("/auth/me", s.apiV1AppAuthHandler.Me)
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
