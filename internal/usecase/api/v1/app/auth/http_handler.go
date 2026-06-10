package auth

import (
	"net/http"

	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/consts"
	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewHttpHandler)
}

type IHttpHandler interface {
	SignUp(c *echo.Context) error
	SignIn(c *echo.Context) error
	SignOut(c *echo.Context) error
	Me(c *echo.Context) error
}

type HttpHandler struct {
	usecase IUsecase
}

var _ IHttpHandler = (*HttpHandler)(nil)

func NewHttpHandler(i do.Injector) (*HttpHandler, error) {
	return &HttpHandler{
		usecase: do.MustInvoke[*Usecase](i),
	}, nil
}

func (h *HttpHandler) SignUp(c *echo.Context) error {
	req := SignUpRequest{
		IpAddress: c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	res, err := h.usecase.SignUp(c.Request().Context(), req)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     consts.CookieUserSession,
		Value:    res.Token,
		Path:     "/",
		HttpOnly: true,
	})

	return api.NewResponse(c).SendOk()
}

func (h *HttpHandler) SignIn(c *echo.Context) error {
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
		Name:     consts.CookieUserSession,
		Value:    res.Token,
		Path:     "/",
		HttpOnly: true,
	})

	return api.NewResponse(c).SendOk()
}

func (h *HttpHandler) SignOut(c *echo.Context) error {
	cookie, err := c.Cookie(consts.CookieUserSession)
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
		Name:     consts.CookieUserSession,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return c.NoContent(http.StatusNoContent)
}

func (h *HttpHandler) Me(c *echo.Context) error {
	res, err := h.usecase.Me(c.Request().Context())
	if err != nil {
		return err
	}

	return api.NewResponse(c).SetData(res).Send()
}
