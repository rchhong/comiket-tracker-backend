package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/repositories"
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

func (reservationService *ReservationService) CreateReservation(ctx context.Context, melonbooksId int, discordId int64) (*models.ReservationWithMetadata, *models.ComiketBackendError) {
	// Create user, doujin if they don't exist yet
	_, err := reservationService.userService.GetUserByDiscordId(ctx, discordId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &models.ComiketBackendError{Err: fmt.Errorf("User with discordId %d not found.", discordId), StatusCode: http.StatusBadRequest}
	}
	_, err = reservationService.doujinService.GetDoujinByMelonbooksId(ctx, melonbooksId)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &models.ComiketBackendError{Err: fmt.Errorf("Doujin with melonbooksId %d not found.", melonbooksId), StatusCode: http.StatusBadRequest}
	}

	// This is a no-op, since nothing really needs to be updated about the reservation itself
	existingReservation, repositoryErr := reservationService.reservationRepository.GetReservationByMelonbooksIdDiscordId(ctx, melonbooksId, discordId)
	if existingReservation == nil {
		if errors.Is(pgx.ErrNoRows, repositoryErr) {
			newReservation, err := reservationService.reservationRepository.CreateReservation(ctx, melonbooksId, discordId)
			if err != nil {
				return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
			}
			return newReservation, nil
		}
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return existingReservation, nil
}

func (reservationService *ReservationService) GetAllReservationsForUser(ctx context.Context, discordId int64) ([]models.DoujinWithMetadata, *models.ComiketBackendError) {
	reservations, err := reservationService.reservationRepository.GetAllReservationsForUser(ctx, discordId)
	if err != nil {
		return nil, &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return reservations, nil
}

func (reservationService *ReservationService) DeleteReservation(ctx context.Context, melonbooksId int, discordId int64) *models.ComiketBackendError {
	err := reservationService.reservationRepository.DeleteReservation(ctx, melonbooksId, discordId)
	if err != nil {
		return &models.ComiketBackendError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return nil
}
