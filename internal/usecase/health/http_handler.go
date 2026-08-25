package health

import (
	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewHttpHandler)
}

type IHttpHandler interface {
	Up(c *echo.Context) error
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

func (h *HttpHandler) Up(c *echo.Context) error {
	res, err := h.usecase.Up(c.Request().Context())
	if err != nil {
		return err
	}

	return api.NewResponse(c).SetData(res).Send()
}
