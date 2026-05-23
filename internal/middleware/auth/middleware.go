package auth

import (
	"slices"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/consts"
	"github.com/anonychun/bibit/internal/current"
	repositoryUser "github.com/anonychun/bibit/internal/repository/user"
	repositoryUserSession "github.com/anonychun/bibit/internal/repository/user_session"
	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewMiddleware)
}

type IMiddleware interface {
	AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc
}

type Middleware struct {
	userRepository        repositoryUser.IRepository
	userSessionRepository repositoryUserSession.IRepository
}

var _ IMiddleware = (*Middleware)(nil)

func NewMiddleware(i do.Injector) (*Middleware, error) {
	return &Middleware{
		userRepository:        do.MustInvoke[*repositoryUser.Repository](i),
		userSessionRepository: do.MustInvoke[*repositoryUserSession.Repository](i),
	}, nil
}

func (m *Middleware) AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		bypassedPaths := []string{
			"/api/v1/app/auth/signup",
			"/api/v1/app/auth/signin",
		}

		if slices.Contains(bypassedPaths, c.Request().URL.Path) {
			return next(c)
		}

		cookie, err := c.Cookie(consts.CookieUserSession)
		if err != nil {
			return consts.ErrUnauthorized
		}

		userSession, err := m.userSessionRepository.FindByToken(c.Request().Context(), cookie.Value)
		if err != nil {
			return consts.ErrUnauthorized
		}

		user, err := m.userRepository.FindById(c.Request().Context(), userSession.UserId.String())
		if err != nil {
			return consts.ErrUnauthorized
		}

		ctx := current.SetUser(c.Request().Context(), user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
