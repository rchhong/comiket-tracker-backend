package service

import (
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
