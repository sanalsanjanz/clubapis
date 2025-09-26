// main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/supabase-community/supabase-go"
)

// Club represents the club data structure
type Club struct {
	ClubName   string  `json:"club_name" db:"club_name"`
	Contact    string  `json:"contact" db:"contact"`
	Location   string  `json:"location" db:"location"`
	MonthlyFee float64 `json:"monthly_fee" db:"monthly_fee"`
	RegNo      string  `json:"reg_no" db:"reg_no"`
}

// Response structure for API responses
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var supaClient *supabase.Client

func init() {
	// Initialize Supabase client
	supaURL := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_ANON_KEY")

	if supaURL == "" || supaKey == "" {
		log.Fatal("Missing SUPABASE_URL or SUPABASE_ANON_KEY environment variables")
	}

	var err error
	supaClient, err = supabase.NewClient(supaURL, supaKey, nil)
	if err != nil {
		log.Fatal("Failed to initialize Supabase client:", err)
	}
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS configuration (adjust for production)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	r.Post("/api/create-club", createClubHandler)

	// Vercel requires this for Go functions
	http.ListenAndServe(":8080", r)
}

func createClubHandler(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse JSON request body
	var club Club
	if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Error: "Invalid JSON format",
		})
		return
	}

	// Validate required fields
	if club.ClubName == "" || club.Contact == "" || club.Location == "" || club.RegNo == "" {
		sendJSON(w, http.StatusBadRequest, Response{
			Error: "Missing required fields: club_name, contact, location, and reg_no are mandatory",
		})
		return
	}

	// Insert into Supabase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use the correct table name (replace "clubs" with your actual table name)
	table := supaClient.DB.From("clubs")

	// Insert data
	_, err := table.Insert(map[string]interface{}{
		"club_name":   club.ClubName,
		"contact":     club.Contact,
		"location":    club.Location,
		"monthly_fee": club.MonthlyFee,
		"reg_no":      club.RegNo,
	}).Execute(ctx)

	if err != nil {
		log.Printf("Supabase insert error: %v", err)
		sendJSON(w, http.StatusInternalServerError, Response{
			Error: "Failed to create club account",
		})
		return
	}

	sendJSON(w, http.StatusCreated, Response{
		Message: "Club account created successfully",
		Data:    club,
	})
}

// Helper function to send JSON responses
func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
