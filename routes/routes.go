package routes

import (
	"product-catalog-service/internal/api"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ph api.ProductHandler) {
	e.GET("/product/:id/stock", ph.GetProductStock)    // Get product stock by ID
	e.POST("/product/reserve", ph.ReserveProductStock) // Reserve product stock
	e.POST("/product/release", ph.ReleaseProductStock) // Release product stock
	e.GET("/products", ph.GetAllProducts)
	e.POST("/product", ph.CreateProduct)
}
