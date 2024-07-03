package routes

import (
	"product-api/handlers"

	"github.com/labstack/echo/v4"
)

func SetupProductRoutes(e *echo.Echo) {
	e.POST("/create_product", handlers.CreateProduct)
	e.GET("/get_all_products", handlers.GetAllProducts)
	e.GET("/get_product/:id", handlers.GetProduct)
	e.PUT("/update_product/:id", handlers.UpdateProduct)
	e.DELETE("/delete_product/:id", handlers.DeleteProduct)
}
