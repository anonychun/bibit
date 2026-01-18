package auth

import (
	"net/http"

	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/consts"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewHandler)
}

type Handler interface {
	SignIn(c echo.Context) error
	SignOut(c echo.Context) error
	Me(c echo.Context) error
}

type HandlerImpl struct {
	usecase Usecase
}

func NewHandler(i do.Injector) (Handler, error) {
	return &HandlerImpl{
		usecase: do.MustInvoke[Usecase](i),
	}, nil
}

func (h *HandlerImpl) SignIn(c echo.Context) error {
	req := SignInRequest{
		IpAddress: c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	res, err := h.usecase.SignIn(c.Request().Context(), req)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     consts.CookieAdminSession,
		Value:    res.Token,
		Path:     "/",
		HttpOnly: true,
	})

	return api.NewResponse(c).SendOk()
}

func (h *HandlerImpl) SignOut(c echo.Context) error {
	cookie, err := c.Cookie(consts.CookieAdminSession)
	if err != nil {
		return err
	}

	req := SignOutRequest{
		Token: cookie.Value,
	}

	err = h.usecase.SignOut(c.Request().Context(), req)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     consts.CookieAdminSession,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return c.NoContent(http.StatusNoContent)
}

func (h *HandlerImpl) Me(c echo.Context) error {
	res, err := h.usecase.Me(c.Request().Context())
	if err != nil {
		return err
	}

	return api.NewResponse(c).SetData(res).Send()
}
