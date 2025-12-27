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

func (userController UserController) getUserByDiscordId(r *http.Request) (int, any, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	user, err := userController.userService.GetUserByDiscordId(discordId)
	if err != nil {
		return err.Status(), nil, err
	}

	return http.StatusOK, user, nil

}

func (userController UserController) upsertUser(r *http.Request) (int, any, error) {
	discordId, parseErr := strconv.ParseInt(r.PathValue("discordId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	var responseBody models.User
	parseErr = json.NewDecoder(r.Body).Decode(&responseBody)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	user, err := userController.userService.UpsertUser(discordId, responseBody)
	if err != nil {
		return err.Status(), nil, err
	}

	return http.StatusAccepted, user, nil
}

func (userController UserController) RegisterUserController(mux *http.ServeMux) {
	userPath := fmt.Sprintf("%s/{discordId}", userController.prefix)
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, userPath, userController.getUserByDiscordId)
	utils.RegisterMethodToHTTPServer(mux, http.MethodPut, userPath, userController.upsertUser)

}
