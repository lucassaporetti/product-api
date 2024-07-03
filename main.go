package main

import (
	"product-api/database"
	"product-api/routes"
	"product-api/settings"

	_ "product-api/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @PRODUCT-API
// @version 1.0
// @description RESTful API for product management - Swagger Documentation
// @host localhost:8080
func main() {
	settings.InitConfig()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database.Connect(settings.AppConfig)

	routes.SetupProductRoutes(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
