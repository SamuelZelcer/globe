package productHandler

import (
	productService "globe/internal/service/product"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
}

type handler struct {
	service productService.Service
}

func Init(service productService.Service) Handler {
	return &handler{service: service}
}