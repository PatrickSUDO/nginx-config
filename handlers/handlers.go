package handlers

import (
	"net/http"

	"github.com/PatrickSUDO/nginx-config/config"
	"github.com/PatrickSUDO/nginx-config/nginx"
	"github.com/gorilla/mux"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig("input.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nginxConfig, err := nginx.GenerateConfig(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=nginx.conf")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte(nginxConfig))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func RegisterHandlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/generate", GenerateHandler).Methods("POST")
	return r
}
