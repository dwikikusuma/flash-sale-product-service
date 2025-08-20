package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"product-catalog-service/infrastructure/log"
	"product-catalog-service/internal/api"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/resource"
	"product-catalog-service/internal/service"
	infrastructure2 "product-catalog-service/middleware"
	"product-catalog-service/routes"
	"time"
)

func main() {
	log.InitLogger()

	redisClient := resource.InitRedis()

	cacheRepo := repository.NewCacheRepository(redisClient)
	productRepo := repository.NewProductRepository(cacheRepo)
	productService := service.NewProductService(productRepo)
	productHandler := api.NewProductHandler(productService)

	e := echo.New()
	e.Use(middleware.RateLimiterWithConfig(infrastructure2.GetRateLimiter()))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeout(10 * time.Second))
	e.Use(echojwt.JWT([]byte("secret")))

	routes.SetupRoutes(e, productHandler)

	e.Logger.Fatal(e.Start(":8081"))
}
