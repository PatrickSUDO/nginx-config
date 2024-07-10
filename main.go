package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PatrickSUDO/nginx-config/config"
	"github.com/PatrickSUDO/nginx-config/nginx"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	r.HandleFunc("/", serveIndex)
	r.HandleFunc("/generate", generateConfigHandler).Methods("POST")

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func generateConfigHandler(w http.ResponseWriter, r *http.Request) {
	var cfg *config.Config

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("yaml_file")
	if err == nil {
		// File was uploaded
		defer file.Close()
		yamlData, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		cfg, err = config.LoadConfigFromBytes(yamlData)
		if err != nil {
			http.Error(w, "Error parsing YAML", http.StatusBadRequest)
			return
		}
	} else {
		// No file uploaded, use JSON data
		configJSON := r.FormValue("config")
		if configJSON == "" {
			http.Error(w, "No configuration provided", http.StatusBadRequest)
			return
		}
		cfg = &config.Config{}
		if err := json.Unmarshal([]byte(configJSON), cfg); err != nil {
			http.Error(w, "Error parsing JSON configuration", http.StatusBadRequest)
			return
		}
	}

	nginxConfig, err := nginx.GenerateConfig(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"config": nginxConfig})
}
