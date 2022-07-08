package controllers

import (
	"gin-learning/models"
	"gin-learning/repositories"
)

type AlbumController struct {
	AlbumRepository repositories.IAlbumRepository
}

func (ac AlbumController) GetAlbums() ([]models.Album, error) {
	return ac.AlbumRepository.GetAlbums()
}

func (ac AlbumController) GetAlbumByID(id int) (models.Album, error) {
	return ac.AlbumRepository.GetAlbumByID(id)
}

func (ac AlbumController) AddAlbum(album models.Album) (int, error) {
	return ac.AlbumRepository.AddAlbum(album)
}
