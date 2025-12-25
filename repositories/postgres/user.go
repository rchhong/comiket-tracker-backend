package postgres

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rchhong/comiket-backend/models"
)

type UserRepository struct {
	dbpool *pgxpool.Pool
}

func NewUserRepository(dbpool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		dbpool: dbpool,
	}
}

func (userRepository *UserRepository) CreateUser(discordId int64, user models.User) (*models.UserWithMetadata, error) {
	var newUserWithMetadata models.UserWithMetadata
	row, err := userRepository.dbpool.Query(context.Background(), `
		INSERT INTO users 
			(discord_id, discord_name, discord_global_name) 
		VALUES
			($1, $2, $3)
		RETURNING *
		`, discordId, user.DiscordName, user.DiscordGlobalName)
	if err != nil {
		return nil, err
	}

	newUserWithMetadata, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
	if err != nil {
		return nil, err
	}

	return &newUserWithMetadata, nil
}

func (userRepository *UserRepository) GetUserByDiscordId(discordId int64) (*models.UserWithMetadata, error) {
	var user models.UserWithMetadata

	row, err := userRepository.dbpool.Query(context.Background(), `
		SELECT * FROM users WHERE discord_id = $1 LIMIT 1
	`, discordId)
	if err != nil {
		return nil, err
	}
	user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
	// TODO: move this logic to service layer, not Repository layer
	// nil, nil -> 404 error
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("[WARNING] user with discordId %d does not exist", discordId)
			return nil, models.StatusError{Err: err, StatusCode: http.StatusNotFound}
		}
		return nil, err
	}
	return &user, nil
}

func (userRepository *UserRepository) UpdateUser(discordId int64, updatedUser models.User) (*models.UserWithMetadata, error) {
	var user models.UserWithMetadata

	row, err := userRepository.dbpool.Query(context.Background(), `
		UPDATE users 
		SET 
			discord_name = $1, 
			discord_global_name = $2,
			updated_at = NOW()
		WHERE discord_id = $3
		RETURNING *
	`, updatedUser.DiscordName, updatedUser.DiscordGlobalName, discordId)
	if err != nil {
		return nil, err
	}
	user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository UserRepository) DeleteUser(discordId int64) error {
	_, err := userRepository.dbpool.Query(context.Background(), `
		DELETE FROM users 
		WHERE discord_id = $1
	`, discordId)

	return err
}
