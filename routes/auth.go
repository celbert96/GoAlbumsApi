package routes

import (
	"errors"
	"fmt"
	"gin-learning/models"
	"gin-learning/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type authresponse struct {
	AuthToken    models.ClientReadableToken `json:"auth_token"`
	RefreshToken models.ClientReadableToken `json:"refresh_token"`
}

func AddAuthRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")

	authGroup.GET("/login", login)
	authGroup.GET("/refreshtoken", refreshAuthToken)
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

func refreshAuthToken(c *gin.Context) {
	authTokenStr, refreshTokenStr, err := utils.GetAuthCookiesFromContext(c)
	if err != nil {
		fmt.Printf("routes > albums.go > refreshAuthToken > %s\n", err.Error())
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	authTokenClaims := models.TokenClaims{}
	_, err = jwt.ParseWithClaims(authTokenStr, &authTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
	})

	// if the error is that the auth token is expired, that is fine
	if !errors.Is(err, jwt.ErrTokenExpired) {
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	refreshTokenClaims := models.TokenClaims{}
	_, err = jwt.ParseWithClaims(refreshTokenStr, &refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
	})
	if err != nil {
		fmt.Printf("routes > albums.go > refreshAuthToken > %s\n", err.Error())
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	newExpiresAt := time.Now().Add(time.Hour * 2)
	newAuthToken, err := models.MintToken(authTokenClaims.Subject, newExpiresAt)
	if err != nil {
		fmt.Printf("routes > albums.go > refreshAuthToken > %s\n", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	clientReadableAuthToken := models.ClientReadableToken{
		ExpiresAt: newExpiresAt.Unix(),
		Roles:     []string{},
	}

	clientReadableRefreshToken := models.ClientReadableToken{
		ExpiresAt: refreshTokenClaims.ExpiresAt.Unix(),
		Roles:     []string{},
	}

	c.SetCookie("authtoken", newAuthToken, int(newExpiresAt.Unix()), "/", "", false, true)
	c.IndentedJSON(http.StatusOK, authresponse{
		AuthToken:    clientReadableAuthToken,
		RefreshToken: clientReadableRefreshToken,
	})
}
