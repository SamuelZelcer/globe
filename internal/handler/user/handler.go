package userHandler

import (
	userService "globe/internal/service/user"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	SignUp(ctx echo.Context) error
	Verification(ctx echo.Context) error
	GetNewCode(ctx echo.Context) error
	SendCodeAgain(ctx echo.Context) error
	SignIn(ctx echo.Context) error
	UpdateUsername(ctx echo.Context) error
	UpdateEmail(ctx echo.Context) error
	VerifyNewEmail(ctx echo.Context) error
}

type handler struct {
	service userService.Service
}

func Init(service userService.Service) Handler {
	return &handler{service: service}
}