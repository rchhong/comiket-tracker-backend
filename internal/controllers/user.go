package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/internal/controllers/utils"
	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (userController UserController) getUserByDiscordId(r *http.Request) (any, int, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	user, err := userController.userService.GetUserByDiscordId(r.Context(), discordId)
	if err != nil {
		return nil, err.Status(), err
	}

	return user, http.StatusOK, nil

}

// TODO: filters
func (userController UserController) getUsers(r *http.Request) (any, int, error) {
	users, err := userController.userService.GetUsers(r.Context())
	if err != nil {
		return nil, err.Status(), err
	}

	return users, http.StatusOK, nil
}

func (userController UserController) upsertUser(r *http.Request) (any, int, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	var responseBody models.User
	parseErr = json.NewDecoder(r.Body).Decode(&responseBody)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	user, wasCreated, err := userController.userService.UpsertUser(r.Context(), discordId, responseBody)
	if err != nil {
		return nil, err.Status(), err
	}

	if wasCreated {
		return user, http.StatusCreated, nil
	}

	return user, http.StatusOK, nil
}

func (userController UserController) RegisterUserController(mux *http.ServeMux) {
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, "/users/{discordId}", userController.getUserByDiscordId)
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, "/users", userController.getUsers)
	utils.RegisterMethodToHTTPServer(mux, http.MethodPut, "/users/{discordId}", userController.upsertUser)

}
