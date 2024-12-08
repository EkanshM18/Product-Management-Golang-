package handlers

import (
    "testing"
    // "net/http"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "github.com/stretchr/testify/mock"
    "product-management/models"
    "context"
    "encoding/json"
)

func BenchmarkGetProductByID_CacheHit(b *testing.B) {
    rdb := new(MockRedisClient)
    ctx := context.Background()
    product := models.Product{ID: 1, UserID: 123, ProductName: "Test Product", ProductDescription: "Description", ProductPrice: 10.99}
    productJSON, _ := json.Marshal(product)

    rdb.On("Get", ctx, "product:1").Return(redis.NewStringResult(string(productJSON), nil))

    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    c.Params = gin.Params{{Key: "id", Value: "1"}}

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        GetProductByID(c)
    }
}

func BenchmarkGetProductByID_CacheMiss(b *testing.B) {
    rdb := new(MockRedisClient)
    ctx := context.Background()
    // product := models.Product{ID: 1, UserID: 123, ProductName: "Test Product", ProductDescription: "Description", ProductPrice: 10.99}

    rdb.On("Get", ctx, "product:1").Return(redis.NewStringResult("", redis.Nil))
    rdb.On("Set", ctx, "product:1", mock.Anything, mock.Anything).Return(redis.NewStatusResult("OK", nil))

    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    c.Params = gin.Params{{Key: "id", Value: "1"}}

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        GetProductByID(c)
    }
}
