package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, refreshToken, err := getAuthCookiesFromContext(c)
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Could not parse auth token(s)"})
			c.Abort()
			return
		}

		/* Validate auth token */
		_, authTokenParseErr := jwt.ParseWithClaims(authToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
		})

		/* Validate session token */
		_, refreshTokenParseErr := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
		})

		sessionTokenIsValid := (refreshTokenParseErr == nil)

		if authTokenParseErr != nil {
			fmt.Println(authTokenParseErr.Error())
			validationError, _ := authTokenParseErr.(*jwt.ValidationError)

			/* If the session token is valid, then refresh the auth token and allow the request to continue */
			if validationError.Errors == jwt.ValidationErrorExpired && sessionTokenIsValid {
				c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Need to refresh the auth token here"})
				c.Abort() // TODO: When token refresh is implemented, don't abort and dont return
			} else {
				errorMsg := getErrorMessage(err)
				c.IndentedJSON(http.StatusForbidden, gin.H{"message": errorMsg})
				c.Abort()
			}

			return
		}

		c.Next()
	}
}

func getAuthCookiesFromContext(c *gin.Context) (string, string, error) {
	/* Get auth token */
	authTokenCookie, err := c.Request.Cookie("authtoken")
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}

	/* Get session token */
	refreshTokenCookie, err := c.Request.Cookie("refreshtoken")
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}

	return authTokenCookie.Value, refreshTokenCookie.Value, nil
}

func getErrorMessage(err error) string {
	if validationError, _ := err.(*jwt.ValidationError); validationError.Errors == jwt.ValidationErrorExpired {
		return "auth token expired"
	}

	// TODO: When implementing roles, provide error message here

	return "invalid auth token"
}
