package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marianodsr/routes"
)

func main() {
	r := chi.NewRouter()

	routes.HandleRoutes(r)
	http.ListenAndServe(":8081", r)
}
