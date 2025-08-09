package dao

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
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

func (userDAO *UserDao) CreateUser(user models.User) (models.UserWithMetadata, error) {
	var newUserWithMetadata models.UserWithMetadata
	row, err := userDAO.dbpool.Query(context.Background(), `
		INSERT INTO users 
		(discord_id, discord_name, discord_global_name) 
		VALUES
		($1, $2, $3)
		RETURNING *
		`, user.Discord_Id, user.Discord_Name, user.Discord_Global_Name)
	if err != nil {
		return newUserWithMetadata, err
	}

	newUserWithMetadata, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
	if err != nil {
		return newUserWithMetadata, err
	}

	return newUserWithMetadata, nil

}

func (userDAO *UserDao) GetUserByDiscordId(discordId int64) (models.UserWithMetadata, error) {
	var user models.UserWithMetadata

	row, err := userDAO.dbpool.Query(context.Background(), `
		SELECT * FROM users WHERE discord_id = $1
	`, discordId)
	if err != nil {
		return user, err
	}
	user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, models.StatusError{Err: err, StatusCode: http.StatusNotFound}
		}
		return user, err
	}
	return user, nil
}
