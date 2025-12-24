package dao

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rchhong/comiket-backend/models"
)

type DoujinDAO struct {
	dbpool *pgxpool.Pool
}

func NewDoujinDao(dbpool *pgxpool.Pool) *DoujinDAO {
	return &DoujinDAO{
		dbpool: dbpool,
	}
}

func (doujinDAO *DoujinDAO) CreateDoujin(doujin models.Doujin) (*models.DoujinWithMetadata, error) {
	var newDoujinWithMetadata models.DoujinWithMetadata
	row, err := doujinDAO.dbpool.Query(context.Background(), `
		INSERT INTO doujins 
		(
				melonbooks_id,
				title,
				price_in_yen,
				price_in_usd,
				is_r18,
				image_preview_url,
				url,
				circle,
				authors,
				genres,
				events
		) 
		VALUES 
		(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11
		)
		RETURNING *
		`, doujin.MelonbooksId,
		doujin.Title,
		doujin.PriceInUsd,
		doujin.PriceInYen,
		doujin.IsR18,
		doujin.ImagePreviewURL,
		doujin.URL,
		doujin.Circle,
		doujin.Authors,
		doujin.Genres,
		doujin.Events,
	)
	if err != nil {
		return nil, err
	}

	newDoujinWithMetadata, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.DoujinWithMetadata])
	if err != nil {
		return nil, err
	}

	return &newDoujinWithMetadata, nil
}

// TODO: Create method to retrieve all doujins

func (doujinDAO *DoujinDAO) GetDoujinByMelonbooksId(melonbooksId int) (*models.DoujinWithMetadata, error) {
	var doujin models.DoujinWithMetadata

	row, err := doujinDAO.dbpool.Query(context.Background(), `
		SELECT * FROM doujins WHERE melonbooks_id = $1
	`, melonbooksId)
	if err != nil {
		return nil, err
	}
	doujin, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.DoujinWithMetadata])
	// TODO: move this logic to service layer, not DAO layer
	// nil, nil -> 404 error
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.StatusError{Err: err, StatusCode: http.StatusNotFound}
		}
		return nil, err
	}
	return &doujin, nil
}

func (doujinDAO *DoujinDAO) UpdateDoujin(melonbooksId int, updatedDoujin models.Doujin) (*models.DoujinWithMetadata, error) {
	var doujin models.DoujinWithMetadata

	row, err := doujinDAO.dbpool.Query(context.Background(), `
		UPDATE doujins 
		SET
				title = $2,
				price_in_yen = $3,
				price_in_usd = $4,
				is_r18 = $5,
				image_preview_url = $6,
				url = $7,
				circle = $8,
				authors = $9,
				genres = $10,
				events = $11,
				updated_at = NOW()
		where melonbooks_id = $1 
		RETURNING *
		`, melonbooksId,
		updatedDoujin.Title,
		updatedDoujin.PriceInUsd,
		updatedDoujin.PriceInYen,
		updatedDoujin.IsR18,
		updatedDoujin.ImagePreviewURL,
		updatedDoujin.URL,
		updatedDoujin.Circle,
		updatedDoujin.Authors,
		updatedDoujin.Genres,
		updatedDoujin.Events,
	)
	if err != nil {
		return nil, err
	}

	doujin, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.DoujinWithMetadata])
	if err != nil {
		return nil, err
	}
	return &doujin, nil
}

func (doujinDAO DoujinDAO) DeleteDoujin(melonbooksId int) error {
	_, err := doujinDAO.dbpool.Query(context.Background(), `
		DELETE FROM doujins 
		WHERE melonbooks_id = $1
	`, melonbooksId)

	return err
}
