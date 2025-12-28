package repositories

import (
	"context"

	"github.com/rchhong/comiket-backend/internal/models"
)

type ExportRepository interface {
	GetRawExportData(ctx context.Context) ([]models.ExportRow, error)
}
