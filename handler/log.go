package handler

import (
	"net/http"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Access Log")
}
