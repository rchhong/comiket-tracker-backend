package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/internal/db"
	"github.com/rchhong/comiket-backend/internal/models"
)

type ExportRepositoryPostgres struct {
	postgresDb db.PostgresDB
}

func NewExportRepositoryPostgres(postgresDb *db.PostgresDB) *ExportRepositoryPostgres {
	return &ExportRepositoryPostgres{
		postgresDb: *postgresDb,
	}
}

func (exportRepository ExportRepositoryPostgres) GetRawExportData(ctx context.Context) ([]models.ExportRow, error) {
	var exportRows []models.ExportRow

	conn, err := exportRepository.postgresDb.Dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Conn().Close(ctx)

	err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
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
			return err
		}

		exportRows, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.ExportRow])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return exportRows, nil
}
