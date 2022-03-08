package routes

import (
	"errors"
	"gin-learning/controllers"
	"gin-learning/models"
	"gin-learning/repositories"
	"gin-learning/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type loginresponse struct {
	AuthTokenDetails    models.ClientReadableToken `json:"auth_token_details"`
	RefreshTokenDetails models.ClientReadableToken `json:"refresh_token_detials"`
}

type gettokenresponse struct {
	AuthToken        string                     `json:"auth_token"`
	AuthTokenDetails models.ClientReadableToken `json:"auth_token_details"`
}

func AddAuthRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")

	authGroup.GET("/login", login)
	authGroup.GET("/refreshtoken", refreshAuthToken)
	authGroup.POST("/getauthtoken", getAuthToken)
	authGroup.POST("/register", register)
}

func getAuthToken(c *gin.Context) {
	authTokenExpiration := time.Now().Add(time.Hour * 2)
	authTokenString, err := models.MintToken("user01", authTokenExpiration)
	if err != nil {
		log.Printf("routes > user.go > getAuthToken > failed to mint auth token")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponse{ErrorMessage: "Could not mint token"})
		return
	}

	authTokenDetails := models.ClientReadableToken{
		ExpiresAt: authTokenExpiration.Unix(),
		Roles:     []string{},
	}

	c.IndentedJSON(http.StatusOK, gettokenresponse{
		AuthToken:        authTokenString,
		AuthTokenDetails: authTokenDetails,
	})

}

func login(c *gin.Context) {
	authTokenExpiration := time.Now().Add(time.Hour * 2)
	refreshTokenExpiration := time.Now().Add(time.Hour * 48)

	authTokenString, err := models.MintToken("user01", authTokenExpiration)
	if err != nil {
		log.Println("routes > user.go > login > failed to mint auth token")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponse{ErrorMessage: "Could not mint token"})
		return
	}

	refreshTokenString, err := models.MintToken("user01", refreshTokenExpiration)
	if err != nil {
		log.Println("routes > user.go > login > failed to mint refresh token")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponse{ErrorMessage: "Could not mint token"})
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

	response := loginresponse{
		clientReadableAuthToken,
		clientReadableSessionToken,
	}

	c.SetCookie("authtoken", authTokenString, int(authTokenExpiration.Unix()), "/", "", false, true)
	c.SetCookie("refreshtoken", refreshTokenString, int(refreshTokenExpiration.Unix()), "/", "", false, true)
	c.IndentedJSON(http.StatusOK, response)
}

func register(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	env, ok := c.MustGet("env").(models.Env)
	if !ok {
		log.Println("routes > albums > postAlbum > env not accessible")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	repo := repositories.UserRepository{DBConn: *env.DB}
	controller := controllers.UserController{UserRepository: repo}

	id, err := controller.AddUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponse{ErrorMessage: err.Error()})
		return
	}

	user.ID = id
	c.IndentedJSON(http.StatusCreated, user)
}

func refreshAuthToken(c *gin.Context) {
	authTokenStr, refreshTokenStr, err := utils.GetAuthCookiesFromContext(c)
	if err != nil {
		log.Printf("routes > albums.go > refreshAuthToken > %s\n", err.Error())
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	authTokenClaims := models.TokenClaims{}
	_, err = jwt.ParseWithClaims(authTokenStr, &authTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
	})

	// if the error is that the auth token is expired, that is fine
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		log.Printf("routes > albums.go > refreshAuthToken > %s\n", err.Error())
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	refreshTokenClaims := models.TokenClaims{}
	_, err = jwt.ParseWithClaims(refreshTokenStr, &refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("DND_JWT_PRIVATE_KEY")), nil
	})
	if err != nil {
		log.Printf("routes > albums.go > refreshAuthToken > %s\n", err.Error())
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	newExpiresAt := time.Now().Add(time.Hour * 2)
	newAuthToken, err := models.MintToken(authTokenClaims.Subject, newExpiresAt)
	if err != nil {
		log.Printf("routes > albums.go > refreshAuthToken > %s\n", err.Error())
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
	c.IndentedJSON(http.StatusOK, loginresponse{
		AuthTokenDetails:    clientReadableAuthToken,
		RefreshTokenDetails: clientReadableRefreshToken,
	})
}
