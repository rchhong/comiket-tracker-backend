package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/internal/db"
	"github.com/rchhong/comiket-backend/internal/models"
)

type UserRepositoryPostgres struct {
	postgresDb *db.PostgresDB
}

func NewUserRepositoryPostgres(postgresDb *db.PostgresDB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		postgresDb: postgresDb,
	}
}

func (userRepository *UserRepositoryPostgres) CreateUser(ctx context.Context, discordId int64, user models.User) (*models.UserWithMetadata, error) {
	var newUserWithMetadata models.UserWithMetadata
	err := pgx.BeginFunc(ctx, userRepository.postgresDb.Dbpool, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			INSERT INTO users
				(discord_id, discord_global_name)
			VALUES
				($1, $2, $3)
			RETURNING *
			`, discordId, user.DiscordGlobalName)

		if err != nil {
			return err
		}

		newUserWithMetadata, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newUserWithMetadata, nil
}

func (userRepository *UserRepositoryPostgres) GetUserByDiscordId(ctx context.Context, discordId int64) (*models.UserWithMetadata, error) {
	var user models.UserWithMetadata

	err := pgx.BeginFunc(ctx, userRepository.postgresDb.Dbpool, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			SELECT * FROM users WHERE discord_id = $1 LIMIT 1
		`, discordId)
		if err != nil {
			return err
		}

		user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepository *UserRepositoryPostgres) GetUsers(ctx context.Context) ([]models.UserWithMetadata, error) {
	var user []models.UserWithMetadata

	err := pgx.BeginFunc(ctx, userRepository.postgresDb.Dbpool, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			SELECT * FROM users
		`)
		if err != nil {
			return err
		}

		user, err = pgx.CollectRows(row, pgx.RowToStructByName[models.UserWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepository *UserRepositoryPostgres) UpdateUser(ctx context.Context, discordId int64, updatedUser models.User) (*models.UserWithMetadata, error) {
	var user models.UserWithMetadata

	err := pgx.BeginFunc(ctx, userRepository.postgresDb.Dbpool, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			UPDATE users
			SET
				discord_global_name = $1,
				updated_at = NOW()
			WHERE discord_id = $3
			RETURNING *
		`, updatedUser.DiscordGlobalName, discordId)

		if err != nil {
			return err
		}

		user, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.UserWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepository UserRepositoryPostgres) DeleteUser(ctx context.Context, discordId int64) error {
	return pgx.BeginFunc(ctx, userRepository.postgresDb.Dbpool, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			DELETE FROM users
			WHERE discord_id = $1
		`, discordId)

		return err
	})
}
