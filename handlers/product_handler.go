package handlers

import (
	"net/http"
	"product-api/database"
	"product-api/models"
	"time"

	"github.com/labstack/echo/v4"
)

// CreateProductRequest is used for Swagger documentation
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}

// CreateProduct
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product details"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Router /create_product [post]
func CreateProduct(c echo.Context) error {
	var input CreateProductRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := database.DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, product)
}

// UpdateProduct
// @Summary Update a product
// @Description Update an existing product by its ID
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body CreateProductRequest true "Updated product details"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /update_product/{id} [put]
func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var input CreateProductRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var existingProduct models.Product
	result := database.DB.First(&existingProduct, "id = ?", id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}
	existingProduct.Name = input.Name
	existingProduct.Description = input.Description
	existingProduct.Price = input.Price
	existingProduct.UpdatedAt = time.Now()
	if err := database.DB.Save(&existingProduct).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, existingProduct)
}

// GetAllProducts
// @Summary Get all products
// @Description Retrieve all products
// @Produce json
// @Success 200 {array} models.Product
// @Router /get_all_products [get]
func GetAllProducts(c echo.Context) error {
	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, products)
}

// GetProduct
// @Summary Get a product by ID
// @Description Retrieve a product by its ID
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /get_product/{id} [get]
func GetProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, product)
}

// DeleteProduct
// @Summary Delete a product
// @Description Delete a product by its ID
// @Produce json
// @Param id path string true "Product ID"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /delete_product/{id} [delete]
func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
