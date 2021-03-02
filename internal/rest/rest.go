package rest

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

func renderErrorResponse(w http.ResponseWriter, msg string, status int) {
	renderResponse(w, ErrorResponse{Error: msg}, status)
}

func renderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		// XXX Do something with the error ;)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		// XXX Do something with the error ;)
	}
}
