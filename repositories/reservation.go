package repositories

import "github.com/rchhong/comiket-backend/models"

type ReservationRepository interface {
	CreateReservation(melonbooksId int, discord int64) (*models.ReservationWithMetadata, error)
	GetReservationByReservationId(reservationId int64) (*models.ReservationWithMetadata, error)
	GetReservationByMelonbooksIdDiscordId(melonbooksId int, discordId int64) (*models.ReservationWithMetadata, error)
	DeleteReservation(melonbooksId int, discordId int64) error
	GetAllReservationsForUser(discordId int64) ([]models.DoujinWithMetadata, error)
}
