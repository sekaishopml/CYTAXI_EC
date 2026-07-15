package http

import "net/http"

type Server interface {
	Start() error
	Shutdown() error
	Route(method, path string, handler http.HandlerFunc)
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    any         `json:"data,omitempty"`
	Errors  []ErrorItem `json:"errors,omitempty"`
}

type ErrorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type Handler func(w http.ResponseWriter, r *http.Request) error
