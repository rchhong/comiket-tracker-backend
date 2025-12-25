package service

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/repositories"
)

type ReservationService struct {
	reservationRepository repositories.ReservationRepository
	userService           *UserService
	doujinService         *DoujinService
}

func NewReservationService(reservationRepository repositories.ReservationRepository, userService *UserService, doujinService *DoujinService) *ReservationService {
	return &ReservationService{
		reservationRepository: reservationRepository,
		userService:           userService,
		doujinService:         doujinService,
	}
}

func (reservationService *ReservationService) CreateReservation(melonbooksId int, discordId int64, user models.User) (*models.ReservationWithMetadata, error) {
	// Create user, doujin if they don't exist yet
	_, err := reservationService.userService.UpsertUser(discordId, user)
	if err != nil {
		return nil, models.StatusError{Err: err, StatusCode: http.StatusInternalServerError}
	}
	_, err = reservationService.doujinService.UpsertDoujin(melonbooksId)
	if err != nil {
		return nil, models.StatusError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	existingReservation, err := reservationService.reservationRepository.GetReservationByMelonbooksIdDiscordId(melonbooksId, discordId)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return reservationService.reservationRepository.CreateReservation(melonbooksId, discordId)
		}
		return nil, models.StatusError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return existingReservation, nil
}

func (reservationService *ReservationService) GetAllReservationsForUser(discordId int64) ([]models.DoujinWithMetadata, error) {
	return reservationService.reservationRepository.GetAllReservationsForUser(discordId)
}

func (reservationService *ReservationService) DeleteReservation(melonbooksId int, discordId int64) error {
	return reservationService.reservationRepository.DeleteReservation(melonbooksId, discordId)
}
