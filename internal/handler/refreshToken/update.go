package refreshTokenHandler

import (
	"globe/internal/repository/dtos"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) Update(ctx echo.Context) error {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Missing authorization header"})
	}
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Bearer" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid authorization header"})
	}
	request := &dtos.UpdateAuthTokensRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	tokens, err := h.service.Update(
		ctx.Request().Context(),
		request.RefreshToken,
		&splitAuthHeader[1],
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"refreshToken": *tokens.RefreshToken,
		"accessToken" : *tokens.AccessToken,
	})
} 