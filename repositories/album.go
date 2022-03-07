package repositories

import (
	"database/sql"
	"gin-learning/models"
)

type IAlbumRepository interface {
	GetAlbums() ([]models.Album, error)
	GetAlbumByID(id int) (models.Album, error)
	AddAlbum(albums models.Album) (int, error)
}

type AlbumRepository struct {
	DBConn sql.DB
}

func (albumRepo AlbumRepository) GetAlbums() ([]models.Album, error) {
	dbConn := albumRepo.DBConn

	results, err := dbConn.Query("select id, title, artist, price from albums")
	if err != nil {
		return nil, err
	}

	var r = []models.Album{}
	for results.Next() {
		var a models.Album

		err = results.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
		if err != nil {
			return nil, err
		}

		r = append(r, a)
	}

	return r, nil
}

func (albumRepo AlbumRepository) GetAlbumByID(id int) (models.Album, error) {
	dbConn := albumRepo.DBConn

	row := dbConn.QueryRow("select id, title, artist, price from albums where id = ?", id)

	var album models.Album
	err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)

	if err != nil {
		return album, err
	}

	return album, nil
}

func (albumRepo AlbumRepository) AddAlbum(album models.Album) (int, error) {
	dbConn := albumRepo.DBConn

	result, err := dbConn.Exec("insert into albums (title, artist, price) values (?, ?, ?)",
		album.Title, album.Artist, album.Price)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
