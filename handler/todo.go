package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// ServeHTTP implements http.Handler.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req := &model.CreateTODORequest{}

		deq := json.NewDecoder(r.Body)
		if err := deq.Decode(req); err != nil { // Decodeはアドレスを渡す必要がある
			log.Println(err)
		}

		if req.Subject == "" { //Subjectが空でないか判定, 空なら400を返す
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		res, err := h.Create(r.Context(), req)
		if err != nil {
			log.Println(err)
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}
		return

	case "PUT":
		req := &model.UpdateTODORequest{}

		deq := json.NewDecoder(r.Body)
		if err := deq.Decode(req); err != nil {
			log.Println(err)
		}

		if req.Subject == "" || req.ID == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		res, err := h.Update(r.Context(), req)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}
		return
	case "GET":
		req := &model.ReadTODORequest{}
		var err error

		prevID := r.URL.Query().Get("prev_id")             // string
		req.PrevID, err = strconv.ParseInt(prevID, 10, 64) // int64に変換, reqに値を代入
		if err != nil {
			log.Println(err)
		}

		size := r.URL.Query().Get("size")
		req.Size, err = strconv.ParseInt(size, 10, 64)
		if req.Size == 0 {
			req.Size = 5 // sizeのdefault値は5
		}
		if err != nil {
			log.Println(err)
		}

		res, err := h.Read(r.Context(), req)
		if err != nil {
			log.Println(err)
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}
		return
	case "DELETE":
		req := &model.DeleteTODORequest{}

		deq := json.NewDecoder(r.Body)
		if err := deq.Decode(req); err != nil {
			log.Println(err)
		}
		if len(req.IDs) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		res, err := h.Delete(r.Context(), req)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}
		return
	}
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *todo}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	return &model.ReadTODOResponse{TODOs: todos}, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	ids := make([]int64, len(req.IDs))
	for i, id := range req.IDs {
		ids[i] = int64(id)
	}
	err := h.svc.DeleteTODO(ctx, ids)
	return &model.DeleteTODOResponse{}, err
}
