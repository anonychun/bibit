package server

import (
	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/internal/bootstrap"
	middlewareAuth "github.com/anonychun/bibit/internal/middleware/auth"
	usecaseApiV1AdminAuth "github.com/anonychun/bibit/internal/usecase/api/v1/admin/auth"
	usecaseApiV1AppAuth "github.com/anonychun/bibit/internal/usecase/api/v1/app/auth"
	"github.com/anonychun/bibit/public"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/samber/do/v2"
)

func namespace(e *echo.Group, path string, f func(e *echo.Group)) {
	f(e.Group(path))
}

func routes(e *echo.Echo) error {
	authMiddleware := do.MustInvoke[middlewareAuth.Middleware](bootstrap.Injector)

	apiV1AdminAuthHandler := do.MustInvoke[usecaseApiV1AdminAuth.Handler](bootstrap.Injector)
	apiV1AppAuthHandler := do.MustInvoke[usecaseApiV1AppAuth.Handler](bootstrap.Injector)

	e.HTTPErrorHandler = api.HttpErrorHandler
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.Logger())

	apiRouter := e.Group("/api")
	namespace(apiRouter, "/v1", func(e *echo.Group) {
		namespace(e, "/admin", func(e *echo.Group) {
			e.Use(authMiddleware.AuthenticateAdmin)

			e.POST("/auth/signin", apiV1AdminAuthHandler.SignIn)
			e.POST("/auth/signout", apiV1AdminAuthHandler.SignOut)
			e.GET("/auth/me", apiV1AdminAuthHandler.Me)
		})

		namespace(e, "/app", func(e *echo.Group) {
			e.Use(authMiddleware.AuthenticateUser)

			e.POST("/auth/signup", apiV1AppAuthHandler.SignUp)
			e.POST("/auth/signin", apiV1AppAuthHandler.SignIn)
			e.POST("/auth/signout", apiV1AppAuthHandler.SignOut)
			e.GET("/auth/me", apiV1AppAuthHandler.Me)
		})

		namespace(e, "/landing", func(e *echo.Group) {
		})
	})

	e.StaticFS("/", public.PublicFs)
	e.GET("/up", func(c echo.Context) error {
		return nil
	})

	return nil
}
