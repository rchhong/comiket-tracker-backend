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

	conn, err := userRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			INSERT INTO users
				(discord_id, discord_name, discord_global_name)
			VALUES
				($1, $2, $3)
			RETURNING *
			`, discordId, user.DiscordName, user.DiscordGlobalName)

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

	conn, err := userRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
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

func (userRepository *UserRepositoryPostgres) UpdateUser(ctx context.Context, discordId int64, updatedUser models.User) (*models.UserWithMetadata, error) {
	var user models.UserWithMetadata

	conn, err := userRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			UPDATE users
			SET
				discord_name = $1,
				discord_global_name = $2,
				updated_at = NOW()
			WHERE discord_id = $3
			RETURNING *
		`, updatedUser.DiscordName, updatedUser.DiscordGlobalName, discordId)

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
	conn, err := userRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Conn().Close(ctx)

	return pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			DELETE FROM users
			WHERE discord_id = $1
		`, discordId)

		return err
	})
}
