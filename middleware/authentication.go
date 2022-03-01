package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("TokenAuthMiddleware ran")
		c.Next()
	}
}
