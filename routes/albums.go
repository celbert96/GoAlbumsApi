package routes

import (
	"database/sql"
	"fmt"
	"gin-learning/controllers"
	"gin-learning/models"
	"log"
	"net/http"
	"strconv"

	"gin-learning/repositories"

	"github.com/gin-gonic/gin"
)

func AddAlbumRoutes(rg *gin.RouterGroup) {
	albumsGroup := rg.Group("/albums")

	albumsGroup.GET("/", getAlbums)
	albumsGroup.GET("/:id", getAlbumByID)
	albumsGroup.POST("/", postAlbum)
}

func getAlbums(c *gin.Context) {
	env, ok := c.MustGet("env").(models.Env)

	if !ok {
		log.Println("routes > albums > getAlbums > env not accessible")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}
	albumController := controllers.AlbumController{AlbumRepository: albumRepo}

	albums, err := albumController.GetAlbums()
	if err != nil {
		log.Printf("routes > albums > getAlbums > failed to get albums: error: \n%s\n", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context) {
	var newAlbum models.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponseForHttpStatus(http.StatusBadRequest))
		return
	}

	if err := newAlbum.AlbumIsValid(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponse{ErrorMessage: err.Error()})
		return
	}

	env, ok := c.MustGet("env").(models.Env)
	if !ok {
		log.Println("routes > albums > postAlbum > env not accessible")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}
	albumController := controllers.AlbumController{AlbumRepository: albumRepo}

	albumID, err := albumController.AddAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	newAlbum.ID = albumID
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	env, ok := c.MustGet("env").(models.Env)

	if !ok {
		log.Println("routes > albums > getAlbumByID > env not accessible")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponse{ErrorMessage: ""})
		return
	}

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ErrResponse{ErrorMessage: "album id must be an integer"})
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}
	albumController := controllers.AlbumController{AlbumRepository: albumRepo}

	album, err := albumController.GetAlbumByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := fmt.Sprintf("no album found with id %d", id)
			c.IndentedJSON(http.StatusNotFound, models.ErrResponse{ErrorMessage: errMsg})
		} else {
			log.Printf("routes > albums > getAlbumByID > failed to get album with id %d: error: \n%s\n", id, err.Error())
			c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		}
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}
