package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP) // 第1引数にエンドポイント, 第2引数にhandlerを登録
	return mux
}
