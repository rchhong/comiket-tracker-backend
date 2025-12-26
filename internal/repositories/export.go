package repositories

import "github.com/rchhong/comiket-backend/internal/models"

type ExportRepository interface {
	GetRawExportData() ([]models.ExportRow, error)
}
