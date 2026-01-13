package productHandler

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
	request := &dtos.UpdateProductRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	tokens, product, err := h.service.Update(ctx.Request().Context(), request, &splitAuthHeader[1]); 
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if request.RefreshToken == nil {
		return ctx.JSON(http.StatusOK, map[string]string{
			"name": *product.Name,
			"price": *product.Price,
			"description": *product.Description,
		})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
			"name": *product.Name,
			"price": *product.Price,
			"description": *product.Description,
			"refreshToken": *tokens.RefreshToken,
			"accessToken": *tokens.AccessToken,
		})
}