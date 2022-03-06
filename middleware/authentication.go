package middleware

import (
	"gin-learning/models"
	"gin-learning/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authTokenStr, _, err := utils.GetAuthCookiesFromContext(c)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, models.ErrResponse{ErrorMessage: "Could not parse auth token(s)"})
			c.Abort()
			return
		}

		/* Validate auth token */
		authTokenClaims := models.TokenClaims{}
		_, authTokenParseErr := jwt.ParseWithClaims(authTokenStr, &authTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
		})

		if authTokenParseErr != nil {
			c.IndentedJSON(http.StatusForbidden, models.ErrResponseForHttpStatus(http.StatusForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}
