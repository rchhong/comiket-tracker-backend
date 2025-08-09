package dao

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rchhong/comiket-backend/models"
)

type UserDao struct {
	dbpool *pgxpool.Pool
}

func NewUserDAO(dbpool *pgxpool.Pool) *UserDao {
	return &UserDao{
		dbpool: dbpool,
	}
}

func (userDAO *UserDao) GetUserByDiscordId(discordId int64) (models.UserWithMetadata, error) {
	var user models.UserWithMetadata
	err := userDAO.dbpool.QueryRow(context.Background(), "SELECT * FROM users WHERE discord_id = $1", discordId).Scan(&user)
	if err != nil {
		log.Printf("Error retrieving user with discordId %d: %s", discordId, err)
		return user, err
	}

	return user, nil
}
