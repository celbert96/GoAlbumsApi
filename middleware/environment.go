package middleware

import (
	"gin-learning/models"

	"github.com/gin-gonic/gin"
)

func EnvMiddleware(env models.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("env", env)
		c.Next()
	}
}
