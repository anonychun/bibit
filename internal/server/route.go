package server

import (
	"github.com/anonychun/bibit/public"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
)

func namespace(e *echo.Group, path string, f func(e *echo.Group)) {
	f(e.Group(path))
}

func (s *HttpServer) routes() error {
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(s.loggerMiddleware.RequestLogger)

	s.echo.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
		TracerProvider: otel.GetTracerProvider(),
		MeterProvider:  otel.GetMeterProvider(),
	}))

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://127.0.0.1:3000",
			"http://localhost:3000",
			"http://127.0.0.1:5173",
			"http://localhost:5173",
		},
		AllowCredentials: true,
	}))

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

	s.echo.GET("/up", s.healthHttpHandler.Up)
	s.echo.StaticFS("/", public.PublicFs)

	s.echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return nil
}
