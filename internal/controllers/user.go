package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/internal/controllers/utils"
	"github.com/rchhong/comiket-backend/internal/models"
	"github.com/rchhong/comiket-backend/internal/service"
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

	user, err := userController.userService.UpsertUser(r.Context(), discordId, responseBody)
	if err != nil {
		return nil, err.Status(), err
	}

	return user, http.StatusAccepted, nil
}

func (userController UserController) RegisterUserController(mux *http.ServeMux) {
	userPath := fmt.Sprintf("%s/{discordId}", userController.prefix)
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, userPath, userController.getUserByDiscordId)
	utils.RegisterMethodToHTTPServer(mux, http.MethodPut, userPath, userController.upsertUser)

}
