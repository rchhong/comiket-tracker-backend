package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rchhong/comiket-backend/internal/models"
)

type ExportRepositoryPostgres struct {
	dbpool *pgxpool.Pool
}

func NewExportRepositoryPostgres(dbpool *pgxpool.Pool) *ExportRepositoryPostgres {
	return &ExportRepositoryPostgres{
		dbpool: dbpool,
	}
}

func (exportRepository ExportRepositoryPostgres) GetRawExportData() ([]models.ExportRow, error) {
	var exportRows []models.ExportRow

	rows, err := exportRepository.dbpool.Query(context.Background(), `
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
