package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/internal/controllers/utils"
	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/service"
)

type ReservationController struct {
	reservationService *service.ReservationService
}

func NewReservationController(reservationService *service.ReservationService) *ReservationController {
	return &ReservationController{
		reservationService: reservationService,
	}
}

func (reservationController ReservationController) getReservationsForUser(r *http.Request) (any, int, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	reservations, err := reservationController.reservationService.GetAllReservationsForUser(discordId)
	if err != nil {
		return nil, err.Status(), err
	}

	return reservations, http.StatusCreated, nil
}

func (reservationController ReservationController) upsertReservation(r *http.Request) (any, int, error) {
	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	var user models.User
	parseErr = json.NewDecoder(r.Body).Decode(&user)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	reservation, err := reservationController.reservationService.CreateReservation(int(melonbooksId), discordId, user)
	if err != nil {
		return nil, err.Status(), err

	}

	return reservation, http.StatusAccepted, nil
}

func (reservationController ReservationController) deleteReservation(r *http.Request) (any, int, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	err := reservationController.reservationService.DeleteReservation(int(melonbooksId), discordId)
	if err != nil {
		return nil, err.Status(), err

	}

	return nil, http.StatusAccepted, nil
}

func (reservationController *ReservationController) RegisterReservationController(mux *http.ServeMux) {
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, "/doujins/{melonbooksId}/reservations", reservationController.getReservationsForUser)
	utils.RegisterMethodToHTTPServer(mux, http.MethodPut, "/doujins/{melonbooksId}/reservations/{discordId}", reservationController.upsertReservation)
	utils.RegisterMethodToHTTPServer(mux, http.MethodDelete, "/doujins/{melonbooksId}/reservations/{discordId}", reservationController.deleteReservation)
}
