package api

import (
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4/middleware"
	"product-catalog-service/internal/entity"
	"product-catalog-service/internal/service"
	"strconv"
)

type ProductHandler interface {
	ReleaseProductStock(c echo.Context) error
	GetProductStock(c echo.Context) error
	ReserveProductStock(c echo.Context) error

	RegisterRoutes(e *echo.Echo)
}

type productHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return &productHandler{
		ProductService: productService,
	}
}

func (ph *productHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/product/:id/stock", ph.GetProductStock)    // Get product stock by ID
	e.POST("/product/reserve", ph.ReserveProductStock) // Reserve product stock
	e.POST("/product/release", ph.ReleaseProductStock) // Release product stock
}

// GetProductStock retrieves the stock information for a specific product by its ID.
// product/{id}/stock
func (ph *productHandler) GetProductStock(c echo.Context) error {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid product ID"})

	}
	productStock, err := ph.ProductService.GetProductStock(productID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to retrieve product stock"})
	}

	return c.JSON(200, map[string]int{"stock": productStock})
}

// ReserveProductStock reserves a specified quantity of stock for a product.
// product/reserve
func (ph *productHandler) ReserveProductStock(c echo.Context) error {
	var request entity.StockReservation

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request format"})
	}

	isSuccess, err := ph.ProductService.ReserveProductStock(request.ProductID, request.Quantity)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to reserve product stock"})
	} else if !isSuccess {
		return c.JSON(400, map[string]string{"error": "Insufficient stock available"})
	}

	return c.JSON(200, map[string]string{"message": "Product stock reserved successfully"})
}

// ReleaseProductStock releases a specified quantity of stock for a product.
// product/release
func (ph *productHandler) ReleaseProductStock(c echo.Context) error {
	var request entity.StockReservation

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request format"})
	}

	isSuccess, err := ph.ProductService.ReleaseProductStock(request.ProductID, request.Quantity)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to release product stock"})
	} else if !isSuccess {
		return c.JSON(400, map[string]string{"error": "Failed to release product stock"})
	}

	return c.JSON(200, map[string]string{"message": "Product stock released successfully"})
}
