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
	r.Post("/api/club-login", controller.LoginClubHandler)

	r.Post("/api/addMembers", controller.AddMemberHandler)              // Add member
	r.Post("/api/updateMembers", controller.UpdateMemberHandler)        // Update member (?id=…)
	r.Get("/api/getAllMembers", controller.GetMembersByAdminHandler)    // Get all members (?admin_id=…)
	r.Patch("/api/update-status", controller.ToggleMemberActiveHandler) // Enable/Disable (?id=…&active=true/false)
	r.Delete("/api/deleteMember", controller.DeleteMemberHandler)       // Delete member (?id=…)

	// ----------------- Collection Routes -----------------
	r.Post("/api/addCollection", controller.AddCollectionHandler)                // Add collection/payment
	r.Post("/api/updateCollection", controller.UpdateCollectionHandler)          // Update collection (?id=…)
	r.Get("/api/deleteCollection", controller.DeleteCollectionHandler)           // Delete collection (?id=…)
	r.Get("/api/getCollectionsByAdmin", controller.GetCollectionsByAdminHandler) // Get all collections for admin (?admin_id=…)

	http.ListenAndServe(":8080", r)
}
