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
	userService *service.UserService
	prefix      string
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
		prefix:      "/user",
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

	mux.HandleFunc(fmt.Sprintf("POST %s", userController.prefix), func(w http.ResponseWriter, r *http.Request) {
		var reponseBody models.User
		err := json.NewDecoder(r.Body).Decode(&reponseBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return
		}

		user, err := userController.userService.CreateUser(reponseBody)
		if err != nil {
			switch e := err.(type) {
			case models.Error:
				w.WriteHeader(e.Status())
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
			return

		}

		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	})
}
