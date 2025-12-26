package userHandler

import (
	"globe/internal/repository/dtos"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Create(ctx echo.Context) error {
	request := &dtos.SignUpRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	token, err := h.userService.Create(request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, *token)
}