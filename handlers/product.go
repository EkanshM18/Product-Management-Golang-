package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"product-management/models"
)

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func CreateProduct(c *gin.Context) {
	var product models.Product
	if c.BindJSON(&product) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	db, err := models.InitializeDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()
	if db.Create(&product).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	rdb.Del(context.Background(), productCacheKey(product.ID))
	c.JSON(http.StatusCreated, product)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()
	cacheKey := productCacheKey(id)
	if cachedProduct, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		var product models.Product
		json.Unmarshal([]byte(cachedProduct), &product)
		logrus.Info("Cache hit for the product ID:", id)
		c.JSON(http.StatusOK, product)
		return
	}
	db, err := models.InitializeDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()
	var product models.Product
	if db.First(&product, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	productJSON, _ := json.Marshal(product)
	rdb.Set(ctx, cacheKey, productJSON, 5*time.Minute)
	c.JSON(http.StatusOK, product)
}

func GetAllProducts(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "products:all"
	if cachedProducts, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		var products []models.Product
		json.Unmarshal([]byte(cachedProducts), &products)
		c.JSON(http.StatusOK, products)
		return
	}
	db, err := models.InitializeDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()
	var products []models.Product
	if db.Find(&products).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	productsJSON, _ := json.Marshal(products)
	rdb.Set(ctx, cacheKey, productsJSON, 5*time.Minute)
	c.JSON(http.StatusOK, products)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if c.BindJSON(&product) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	db, err := models.InitializeDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()
	var existingProduct models.Product
	if db.First(&existingProduct, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	existingProduct.ProductName = product.ProductName
	existingProduct.ProductDescription = product.ProductDescription
	existingProduct.ProductPrice = product.ProductPrice
	if db.Save(&existingProduct).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	rdb.Del(context.Background(), productCacheKey(id))
	rdb.Del(context.Background(), "products:all")
	c.JSON(http.StatusOK, existingProduct)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	db, err := models.InitializeDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()
	if db.Delete(&models.Product{}, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	rdb.Del(context.Background(), productCacheKey(id))
	rdb.Del(context.Background(), "products:all")
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func productCacheKey(id interface{}) string {
    switch v := id.(type) {
    case string:
        return "product:" + v
    case uint:
        return "product:" + strconv.Itoa(int(v))  // Convert uint to string
    default:
        return "product:unknown"  // Default case in case of an unsupported type
    }
}

