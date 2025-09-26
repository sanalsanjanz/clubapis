package main

import (
	"clubapis/controller"
	"clubapis/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func init() {
	// Load .env file (only during local development)
	db.InitDB()

}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         300,
	}))

	r.Post("/api/create-club", controller.CreateClubHandler)
	r.Get("/api/get-all-clubs", controller.GetClubsHandler)
	http.ListenAndServe(":8080", r)
}
