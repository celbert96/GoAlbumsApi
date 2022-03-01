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
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not connect to database"})
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}
	albums, err := albumRepo.GetAlbums()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not retrieve albums"})
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbum(c *gin.Context) {
	var newAlbum models.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album is not in a valid format"})
		return
	}

	if err := newAlbum.AlbumIsValid(); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	env, ok := c.MustGet("env").(models.Env)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not connect to database"})
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}
	albumID, err := albumRepo.AddAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	newAlbum.ID = albumID
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	env, ok := c.MustGet("env").(models.Env)

	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not connect to database"})
		return
	}

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album id must be an integer"})
		return
	}

	albumRepo := repositories.AlbumRepository{DBConn: *env.DB}

	album, err := albumRepo.GetAlbumByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := fmt.Sprintf("no album found with id %d", id)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": errMsg})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}
