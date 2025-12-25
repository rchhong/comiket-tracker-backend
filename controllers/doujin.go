package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/service"
)

type DoujinController struct {
	doujinService      *service.DoujinService
	reservationService *service.ReservationService
	prefix             string
}

func NewDoujinController(doujinService *service.DoujinService, reservationService *service.ReservationService) *DoujinController {
	return &DoujinController{
		doujinService:      doujinService,
		reservationService: reservationService,
		prefix:             "/doujins",
	}
}

func (doujinController DoujinController) RegisterDoujinController(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("GET %s/{melonbooksId}", doujinController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		melonbooksId, err := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		doujin, err := doujinController.doujinService.GetDoujinByMelonbooksId(int(melonbooksId))
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
		jsonEncoder.Encode(doujin)

	})

	mux.HandleFunc(fmt.Sprintf("PUT %s/{melonbooksId}", doujinController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		melonbooksId, err := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		doujin, err := doujinController.doujinService.UpsertDoujin(int(melonbooksId))
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
		jsonEncoder.Encode(doujin)

	})

	mux.HandleFunc(fmt.Sprintf("PUT %s/{melonbooksId}/reservations/{discordId}", doujinController.prefix), func(w http.ResponseWriter, r *http.Request) {
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

		reservation, err := doujinController.reservationService.CreateReservation(int(melonbooksId), discordId, user)
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

}
