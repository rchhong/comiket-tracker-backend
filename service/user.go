package service

import (
	"errors"
	"net/http"

	"github.com/rchhong/comiket-backend/dao"
	"github.com/rchhong/comiket-backend/models"
)

type UserService struct {
	userDAO *dao.UserDAO
}

func NewUserService(userDAO *dao.UserDAO) *UserService {
	return &UserService{
		userDAO: userDAO,
	}
}

func (userService UserService) GetUserByDiscordId(discordId int64) (*models.UserWithMetadata, error) {
	return userService.userDAO.GetUserByDiscordId(discordId)
}

func (userService UserService) CreateUser(user models.User) (*models.UserWithMetadata, error) {
	return userService.userDAO.CreateUser(user)
}

func (userService UserService) UpdateUser(discordId int64, user models.User) (*models.UserWithMetadata, error) {
	return userService.userDAO.UpdateUser(discordId, user)
}

func (userService UserService) UpsertUser(discordId int64, user models.User) (*models.UserWithMetadata, error) {
	_, err := userService.GetUserByDiscordId(discordId)
	if err == nil {
		return userService.UpdateUser(discordId, user)
	}

	var statusError models.StatusError
	if errors.As(err, &statusError) {
		if statusError.StatusCode == http.StatusNotFound {
			return userService.CreateUser(user)
		}
	}
	return nil, err
}

func (userService UserService) DeleteUser(discordId int64) error {
	_, err := userService.GetUserByDiscordId(discordId)
	if err == nil {
		return userService.userDAO.DeleteUser(discordId)
	}

	return err
}
