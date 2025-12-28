package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rchhong/comiket-backend/internal/controllers/dto"
)

type handlerFunc func(r *http.Request) (any, int, error)

func (h handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseBody, statusCode, err := h(r)

	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.SetEscapeHTML(false)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err != nil {
		log.Printf("Handler error: %v", err)
		jsonEncoder.Encode(dto.ComiketBackendErrorResponse{Message: err.Error()})
	} else {
		jsonEncoder.Encode(responseBody)
	}

}

func RegisterMethodToHTTPServer(mux *http.ServeMux, httpMethod string, path string, handler handlerFunc) {
	matchString := fmt.Sprintf("%s %s", httpMethod, path)
	mux.Handle(matchString, handler)
}
