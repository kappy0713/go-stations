package handler

import (
	"net/http"
)

type BasicAuthHandler struct{}

func NewBasicAuthHandler() *BasicAuthHandler {
	return &BasicAuthHandler{}
}

func (h *BasicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
