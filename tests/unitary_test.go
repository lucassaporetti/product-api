package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-api/database"
	"product-api/handlers"
	"product-api/models"
	"product-api/settings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func setupTestDB() *gorm.DB {
	settings.InitConfig()
	cfg := settings.AppConfig

	dsn := cfg.DBUsername + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		panic("Failed to perform database migrations: " + err.Error())
	}

	testDB = db
	database.DB = db
	return db
}

func beginTransaction() *gorm.DB {
	tx := testDB.Begin()
	database.DB = tx
	return tx
}

func rollbackTransaction(tx *gorm.DB) {
	tx.Rollback()
	database.DB = testDB
}

func setupEcho() *echo.Echo {
	e := echo.New()
	return e
}

func TestCreateProduct_WhenValidRequest_ShouldReturnCreatedProduct(t *testing.T) {
	e := setupEcho()
	setupTestDB()
	tx := beginTransaction()
	defer rollbackTransaction(tx)

	reqBody := `{"name": "Test Product", "description": "This is a test product", "price": 9.99}`
	req := httptest.NewRequest(http.MethodPost, "/create_product", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.CreateProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var product models.Product
		json.Unmarshal(rec.Body.Bytes(), &product)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, "This is a test product", product.Description)
		assert.Equal(t, 9.99, product.Price)
	}
}

func TestUpdateProduct_WhenValidRequest_ShouldReturnUpdatedProduct(t *testing.T) {
	e := setupEcho()
	setupTestDB()
	tx := beginTransaction()
	defer rollbackTransaction(tx)

	product := models.Product{
		Name:        "Old Product",
		Description: "Old Description",
		Price:       19.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	database.DB.Create(&product)

	reqBody := `{"name": "Updated Product", "description": "Updated Description", "price": 29.99}`
	req := httptest.NewRequest(http.MethodPut, "/update_product/"+product.ID, bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(product.ID)

	if assert.NoError(t, handlers.UpdateProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var updatedProduct models.Product
		json.Unmarshal(rec.Body.Bytes(), &updatedProduct)
		assert.Equal(t, "Updated Product", updatedProduct.Name)
		assert.Equal(t, "Updated Description", updatedProduct.Description)
		assert.Equal(t, 29.99, updatedProduct.Price)
	}
}

func TestGetAllProducts_WhenProductsExist_ShouldReturnAllProducts(t *testing.T) {
	e := setupEcho()
	setupTestDB()
	tx := beginTransaction()
	defer rollbackTransaction(tx)

	products := []models.Product{
		{Name: "Product 1", Description: "Description 1", Price: 10.0, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Product 2", Description: "Description 2", Price: 20.0, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	database.DB.Create(&products)

	req := httptest.NewRequest(http.MethodGet, "/get_all_products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.GetAllProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var fetchedProducts []models.Product
		json.Unmarshal(rec.Body.Bytes(), &fetchedProducts)
		assert.Len(t, fetchedProducts, 2)
	}
}

func TestGetProduct_WhenValidID_ShouldReturnProduct(t *testing.T) {
	e := setupEcho()
	setupTestDB()
	tx := beginTransaction()
	defer rollbackTransaction(tx)

	product := models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       9.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	database.DB.Create(&product)

	req := httptest.NewRequest(http.MethodGet, "/get_product/"+product.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(product.ID)

	if assert.NoError(t, handlers.GetProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var fetchedProduct models.Product
		json.Unmarshal(rec.Body.Bytes(), &fetchedProduct)
		assert.Equal(t, product.Name, fetchedProduct.Name)
		assert.Equal(t, product.Description, fetchedProduct.Description)
		assert.Equal(t, product.Price, fetchedProduct.Price)
	}
}

func TestDeleteProduct_WhenValidID_ShouldDeleteProduct(t *testing.T) {
	e := setupEcho()
	setupTestDB()
	tx := beginTransaction()
	defer rollbackTransaction(tx)

	product := models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       9.99,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	database.DB.Create(&product)

	req := httptest.NewRequest(http.MethodDelete, "/delete_product/"+product.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(product.ID)

	if assert.NoError(t, handlers.DeleteProduct(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		var deletedProduct models.Product
		result := database.DB.First(&deletedProduct, "id = ?", product.ID)
		assert.Error(t, result.Error)
	}
}
