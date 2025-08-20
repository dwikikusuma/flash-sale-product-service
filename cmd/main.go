package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"product-catalog-service/internal/api"
	"product-catalog-service/internal/infrastructure"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/service"
)

func main() {
	infrastructure.InitLogger()

	productRepo := repository.NewProductRepository()
	productService := service.NewProductService(productRepo)
	productHandler := api.NewProductHandler(productService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.JWT([]byte("secret")))

	productHandler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8081"))
}
