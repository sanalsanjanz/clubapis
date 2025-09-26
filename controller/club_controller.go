package controller

import (
	"clubapis/db"
	schemas "clubapis/schema"
	"encoding/json"
	"log"
	"net/http"
)

func CreateClubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var club schemas.Club
	if err := json.NewDecoder(r.Body).Decode(&club); err != nil {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Invalid JSON"})
		return
	}

	// Validate required fields
	if club.ClubName == "" || club.Contact == "" || club.Location == "" || club.RegNo == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{
			Error: "Missing required fields: club_name, contact, location, reg_no",
		})
		return
	}

	// ✅ No context needed in Execute() in v0.10+
	_, _, err := db.SupaClient.From("clubs").Insert(
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
		sendJSON(w, http.StatusInternalServerError, schemas.Response{
			Error: "Failed to create club account",
		})
		return
	}

	sendJSON(w, http.StatusCreated, schemas.Response{
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
