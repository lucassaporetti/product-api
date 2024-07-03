package main

import (
	"product-api/database"
	"product-api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database.Connect()

	e.POST("/create_product", handlers.CreateProduct)
	e.GET("/get_all_products", handlers.GetProducts)
	e.GET("/get_product/:id", handlers.GetProduct)
	e.PUT("/update_product/:id", handlers.UpdateProduct)
	e.DELETE("/delete_product/:id", handlers.DeleteProduct)

	e.Logger.Fatal(e.Start(":8080"))
}
