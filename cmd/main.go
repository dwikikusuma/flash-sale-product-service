package main

import (
	"github.com/labstack/echo/v4"
	"product-catalog-service/internal/api"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/service"
)

func main() {
	productRepo := repository.NewProductRepository()
	productService := service.NewProductService(productRepo)
	productHandler := api.NewProductHandler(productService)

	e := echo.New()
	productHandler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
