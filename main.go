package main

import (
	"product-api/database"
	"product-api/routes"
	"product-api/settings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	settings.InitConfig()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database.Connect(settings.AppConfig)

	routes.SetupProductRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
