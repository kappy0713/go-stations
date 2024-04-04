package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", middleware.Log(middleware.OS(handler.NewHealthzHandler())).ServeHTTP) // 第1引数にエンドポイント, 第2引数にhandlerを登録
	mux.HandleFunc("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)).ServeHTTP)
	mux.HandleFunc("/do-panic", middleware.Recovery(handler.NewPanicHandler()).ServeHTTP)
	// ルータを分ける必要はない
	// mux.HandleFunc("/os", middleware.OS(handler.NewOSHandler()).ServeHTTP)
	mux.HandleFunc("/log", middleware.Log(middleware.OS(handler.NewLogHandler())).ServeHTTP)
	mux.HandleFunc("/auth", middleware.BasicAuth(handler.NewBasicAuthHandler()).ServeHTTP)
	return mux
}
