package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/service"
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

		melonbooksId, err := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		discordId, err := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		var user models.User
		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		reservation, err := reservationController.reservationService.CreateReservation(int(melonbooksId), discordId, user)
		if err != nil {
			switch e := err.(type) {
			case models.Error:
				w.WriteHeader(e.Status())
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return

		}

		w.WriteHeader(http.StatusAccepted)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(reservation)

	})

	mux.HandleFunc("GET /users/{discordId}/reservations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, err := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		reservations, err := reservationController.reservationService.GetAllReservationsForUser(discordId)
		if err != nil {
			switch e := err.(type) {
			case models.Error:
				w.WriteHeader(e.Status())
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(reservations)

	})

	mux.HandleFunc("DELETE /users/{discordId}/reservations/{melonbooksId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, err := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		melonbooksId, err := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		err = reservationController.reservationService.DeleteReservation(int(melonbooksId), discordId)
		if err != nil {
			switch e := err.(type) {
			case models.Error:
				w.WriteHeader(e.Status())
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return

		}

		w.WriteHeader(http.StatusAccepted)
	})
}
