package main

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"product-management/handlers"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	r := gin.Default()
	r.Use(loggingMiddleware())

	r.POST("/products", handlers.CreateProduct)
	r.GET("/products/:id", handlers.GetProductByID) 
	r.GET("/products", handlers.GetAllProducts)
	r.PUT("/products/:id", handlers.UpdateProduct) 
	r.DELETE("/products/:id", handlers.DeleteProduct)

	logrus.Info("Starting server on :8081")

	r.Run(":8081")
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("Request processed")
}}