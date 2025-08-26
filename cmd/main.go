package main

import (
	"product-catalog-service/config"
	"product-catalog-service/infrastructure/log"
	"product-catalog-service/internal/api"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/resource"
	"product-catalog-service/internal/service"
	infrastructure "product-catalog-service/middleware"
	"product-catalog-service/routes"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.InitLogger()
	appConfig := config.LoadConfig(
		config.WithConfigFolder([]string{"./files/config"}),
		config.WithConfigFile("config"),
		config.WithConfigType("yaml"),
	)

	redisClient := resource.InitRedis(appConfig)
	db := resource.InitDB(appConfig)

	cacheRepo := repository.NewCacheRepository(redisClient)
	productRepo := repository.NewProductRepository(cacheRepo, db)
	productService := service.NewProductService(productRepo)
	productHandler := api.NewProductHandler(productService)

	e := echo.New()
	e.Use(middleware.RateLimiterWithConfig(infrastructure.GetRateLimiter()))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeout(10 * time.Second))
	e.Use(echojwt.JWT([]byte(appConfig.Secret.JWTSecret)))

	routes.SetupRoutes(e, productHandler)

	e.Logger.Fatal(e.Start(":" + appConfig.App.Port))
}
