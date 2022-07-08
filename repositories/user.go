package repositories

import (
	"database/sql"
	"fmt"
	"gin-learning/models"
	"log"

	"github.com/go-sql-driver/mysql"
)

type IUserRepository interface {
	AddUser(models.User) (int, error)
	DeleteUser(int) error
	GetUserByID(int) (models.User, error)
}

type UserRepository struct {
	DBConn sql.DB
}

func (repo UserRepository) AddUser(user models.User) (int, error) {
	dbConn := repo.DBConn

	result, err := dbConn.Exec("insert into users (username, password) values (?, ?)",
		user.Username, user.Password)

	if err != nil {
		mysqlerr, _ := err.(*mysql.MySQLError)
		if mysqlerr != nil {
			if mysqlerr.Number == 1062 {
				return 0, fmt.Errorf("user already exists")
			}
		}

		log.Printf("repositories > user.go > AddUser > error: %s", err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (repo UserRepository) GetUserByID(id int) (models.User, error) {
	dbConn := repo.DBConn

	row := dbConn.QueryRow("select id, username from users where id = ?", id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo UserRepository) DeleteUser(id int) error {
	dbConn := repo.DBConn
	_, err := dbConn.Exec("delete from users where id = ?", id)

	return err
}
