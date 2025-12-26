package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/internal/controllers/dto"
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

func (reservationController *ReservationController) RegisterReservationController(mux *http.ServeMux) {
	mux.HandleFunc("PUT /doujins/{melonbooksId}/reservations/{discordId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		var user models.User
		parseErr = json.NewDecoder(r.Body).Decode(&user)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		reservation, err := reservationController.reservationService.CreateReservation(int(melonbooksId), discordId, user)
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return

		}

		w.WriteHeader(http.StatusAccepted)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(reservation)

	})

	mux.HandleFunc("GET /users/{discordId}/reservations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		reservations, err := reservationController.reservationService.GetAllReservationsForUser(discordId)
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(reservations)

	})

	mux.HandleFunc("DELETE /users/{discordId}/reservations/{melonbooksId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		err := reservationController.reservationService.DeleteReservation(int(melonbooksId), discordId)
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return

		}

		w.WriteHeader(http.StatusAccepted)
	})
}
