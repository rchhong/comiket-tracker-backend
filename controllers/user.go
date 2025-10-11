package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/models"
	"github.com/rchhong/comiket-backend/service"
)

type UserController struct {
	userService        *service.UserService
	reservationService *service.ReservationService
	prefix             string
}

func NewUserController(userService *service.UserService, reservationService *service.ReservationService) *UserController {
	return &UserController{
		userService:        userService,
		reservationService: reservationService,
		prefix:             "/users",
	}
}

func (userController UserController) RegisterUserController(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("GET %s/{discordId}", userController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, err := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		user, err := userController.userService.GetUserByDiscordId(discordId)
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
		json.NewEncoder(w).Encode(user)

	})

	mux.HandleFunc(fmt.Sprintf("PUT %s/{discordId}", userController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, err := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		var responseBody models.User
		err = json.NewDecoder(r.Body).Decode(&responseBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		user, err := userController.userService.UpsertUser(discordId, responseBody)
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
		json.NewEncoder(w).Encode(user)

	})

	mux.HandleFunc(fmt.Sprintf("GET %s/{discordId}/reservations", userController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, err := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		reservations, err := userController.reservationService.GetAllReservationsForUser(discordId)
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
		json.NewEncoder(w).Encode(reservations)

	})

	mux.HandleFunc(fmt.Sprintf("DELETE %s/{discordId}/reservations/{melonbooksId}", userController.prefix), func(w http.ResponseWriter, r *http.Request) {
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

		err = userController.reservationService.DeleteReservation(int(melonbooksId), discordId)
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
