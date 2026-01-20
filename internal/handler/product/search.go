package productHandler

import (
	"globe/internal/repository/dtos"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Search(ctx echo.Context) error {
	request := &dtos.SearchRequest{}
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
	}
	products, err := h.service.Search(ctx.Request().Context(), *request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if len(*products.Products) == 0 {
		return ctx.JSON(http.StatusOK, map[string]string{"result": "nothing was found"})
	}
	return ctx.JSON(http.StatusOK, *products)
}