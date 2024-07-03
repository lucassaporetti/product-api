package handlers

import (
	"net/http"
	"product-api/database"
	"product-api/models"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	product := new(models.Product)
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	if err := database.DB.Create(product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, product)
}

func GetAllProducts(c echo.Context) error {
	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid UUID format"})
	}
	var existingProduct models.Product
	result := database.DB.First(&existingProduct, "id = ?", id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}
	if err := c.Bind(&existingProduct); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	existingProduct.UpdatedAt = time.Now()
	if err := database.DB.Save(&existingProduct).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, existingProduct)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
