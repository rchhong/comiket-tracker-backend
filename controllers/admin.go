package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rchhong/comiket-backend/service"

	"github.com/rchhong/comiket-backend/models"
)

type AdminController struct {
	adminService *service.AdminService
	prefix       string
}

func NewAdminController(adminService *service.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
		prefix:       "/admin",
	}
}

func (adminController AdminController) RegisterAdminController(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s/export", adminController.prefix), func(w http.ResponseWriter, r *http.Request) {
		export, err := adminController.adminService.GenerateExport()
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

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=export_%d.csv", time.Now().Unix()))
		csvWriter := csv.NewWriter(w)
		csvWriter.WriteAll(export)

	})

}
