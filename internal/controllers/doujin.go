package controllers

import (
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/internal/controllers/utils"
	"github.com/rchhong/comiket-backend/internal/service"
)

type DoujinController struct {
	doujinService *service.DoujinService
}

func NewDoujinController(doujinService *service.DoujinService) *DoujinController {
	return &DoujinController{
		doujinService: doujinService,
	}
}

func (doujinController DoujinController) getDoujinByMelonbooksId(r *http.Request) (any, int, error) {
	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	doujin, err := doujinController.doujinService.GetDoujinByMelonbooksId(r.Context(), int(melonbooksId))
	if err != nil {
		return nil, err.Status(), err
	}

	return doujin, http.StatusOK, nil
}

func (doujinController DoujinController) upsertDoujin(r *http.Request) (any, int, error) {
	melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
	if parseErr != nil {
		return nil, http.StatusBadRequest, parseErr
	}

	doujin, wasCreated, err := doujinController.doujinService.UpsertDoujin(r.Context(), int(melonbooksId))
	if err != nil {
		return nil, err.Status(), err
	}

	if wasCreated {
		return doujin, http.StatusCreated, nil
	}

	return doujin, http.StatusOK, nil
}

func (doujinController DoujinController) RegisterDoujinController(mux *http.ServeMux) {
	utils.RegisterMethodToHTTPServer(mux, http.MethodGet, "/doujins/{melonbooksId}", doujinController.getDoujinByMelonbooksId)
	utils.RegisterMethodToHTTPServer(mux, http.MethodPut, "/doujins/{melonbooksId}", doujinController.upsertDoujin)

}
