package controllers

import (
	"fmt"
	"gin-learning/models"
	"gin-learning/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserRepository repositories.IUserRepository
}

func (uc UserController) GetUserByID(id int) (models.User, error) {
	return uc.UserRepository.GetUserByID(id)
}

func (uc UserController) AddUser(user models.User) (int, error) {
	if err := user.IsValid(); err != nil {
		return 0, err
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt password")
	}

	user.Password = string(hashedPass)
	return uc.UserRepository.AddUser(user)
}

func (uc UserController) DeleteUser(id int) error {
	return uc.UserRepository.DeleteUser(id)
}
