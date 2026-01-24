package userHandler

import (
	"globe/internal/repository/dtos"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) UpdateUsername(ctx echo.Context) error {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Missing authorization header"})
	}
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Bearer" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid authorization header"})
	}
	request := &dtos.UpdateUsernameRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	if err := h.service.UpdateUsername(ctx.Request().Context(), request, splitAuthHeader[1]); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"error": "Username was updated"})
}