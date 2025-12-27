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

func (reservationController ReservationController) getReservationsForUser(r *http.Request) (int, any, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	reservations, err := reservationController.reservationService.GetAllReservationsForUser(discordId)
	if err != nil {
		return err.Status(), nil, err
	}

	return http.StatusCreated, reservations, nil
}

func (reservationController ReservationController) upsertReservation(r *http.Request) (int, any, error) {
	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	var user models.User
	parseErr = json.NewDecoder(r.Body).Decode(&user)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	reservation, err := reservationController.reservationService.CreateReservation(int(melonbooksId), discordId, user)
	if err != nil {
		return err.Status(), nil, err

	}

	return http.StatusAccepted, reservation, nil
}

func (reservationController ReservationController) deleteReservation(r *http.Request) (int, any, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	err := reservationController.reservationService.DeleteReservation(int(melonbooksId), discordId)
	if err != nil {
		return err.Status(), nil, err

	}

	return http.StatusAccepted, nil, nil
}

func (reservationController *ReservationController) RegisterReservationController(mux *http.ServeMux) {
	utils.RegisterMethodToHTTPServer(mux, http.MethodPut, "/doujins/{melonbooksId}/reservations/{discordId}", reservationController.upsertReservation)
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, "/doujins/{melonbooksId}/reservations", reservationController.getReservationsForUser)
	utils.RegisterMethodToHTTPServer(mux, http.MethodDelete, "/doujins/{melonbooksId}/reservations/{discordId}", reservationController.deleteReservation)
}
