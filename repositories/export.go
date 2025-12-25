package repositories

import "github.com/rchhong/comiket-backend/models"

type ExportRepository interface {
	GetRawExportData() ([]models.ExportRow, error)
}
