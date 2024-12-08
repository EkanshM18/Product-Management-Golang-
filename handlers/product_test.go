package handlers

import (
	"time"
	"context"
	"testing"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"product-management/models"
	"net/http"
	"net/http/httptest"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func TestGetProductByID_CacheHit(t *testing.T) {
	rdb := new(MockRedisClient)
	ctx := context.Background()
	product := models.Product{ID: 1, UserID: 123, ProductName: "Test Product", ProductDescription: "Description", ProductPrice: 10.99}
	productJSON, _ := json.Marshal(product)

	rdb.On("Get", ctx, "product:1").Return(redis.NewStringResult(string(productJSON), nil))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	GetProductByID(c)

	rdb.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
}

func TestGetProductByID_CacheMiss(t *testing.T) {
	rdb := new(MockRedisClient)
	ctx := context.Background()
	product := models.Product{ID: 1, UserID: 123, ProductName: "Test Product", ProductDescription: "Description", ProductPrice: 10.99}

	rdb.On("Get", ctx, "product:1").Return(redis.NewStringResult("", redis.Nil))
	rdb.On("Set", ctx, "product:1", mock.Anything, mock.Anything).Return(redis.NewStatusResult("OK", nil))

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Simulate DB retrieval
	db := models.MockDB{}
	db.On("First", mock.Anything, "1").Return(&product, nil)

	GetProductByID(c)

	rdb.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
}
