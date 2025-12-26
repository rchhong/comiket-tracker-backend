package service

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (userService UserService) GetUserByDiscordId(discordId int64) (*models.UserWithMetadata, *models.ComiketBackendError) {
	existingUser, err := userService.userRepository.GetUserByDiscordId(discordId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusNotFound}
		} else {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
	}

	return existingUser, nil
}

func (userService UserService) CreateUser(discordId int64, user models.User) (*models.UserWithMetadata, *models.ComiketBackendError) {
	createdUser, err := userService.userRepository.CreateUser(discordId, user)
	if err != nil {
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}
	return createdUser, nil
}

func (userService UserService) UpdateUser(discordId int64, user models.User) (*models.UserWithMetadata, *models.ComiketBackendError) {
	updatedUser, err := userService.userRepository.UpdateUser(discordId, user)
	if err != nil {
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}
	return updatedUser, nil
}

func (userService UserService) UpsertUser(discordId int64, user models.User) (*models.UserWithMetadata, *models.ComiketBackendError) {
	_, err := userService.GetUserByDiscordId(discordId)
	if err == nil {
		return userService.UpdateUser(discordId, user)
	}

	var statusError models.ComiketBackendError
	if errors.As(err, &statusError) {
		if statusError.StatusCode == http.StatusNotFound {
			return userService.CreateUser(discordId, user)
		}
	}
	return nil, err
}

func (userService UserService) DeleteUser(discordId int64) *models.ComiketBackendError {
	existingUser, err := userService.GetUserByDiscordId(discordId)
	if existingUser != nil {
		err := userService.userRepository.DeleteUser(discordId)
		if err != nil {
			return &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
	}

	return err
}
