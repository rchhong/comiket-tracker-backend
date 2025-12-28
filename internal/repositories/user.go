package repositories

import (
	"context"

	"github.com/rchhong/comiket-backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, discordId int64, user models.User) (*models.UserWithMetadata, error)
	GetUserByDiscordId(ctx context.Context, discordId int64) (*models.UserWithMetadata, error)
	UpdateUser(ctx context.Context, discordId int64, updatedUser models.User) (*models.UserWithMetadata, error)
	DeleteUser(ctx context.Context, discordId int64) error
}
