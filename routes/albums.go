package routes

import (
	"database/sql"
	"fmt"
	"models"
	"net/http"
	"strconv"

	"repositories"

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
		fmt.Println("routes > albums > getAlbums > env not accessible")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}
	albums, err := albumRepo.GetAlbums()
	if err != nil {
		fmt.Printf("routes > albums > getAlbums > failed to get albums: error: \n%s\n", err.Error())
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
		fmt.Println("routes > albums > getAlbums > env not accessible")
		c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}
	albumID, err := albumRepo.AddAlbum(newAlbum)
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
		fmt.Println("routes > albums > getAlbums > env not accessible")
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

	album, err := albumRepo.GetAlbumByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := fmt.Sprintf("no album found with id %d", id)
			c.IndentedJSON(http.StatusNotFound, models.ErrResponse{ErrorMessage: errMsg})
		} else {
			fmt.Printf("routes > albums > getAlbums > failed to get album with id %d: error: \n%s\n", id, err.Error())
			c.IndentedJSON(http.StatusInternalServerError, models.ErrResponseForHttpStatus(http.StatusInternalServerError))
		}
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}
