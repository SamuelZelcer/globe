package userUpdateHandler

import (
	"globe/internal/repository/dtos"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *updateHandler) NewEmailVerification(ctx echo.Context) error {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Missing authorization header"})
	}
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Bearer" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid authorization header"})
	}
	request := &dtos.VerificationRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	tokens, err := h.service.NewEmailVerification(ctx.Request().Context(), request, splitAuthHeader[1])
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if tokens.RefreshToken == "" {
		return ctx.JSON(http.StatusOK, map[string]string{"accessToken": tokens.AccessToken})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"refreshToken": tokens.RefreshToken,
		"accessToken": tokens.AccessToken,
	})
}