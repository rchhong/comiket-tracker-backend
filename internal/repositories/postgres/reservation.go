package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/internal/db"
	"github.com/rchhong/comiket-backend/internal/models"
)

type ReservationRepositoryPostgres struct {
	postgresDb *db.PostgresDB
}

func NewReservationRepositoryPostgres(postgresDb *db.PostgresDB) *ReservationRepositoryPostgres {
	return &ReservationRepositoryPostgres{
		postgresDb: postgresDb,
	}
}

func (reservationRepository *ReservationRepositoryPostgres) CreateReservation(ctx context.Context, melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error) {
	var newReservation models.ReservationWithMetadata

	conn, err := reservationRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			INSERT INTO reservations
				(melonbooks_id, discord_id)
			VALUES
				($1, $2)
			RETURNING *
			`, melonbooksId, discordId)

		if err != nil {
			return err
		}

		newReservation, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.ReservationWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newReservation, nil
}

func (reservationRepository *ReservationRepositoryPostgres) GetReservationByReservationId(ctx context.Context, reservationId int64) (*models.ReservationWithMetadata, error) {
	var reservation models.ReservationWithMetadata

	conn, err := reservationRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			SELECT * FROM reservations WHERE reservation_id = $1 LIMIT 1
		`, reservationId)
		if err != nil {
			return err
		}

		reservation, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.ReservationWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (reservationRepository *ReservationRepositoryPostgres) GetReservationByMelonbooksIdDiscordId(ctx context.Context, melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error) {
	var reservation models.ReservationWithMetadata

	conn, err := reservationRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		row, err := tx.Query(ctx, `
			SELECT * FROM reservations WHERE melonbooks_id = $1 AND discord_id = $2 LIMIT 1
		`, melonbooksId, discordId)
		if err != nil {
			return err
		}

		reservation, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.ReservationWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (reservationRepository ReservationRepositoryPostgres) DeleteReservation(ctx context.Context, melonbooksId int, discordId int64) error {
	// TODO: should this be a no-op if the resource doesn't exist
	conn, err := reservationRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Conn().Close(ctx)

	return pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			DELETE FROM reservations
			WHERE melonbooks_id = $1 AND discord_id = $2
		`, melonbooksId, discordId)

		return err
	})
}

// TODO: Create function to get all reservations for doujin
func (reservationRepository ReservationRepositoryPostgres) GetAllReservationsForUser(ctx context.Context, discordId int64) ([]models.DoujinWithMetadata, error) {
	var reservations []models.DoujinWithMetadata

	conn, err := reservationRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			WITH user_reservations AS (
				SELECT melonbooks_id FROM reservations WHERE discord_id = $1
			)
			SELECT * FROM doujins WHERE melonbooks_id IN (SELECT melonbooks_id FROM user_reservations)
		`, discordId)
		if err != nil {
			return err
		}

		reservations, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.DoujinWithMetadata])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return reservations, nil
}
