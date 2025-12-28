package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (userService UserService) GetUserByDiscordId(ctx context.Context, discordId int64) (*models.UserWithMetadata, *models.ComiketBackendError) {
	existingUser, err := userService.userRepository.GetUserByDiscordId(ctx, discordId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusNotFound}
		} else {
			return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
	}

	return existingUser, nil
}

func (userService UserService) CreateUser(ctx context.Context, discordId int64, user models.User) (*models.UserWithMetadata, *models.ComiketBackendError) {
	createdUser, err := userService.userRepository.CreateUser(ctx, discordId, user)
	if err != nil {
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}
	return createdUser, nil
}

func (userService UserService) UpdateUser(ctx context.Context, discordId int64, user models.User) (*models.UserWithMetadata, *models.ComiketBackendError) {
	updatedUser, err := userService.userRepository.UpdateUser(ctx, discordId, user)
	if err != nil {
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}
	return updatedUser, nil
}

func (userService UserService) UpsertUser(ctx context.Context, discordId int64, user models.User) (*models.UserWithMetadata, *models.ComiketBackendError) {
	_, err := userService.GetUserByDiscordId(ctx, discordId)
	if err == nil {
		return userService.UpdateUser(ctx, discordId, user)
	}

	var statusError models.ComiketBackendError
	if errors.As(err, &statusError) {
		if statusError.StatusCode == http.StatusNotFound {
			return userService.CreateUser(ctx, discordId, user)
		}
	}
	return nil, err
}

func (userService UserService) DeleteUser(ctx context.Context, discordId int64) *models.ComiketBackendError {
	existingUser, err := userService.GetUserByDiscordId(ctx, discordId)
	if existingUser != nil {
		err := userService.userRepository.DeleteUser(ctx, discordId)
		if err != nil {
			return &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
		}
	}

	return err
}
