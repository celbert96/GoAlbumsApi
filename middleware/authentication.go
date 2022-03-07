package middleware

import (
	"gin-learning/models"
	"gin-learning/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CookieTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authTokenStr, _, err := utils.GetAuthCookiesFromContext(c)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, models.ErrResponse{ErrorMessage: "Could not parse auth token(s)"})
			c.Abort()
			return
		}

		_, err = models.ValidateToken(authTokenStr)

		if err != nil {
			c.IndentedJSON(http.StatusForbidden, models.ErrResponseForHttpStatus(http.StatusForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}

func BearerTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authTokenStr, err := utils.GetBearerTokenFromContext(c)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, models.ErrResponse{ErrorMessage: err.Error()})
			c.Abort()
			return
		}

		_, err = models.ValidateToken(authTokenStr)

		if err != nil {
			c.IndentedJSON(http.StatusForbidden, models.ErrResponseForHttpStatus(http.StatusForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}
