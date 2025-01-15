package handler

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		latency := time.Since(start)

		fmt.Printf("[%s] %s %s %d %s %s\n",
			time.Now().Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			string(bodyBytes),
		)
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				fmt.Printf("Panic recovered: %v\n", err)

				c.JSON(500, gin.H{
					"error": "internal server error",
				})
			}
		}()

		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
