package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rchhong/comiket-backend/internal/models"
)

type DoujinRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NewDoujinRepositoryPostgres(dbpool *pgxpool.Pool) *DoujinRepositoryPostgres {
	return &DoujinRepositoryPostgres{
		dbpool: dbpool,
	}
}

func (doujinRepository *DoujinRepositoryPostgres) CreateDoujin(doujin models.Doujin) (*models.DoujinWithMetadata, error) {
	var newDoujinWithMetadata models.DoujinWithMetadata
	row, err := doujinRepository.dbpool.Query(context.Background(), `
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
		doujin.PriceInYen,
		doujin.PriceInUsd,
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

func (doujinRepository *DoujinRepositoryPostgres) GetDoujinByMelonbooksId(melonbooksId int) (*models.DoujinWithMetadata, error) {
	var doujin models.DoujinWithMetadata

	row, err := doujinRepository.dbpool.Query(context.Background(), `
		SELECT * FROM doujins WHERE melonbooks_id = $1
	`, melonbooksId)
	if err != nil {
		return nil, err
	}

	doujin, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.DoujinWithMetadata])
	if err != nil {
		return nil, err
	}

	return &doujin, nil
}

func (doujinRepository *DoujinRepositoryPostgres) UpdateDoujin(melonbooksId int, updatedDoujin models.Doujin) (*models.DoujinWithMetadata, error) {
	var doujin models.DoujinWithMetadata

	row, err := doujinRepository.dbpool.Query(context.Background(), `
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
		updatedDoujin.PriceInYen,
		updatedDoujin.PriceInUsd,
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

func (doujinRepository DoujinRepositoryPostgres) DeleteDoujin(melonbooksId int) error {
	_, err := doujinRepository.dbpool.Query(context.Background(), `
		DELETE FROM doujins
		WHERE melonbooks_id = $1
	`, melonbooksId)

	return err
}
