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
			"password":    club.Password,
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

func GetClubsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query all rows from "clubs"
	data, _, err := db.SupaClient.
		From("clubs").
		Select("*", "", false).
		Execute()

	if err != nil {
		log.Printf("Supabase error: %v", err)
		sendJSON(w, http.StatusInternalServerError, schemas.Response{
			Error: "Failed to fetch clubs",
		})
		return
	}

	// Unmarshal JSON into slice of Club
	var clubs []schemas.Club
	if err := json.Unmarshal(data, &clubs); err != nil {
		log.Printf("Unmarshal error: %v", err)
		sendJSON(w, http.StatusInternalServerError, schemas.Response{
			Error: "Failed to parse clubs data",
		})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{
		Message: "Clubs fetched successfully",
		Data:    clubs,
	})
}

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("JSON encode error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
func LoginClubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse incoming JSON
	var creds struct {
		RegNo    string `json:"reg_no"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendJSON(w, http.StatusBadRequest, schemas.Response{
			Error: "Invalid JSON",
		})
		return
	}

	// Validate required fields
	if creds.RegNo == "" || creds.Password == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{
			Error: "Missing required fields: reg_no and password",
		})
		return
	}

	// Query Supabase for club with reg_no and password
	data, _, err := db.SupaClient.
		From("clubs").
		Select("*", "", false).
		Eq("reg_no", creds.RegNo).
		Eq("password", creds.Password). // simple password check
		Single().                       // ensures single record
		Execute()

	if err != nil {
		log.Printf("Supabase login error: %v", err)
		sendJSON(w, http.StatusUnauthorized, schemas.Response{
			Error: "Invalid credentials",
		})
		return
	}

	// Unmarshal returned club JSON
	var club schemas.Club
	if err := json.Unmarshal(data, &club); err != nil {
		log.Printf("Unmarshal error: %v", err)
		sendJSON(w, http.StatusInternalServerError, schemas.Response{
			Error: "Failed to parse club data",
		})
		return
	}

	// Success response
	sendJSON(w, http.StatusOK, schemas.Response{
		Message: "Login successful",
		Data: map[string]interface{}{
			"reg_no": creds.RegNo,
			"club":   club,
		},
	})
}
