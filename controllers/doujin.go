package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rchhong/comiket-backend/controllers/dto"
	"github.com/rchhong/comiket-backend/service"
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

func (doujinController DoujinController) RegisterDoujinController(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("GET %s/{melonbooksId}", doujinController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		doujin, err := doujinController.doujinService.GetDoujinByMelonbooksId(int(melonbooksId))
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(doujin)

	})

	mux.HandleFunc(fmt.Sprintf("PUT %s/{melonbooksId}", doujinController.prefix), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		melonbooksId, parseErr := strconv.ParseInt(r.PathValue("melonbooksId"), 10, 64)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: parseErr.Error()})
			return
		}

		doujin, err := doujinController.doujinService.UpsertDoujin(int(melonbooksId))
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return
		}

		w.WriteHeader(http.StatusAccepted)
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.Encode(doujin)

	})

}
