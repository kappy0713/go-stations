package handler

import (
	"net/http"
)

type OSHandler struct{}

func NewOSHandler() *OSHandler {
	return &OSHandler{}
}

func (h *OSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("OS")
}
