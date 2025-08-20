package routes

import (
	"github.com/labstack/echo/v4"
	"product-catalog-service/internal/api"
)

func SetupRoutes(e *echo.Echo, ph api.ProductHandler) {
	e.GET("/product/:id/stock", ph.GetProductStock)    // Get product stock by ID
	e.POST("/product/reserve", ph.ReserveProductStock) // Reserve product stock
	e.POST("/product/release", ph.ReleaseProductStock) // Release product stock
}
