package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rchhong/comiket-backend/internal/controllers/dto"
	"github.com/rchhong/comiket-backend/internal/service"
)

type ExportController struct {
	exportService *service.ExportService
	prefix        string
}

func NewExportController(exportService *service.ExportService) *ExportController {
	return &ExportController{
		exportService: exportService,
		prefix:        "/admin",
	}
}

// TODO: consider just moving this to a python script or something like that
func (exportController ExportController) RegisterExportController(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s/export", exportController.prefix), func(w http.ResponseWriter, r *http.Request) {
		export, err := exportController.exportService.GenerateExport()
		if err != nil {
			w.WriteHeader(err.Status())
			json.NewEncoder(w).Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=export_%d.csv", time.Now().Unix()))
		csvWriter := csv.NewWriter(w)
		csvWriter.WriteAll(export)

	})

}
