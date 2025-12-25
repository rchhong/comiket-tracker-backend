package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rchhong/comiket-backend/models"
)

type ReservationRepository struct {
	dbpool *pgxpool.Pool
}

func NewReservationRepository(dbpool *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{
		dbpool: dbpool,
	}
}

func (reservationRepository *ReservationRepository) CreateReservation(melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error) {
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

func (reservationRepository *ReservationRepository) GetReservationByReservationId(reservationId int64) (*models.ReservationWithMetadata, error) {
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

func (reservationRepository *ReservationRepository) GetReservationByMelonbooksIdDiscordId(melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error) {
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

func (reservationRepository ReservationRepository) DeleteReservation(melonbooksId int, discordId int64) error {
	// TODO: should this be a no-op if the resource doesn't exist
	_, err := reservationRepository.dbpool.Query(context.Background(), `
		DELETE FROM reservations 
		WHERE melonbooks_id = $1 AND discord_id = $2
	`, melonbooksId, discordId)

	return err
}

// TODO: Create function to get all reservations for doujin
func (reservationRepository ReservationRepository) GetAllReservationsForUser(discordId int64) ([]models.DoujinWithMetadata, error) {
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

func (reservationRepository ReservationRepository) GetRawExportData() ([]models.ExportRow, error) {
	var exportRows []models.ExportRow

	rows, err := reservationRepository.dbpool.Query(context.Background(), `
		SELECT 
			r.melonbooks_id, 
			r.discord_id, 
			d.url,
			d.title,
			d.price_in_yen,
			d.price_in_usd,
			u.discord_name 
		FROM (reservations r 
			  LEFT JOIN doujins d on r.melonbooks_id = d.melonbooks_id 
			  LEFT JOIN users u ON u.discord_id = r.discord_id
		)
		ORDER BY r.discord_id, r.melonbooks_id;
	`)

	if err != nil {
		return nil, err
	}
	exportRows, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.ExportRow])

	if err != nil {
		return nil, err
	}
	return exportRows, nil
}
