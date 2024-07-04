package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PatrickSUDO/nginx-config/config"
	"github.com/PatrickSUDO/nginx-config/nginx"
	"github.com/gorilla/mux"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	yamlPath := r.FormValue("yaml_path")

	cfg, err := config.LoadConfig(yamlPath)
	if err != nil {
		http.Error(w, "Error loading YAML: "+err.Error(), http.StatusBadRequest)
		return
	}

	nginxConfig, err := nginx.GenerateConfig(cfg)
	if err != nil {
		http.Error(w, "Error generating NGINX config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"config": nginxConfig})
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func RegisterHandlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/generate", GenerateHandler).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return r
}
