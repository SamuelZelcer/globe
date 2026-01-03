package userHandler

import (
	userService "globe/internal/service/user"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	SignUp(ctx echo.Context) error
	Verification(ctx echo.Context) error
}

type handler struct {
	service userService.Service
}

func Init(service userService.Service) Handler {
	return &handler{service: service}
}