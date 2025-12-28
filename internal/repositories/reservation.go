package repositories

import (
	"context"

	"github.com/rchhong/comiket-backend/internal/models"
)

type ReservationRepository interface {
	CreateReservation(ctx context.Context, melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error)
	GetReservationByReservationId(ctx context.Context, reservationId int64) (*models.ReservationWithMetadata, error)
	GetReservationByMelonbooksIdDiscordId(ctx context.Context, melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error)
	DeleteReservation(ctx context.Context, melonbooksId int, discordId int64) error
	GetAllReservationsForUser(ctx context.Context, discordId int64) ([]models.DoujinWithMetadata, error)
}
