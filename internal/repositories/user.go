package repositories

import "github.com/rchhong/comiket-backend/internal/models"

type UserRepository interface {
	CreateUser(discordId int64, user models.User) (*models.UserWithMetadata, error)
	GetUserByDiscordId(discordId int64) (*models.UserWithMetadata, error)
	UpdateUser(discordId int64, updatedUser models.User) (*models.UserWithMetadata, error)
	DeleteUser(discordId int64) error
}
