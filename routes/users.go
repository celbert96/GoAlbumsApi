package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type authresponse struct {
	AuthToken    tokenclaims `json:"auth_token"`
	RefreshToken tokenclaims `json:"refresh_token"`
}

type tokenclaims struct {
	Issuer     string `json:"issuer"`
	Subscriber string `json:"subscriber"`
	Expires    int64  `json:"expires"`
	Issued     int64  `json:"issued"`
}

func AddUserRoutes(rg *gin.RouterGroup) {
	usersGroup := rg.Group("/users")

	usersGroup.GET("/login", login)
}

func login(c *gin.Context) {
	authTokenClaims := tokenclaims{
		"gin-api",
		"user01",
		time.Now().Add(time.Hour * 2).Unix(),
		time.Now().Unix(),
	}

	refreshTokenClaims := tokenclaims{
		authTokenClaims.Issuer,
		authTokenClaims.Subscriber,
		time.Now().Add(time.Hour * 48).Unix(),
		time.Now().Unix(),
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": authTokenClaims.Issuer,
		"sub": authTokenClaims.Subscriber,
		"exp": authTokenClaims.Expires,
		"iat": authTokenClaims.Issued,
	})

	authTokenString, err := authToken.SignedString([]byte(os.Getenv("DND_JWT_PRIVATE_KEY")))

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not sign token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": refreshTokenClaims.Issuer,
		"exp": refreshTokenClaims.Expires,
		"iat": refreshTokenClaims.Issued,
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("DND_JWT_PRIVATE_KEY")))

	response := authresponse{
		authTokenClaims,
		refreshTokenClaims,
	}

	c.SetCookie("authtoken", authTokenString, int(time.Now().Add(time.Hour*3).Unix()), "/", "", false, true)
	c.SetCookie("refreshtoken", refreshTokenString, int(time.Now().Add(time.Hour*48).Unix()), "/", "", false, true)
	c.IndentedJSON(http.StatusOK, response)
}
