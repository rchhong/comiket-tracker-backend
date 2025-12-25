package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rchhong/comiket-backend/models"
)

type ReservationRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NewReservationRepositoryPostgres(dbpool *pgxpool.Pool) *ReservationRepositoryPostgres {
	return &ReservationRepositoryPostgres{
		dbpool: dbpool,
	}
}

func (reservationRepository *ReservationRepositoryPostgres) CreateReservation(melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error) {
	var newReservation models.ReservationWithMetadata
	row, err := reservationRepository.dbpool.Query(context.Background(), `
		INSERT INTO reservations 
			(melonbooks_id, discord_id) 
		VALUES
			($1, $2)
		RETURNING *
		`, melonbooksId, discordId)
	if err != nil {
		return nil, err
	}

	newReservation, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.ReservationWithMetadata])
	if err != nil {
		return nil, err
	}

	return &newReservation, nil
}

func (reservationRepository *ReservationRepositoryPostgres) GetReservationByReservationId(reservationId int64) (*models.ReservationWithMetadata, error) {
	var reservation models.ReservationWithMetadata

	row, err := reservationRepository.dbpool.Query(context.Background(), `
		SELECT * FROM reservations WHERE reservation_id = $1 LIMIT 1
	`, reservationId)
	if err != nil {
		return nil, err
	}
	reservation, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.ReservationWithMetadata])
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (reservationRepository *ReservationRepositoryPostgres) GetReservationByMelonbooksIdDiscordId(melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error) {
	var reservation models.ReservationWithMetadata

	row, err := reservationRepository.dbpool.Query(context.Background(), `
		SELECT * FROM reservations WHERE melonbooks_id = $1 AND discord_id = $2 LIMIT 1
	`, melonbooksId, discordId)
	if err != nil {
		return nil, err
	}
	reservation, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.ReservationWithMetadata])
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (reservationRepository ReservationRepositoryPostgres) DeleteReservation(melonbooksId int, discordId int64) error {
	// TODO: should this be a no-op if the resource doesn't exist
	_, err := reservationRepository.dbpool.Query(context.Background(), `
		DELETE FROM reservations 
		WHERE melonbooks_id = $1 AND discord_id = $2
	`, melonbooksId, discordId)

	return err
}

// TODO: Create function to get all reservations for doujin
func (reservationRepository ReservationRepositoryPostgres) GetAllReservationsForUser(discordId int64) ([]models.DoujinWithMetadata, error) {
	var reservations []models.DoujinWithMetadata

	rows, err := reservationRepository.dbpool.Query(context.Background(), `
		WITH user_reservations AS (
			SELECT melonbooks_id FROM reservations WHERE discord_id = $1
		) 
		SELECT * FROM doujins WHERE melonbooks_id IN (SELECT melonbooks_id FROM user_reservations)
	`, discordId)
	if err != nil {
		return nil, err
	}
	reservations, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.DoujinWithMetadata])

	if err != nil {
		return nil, err
	}
	return reservations, nil
}
