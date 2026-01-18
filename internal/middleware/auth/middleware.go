package auth

import (
	"slices"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/consts"
	"github.com/anonychun/bibit/internal/current"
	repositoryAdmin "github.com/anonychun/bibit/internal/repository/admin"
	repositoryAdminSession "github.com/anonychun/bibit/internal/repository/admin_session"
	repositoryUser "github.com/anonychun/bibit/internal/repository/user"
	repositoryUserSession "github.com/anonychun/bibit/internal/repository/user_session"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewMiddleware)
}

type Middleware interface {
	AuthenticateAdmin(next echo.HandlerFunc) echo.HandlerFunc
	AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc
}

type MiddlewareImpl struct {
	adminRepository        repositoryAdmin.Repository
	adminSessionRepository repositoryAdminSession.Repository
	userRepository         repositoryUser.Repository
	userSessionRepository  repositoryUserSession.Repository
}

func NewMiddleware(i do.Injector) (Middleware, error) {
	return &MiddlewareImpl{
		adminRepository:        do.MustInvoke[repositoryAdmin.Repository](i),
		adminSessionRepository: do.MustInvoke[repositoryAdminSession.Repository](i),
		userRepository:         do.MustInvoke[repositoryUser.Repository](i),
		userSessionRepository:  do.MustInvoke[repositoryUserSession.Repository](i),
	}, nil
}

func (m *MiddlewareImpl) AuthenticateAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bypassedPaths := []string{
			"/api/v1/admin/auth/signin",
		}

		if slices.Contains(bypassedPaths, c.Request().URL.Path) {
			return next(c)
		}

		cookie, err := c.Cookie(consts.CookieAdminSession)
		if err != nil {
			return consts.ErrUnauthorized
		}

		adminSession, err := m.adminSessionRepository.FindByToken(c.Request().Context(), cookie.Value)
		if err != nil {
			return consts.ErrUnauthorized
		}

		admin, err := m.adminRepository.FindById(c.Request().Context(), adminSession.AdminId.String())
		if err != nil {
			return consts.ErrUnauthorized
		}

		ctx := current.SetAdmin(c.Request().Context(), admin)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (m *MiddlewareImpl) AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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
