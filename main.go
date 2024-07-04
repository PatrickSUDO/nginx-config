package main

import (
	"log"
	"net/http"

	"github.com/PatrickSUDO/nginx-config/handlers"
)

func main() {
	r := handlers.RegisterHandlers()
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
