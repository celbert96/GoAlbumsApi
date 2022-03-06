package middleware

import (
	"fmt"
	"gin-learning/models"
	"gin-learning/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authTokenStr, refreshTokenStr, err := utils.GetAuthCookiesFromContext(c)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Could not parse auth token(s)"})
			c.Abort()
			return
		}

		/* Validate auth token */
		authTokenClaims := jwt.StandardClaims{}
		_, authTokenParseErr := jwt.ParseWithClaims(authTokenStr, &authTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
		})

		/* Validate session token */
		_, refreshTokenParseErr := jwt.ParseWithClaims(refreshTokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
		})

		refreshTokenIsValid := (refreshTokenParseErr == nil)

		if authTokenParseErr != nil {
			fmt.Println(authTokenParseErr.Error())
			validationError, _ := authTokenParseErr.(*jwt.ValidationError)

			/* If the session token is valid, then refresh the auth token and allow the request to continue */
			if validationError.Errors == jwt.ValidationErrorExpired && refreshTokenIsValid {
				newExpiresAt := time.Now().Add(time.Hour * 2)
				newAuthToken, err := models.MintToken(authTokenClaims.Subject, newExpiresAt)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
					c.Abort()
					return
				} else {
					fmt.Printf("middleware > authentication.go > refreshed auth token for user %s\n", authTokenClaims.Subject)
					c.SetCookie("authtoken", newAuthToken, int(newExpiresAt.Unix()), "/", "", false, true)
				}
			} else {
				c.IndentedJSON(http.StatusForbidden, models.ErrResponseForHttpStatus(http.StatusForbidden))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
