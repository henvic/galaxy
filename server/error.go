package server

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ErrorResponse from the API.
type ErrorResponse struct {
	Status  int    `json:"status" example:"404"`
	Message string `json:"message" example:"Not Found"`
	Errors  string `json:"errors,omitempty" example:"reference not found"`
}

// ErrorHandler for errors.
// Use code = -1 to bypass writing header or trying to set Content-Type. Useful when headers are already sent.
func ErrorHandler(w http.ResponseWriter, r *http.Request, code int, errors ...string) {
	if code >= 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.WriteHeader(code)
	}

	if code != http.StatusNotFound {
		log.Debugf("Request error: %v (%v)", code, http.StatusText(code))
	}

	e := json.NewEncoder(w)

	e.SetIndent("", "    ")

	if err := e.Encode(ErrorResponse{
		Status:  code,
		Message: http.StatusText(code),
		Errors:  strings.Join(errors, "\n"),
	}); err != nil {
		log.Error(err)
	}
}
