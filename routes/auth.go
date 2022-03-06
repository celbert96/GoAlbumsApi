package routes

import (
	"fmt"
	"models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type authresponse struct {
	AuthToken    models.ClientReadableToken `json:"auth_token"`
	RefreshToken models.ClientReadableToken `json:"refresh_token"`
}

func AddAuthRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")

	authGroup.GET("/login", login)
}

func login(c *gin.Context) {
	authTokenExpiration := time.Now().Add(time.Hour * 2)
	refreshTokenExpiration := time.Now().Add(time.Hour * 48)

	authTokenString, err := models.MintToken("user01", authTokenExpiration)
	if err != nil {
		fmt.Printf("routes > user.go > login > failed to mint auth token")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not mint token"})
		return
	}

	refreshTokenString, err := models.MintToken("user01", refreshTokenExpiration)
	if err != nil {
		fmt.Printf("routes > user.go > login > failed to mint refresh token")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not mint token"})
		return
	}

	clientReadableAuthToken := models.ClientReadableToken{
		ExpiresAt: authTokenExpiration.Unix(),
		Roles:     []string{},
	}

	clientReadableSessionToken := models.ClientReadableToken{
		ExpiresAt: refreshTokenExpiration.Unix(),
		Roles:     []string{},
	}

	response := authresponse{
		clientReadableAuthToken,
		clientReadableSessionToken,
	}

	c.SetCookie("authtoken", authTokenString, int(authTokenExpiration.Unix()), "/", "", false, true)
	c.SetCookie("refreshtoken", refreshTokenString, int(refreshTokenExpiration.Unix()), "/", "", false, true)
	c.IndentedJSON(http.StatusOK, response)
}
