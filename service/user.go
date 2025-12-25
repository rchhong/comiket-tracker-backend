package service

import (
	"errors"
	"net/http"

	"github.com/rchhong/comiket-backend/repositories"
	"github.com/rchhong/comiket-backend/models"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (userService UserService) GetUserByDiscordId(discordId int64) (*models.UserWithMetadata, error) {
	return userService.userRepository.GetUserByDiscordId(discordId)
}

func (userService UserService) CreateUser(discordId int64, user models.User) (*models.UserWithMetadata, error) {
	return userService.userRepository.CreateUser(discordId, user)
}

func (userService UserService) UpdateUser(discordId int64, user models.User) (*models.UserWithMetadata, error) {
	return userService.userRepository.UpdateUser(discordId, user)
}

func (userService UserService) UpsertUser(discordId int64, user models.User) (*models.UserWithMetadata, error) {
	_, err := userService.GetUserByDiscordId(discordId)
	if err == nil {
		return userService.UpdateUser(discordId, user)
	}

	var statusError models.StatusError
	if errors.As(err, &statusError) {
		if statusError.StatusCode == http.StatusNotFound {
			return userService.CreateUser(discordId, user)
		}
	}
	return nil, err
}

func (userService UserService) DeleteUser(discordId int64) error {
	_, err := userService.GetUserByDiscordId(discordId)
	if err == nil {
		return userService.userRepository.DeleteUser(discordId)
	}

	return err
}
