package handlers

import (

	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
	"product-management/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	product := models.Product{UserID: 123, ProductName: "Test Product", ProductDescription: "Description", ProductPrice: 10.99}
	productJSON, _ := json.Marshal(product)

	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)

	c.Request = req

	// Mocking database call
	db := models.MockDB{}
	db.On("Create", mock.Anything).Return(nil)

	CreateProduct(c)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var createdProduct models.Product
	json.NewDecoder(rr.Body).Decode(&createdProduct)
	assert.Equal(t, "Test Product", createdProduct.ProductName)
}
