package main

import (
	"product-api/database"
	"product-api/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database.Connect()

	routes.SetupProductRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
