package userUpdateHandler

import (
	userService "globe/internal/service/user"

	"github.com/labstack/echo/v4"
)

type UpdateHandler interface {
	UpdateUsername(ctx echo.Context) error
	UpdateEmail(ctx echo.Context) error
	NewEmailVerification(ctx echo.Context) error
	UpdatePassword(ctx echo.Context) error
	NewPasswordVerification(ctx echo.Context) error
}

type updateHandler struct {
	service userService.Service
}

func Init(service userService.Service) UpdateHandler {
	return &updateHandler{service: service}
}