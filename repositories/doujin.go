package repositories

import "github.com/rchhong/comiket-backend/models"

type DoujinRepository interface {
	CreateDoujin(doujin models.Doujin) (*models.DoujinWithMetadata, error)
	GetDoujinByMelonbooksId(melonbooksId int) (*models.DoujinWithMetadata, error)
	UpdateDoujin(melonbooksId int, updatedDoujin models.Doujin) (*models.DoujinWithMetadata, error)
	DeleteDoujin(melonbooksId int) error
}
