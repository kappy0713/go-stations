package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("healthzhandler")
	var x = &model.HealthzResponse{}
	x.Message = "OK"

	err := json.NewEncoder(w).Encode(x)
	if err != nil {
		log.Println(err)
	}
}
