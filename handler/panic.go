package handler

import (
	"net/http"
)

// A HealthzHandler implements health check endpoint.
type PanicHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic発生") // panicが起きた時のメッセージ
}
