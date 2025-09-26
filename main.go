package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func init() {
	// Load .env file (only during local development)
	if os.Getenv("VERCEL") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found (this is normal in production)")
		}
	}

	supaURL := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_ANON_KEY")

	// DEBUG: Uncomment these lines to see what's being loaded
	// log.Printf("SUPABASE_URL: '%s'", supaURL)
	// log.Printf("SUPABASE_ANON_KEY length: %d", len(supaKey))

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
	r.Use(middleware.Logger, middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         300,
	}))

	r.Post("/api/create-club", apis.createClubHandler)
	http.ListenAndServe(":8080", r)
} /*

func createClubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var club Club
	if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Error: "Invalid JSON"})
		return
	}

	// Validate required fields
	if club.ClubName == "" || club.Contact == "" || club.Location == "" || club.RegNo == "" {
		sendJSON(w, http.StatusBadRequest, Response{
			Error: "Missing required fields: club_name, contact, location, reg_no",
		})
		return
	}

	// ✅ No context needed in Execute() in v0.10+
	_, _, err := supaClient.From("clubs").Insert(
		map[string]interface{}{
			"club_name":   club.ClubName,
			"contact":     club.Contact,
			"location":    club.Location,
			"monthly_fee": club.MonthlyFee,
			"reg_no":      club.RegNo,
		},
		false,            // upsert
		"",               // onConflict
		"representation", // returning
		"",               // count
	).Execute()

	if err != nil {
		log.Printf("Supabase error: %v", err)
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

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("JSON encode error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
*/
