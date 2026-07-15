package http

import (
	"encoding/json"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/errors"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status: status,
		Data:   data,
	})
}

func WriteError(w http.ResponseWriter, err error) {
	status, body := errors.EncodeHTTP(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func WriteValidationError(w http.ResponseWriter, fieldErrors any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Status: http.StatusBadRequest,
		Errors: fieldErrors,
	})
}

func OK(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
