package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/internal/controllers/utils"
	"github.com/rchhong/comiket-backend/internal/service"
)

type DoujinController struct {
	doujinService *service.DoujinService
	prefix        string
}

func NewDoujinController(doujinService *service.DoujinService) *DoujinController {
	return &DoujinController{
		doujinService: doujinService,
		prefix:        "/doujins",
	}
}

func (doujinController DoujinController) getDoujinByMelonbooksId(r *http.Request) (int, any, error) {
	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	doujin, err := doujinController.doujinService.GetDoujinByMelonbooksId(int(melonbooksId))
	if err != nil {
		return err.Status(), nil, err
	}

	return http.StatusCreated, doujin, nil
}

func (doujinController DoujinController) upsertDoujin(r *http.Request) (int, any, error) {
	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return http.StatusBadRequest, nil, parseErr
	}

	doujin, err := doujinController.doujinService.UpsertDoujin(int(melonbooksId))
	if err != nil {
		return err.Status(), nil, err
	}

	return http.StatusAccepted, doujin, nil
}

func (doujinController DoujinController) RegisterDoujinController(mux *http.ServeMux) {
	doujinPath := fmt.Sprintf("%s/{melonbooksId}", doujinController.prefix)
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, doujinPath, doujinController.getDoujinByMelonbooksId)
	utils.RegisterMethodToHTTPServer(mux, http.MethodPut, doujinPath, doujinController.upsertDoujin)

}
