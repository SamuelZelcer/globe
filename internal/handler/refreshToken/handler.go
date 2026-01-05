package refreshTokenHandler

import (
	refreshTokenService "globe/internal/service/refreshToken"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Update(ctx echo.Context) error
}

type handler struct {
	service refreshTokenService.Service
}

func Init(service refreshTokenService.Service) Handler {
	return &handler{
		service: service,
	}
}