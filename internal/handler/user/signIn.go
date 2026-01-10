package userHandler

import (
	"globe/internal/repository/dtos"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) SignIn(ctx echo.Context) error {
	request := &dtos.SignInRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	authenticationTokens, err := h.service.SignIn(request, ctx.Request().Context())
	if err != nil || authenticationTokens == nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"refreshToken" : *authenticationTokens.RefreshToken,
		"accessToken": *authenticationTokens.AccessToken,
	})
}