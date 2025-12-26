package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/controllers/dto"
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
		prefix:      "/users",
	}
}

func (userController UserController) RegisterUserController(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("GET %s/{discordId}", userController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		user, err := userController.userService.GetUserByDiscordId(discordId)
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return

		}

		w.WriteHeader(http.StatusCreated)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(user)

	})

	mux.HandleFunc(fmt.Sprintf("PUT %s/{discordId}", userController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		var responseBody models.User
		parseErr = json.NewDecoder(r.Body).Decode(&responseBody)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		user, err := userController.userService.UpsertUser(discordId, responseBody)
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return

		}

		w.WriteHeader(http.StatusAccepted)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(user)

	})

}
