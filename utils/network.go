package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetAuthCookiesFromContext(c *gin.Context) (string, string, error) {
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
