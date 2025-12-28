package repositories

import (
	"context"

	"github.com/rchhong/comiket-backend/internal/models"
)

type DoujinRepository interface {
	CreateDoujin(ctx context.Context, doujin models.Doujin) (*models.DoujinWithMetadata, error)
	GetDoujinByMelonbooksId(ctx context.Context, melonbooksId int) (*models.DoujinWithMetadata, error)
	UpdateDoujin(ctx context.Context, melonbooksId int, updatedDoujin models.Doujin) (*models.DoujinWithMetadata, error)
	DeleteDoujin(ctx context.Context, melonbooksId int) error
}
