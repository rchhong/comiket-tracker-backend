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

func (reservationService *ReservationService) CreateReservation(melonbooksId int, discordId int64, user models.User) (*models.ReservationWithMetadata, *models.ComiketBackendError) {
	// Create user, doujin if they don't exist yet
	_, err := reservationService.userService.UpsertUser(discordId, user)
	if err != nil {
		return nil, err
	}
	_, err = reservationService.doujinService.UpsertDoujin(melonbooksId)
	if err != nil {
		return nil, err
	}

	existingReservation, repositoryErr := reservationService.reservationRepository.GetReservationByMelonbooksIdDiscordId(melonbooksId, discordId)
	if existingReservation == nil {
		if errors.Is(pgx.ErrNoRows, repositoryErr) {
			newReservation, err := reservationService.reservationRepository.CreateReservation(melonbooksId, discordId)
			if err != nil {
				return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
			}
			return newReservation, nil
		}
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return existingReservation, nil
}

func (reservationService *ReservationService) GetAllReservationsForUser(discordId int64) ([]models.DoujinWithMetadata, *models.ComiketBackendError) {
	reservations, err := reservationService.reservationRepository.GetAllReservationsForUser(discordId)
	if err != nil {
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return reservations, nil
}

func (reservationService *ReservationService) DeleteReservation(melonbooksId int, discordId int64) *models.ComiketBackendError {
	err := reservationService.reservationRepository.DeleteReservation(melonbooksId, discordId)
	if err != nil {
		return &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return nil
}
