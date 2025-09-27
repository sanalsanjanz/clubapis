package controller

import (
	"clubapis/db"
	schemas "clubapis/schema"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddMemberHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var member schemas.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Invalid JSON"})
		return
	}

	// Required fields
	if member.AdminID == "" || member.Name == "" || member.Mobile == "" || member.Email == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Missing required fields"})
		return
	}

	// Insert into Supabase
	_, _, err := db.SupaClient.From("members").Insert(
		map[string]interface{}{
			"admin_id": member.AdminID,
			"name":     member.Name,
			"mobile":   member.Mobile,
			"email":    member.Email,
			"role":     member.Role,
			"active":   true,
			// "club_name": member.ClubName,
			// "contact":   member.Contact,
			// "location":  member.Location,
			// "reg_no":    member.RegNo,
		},
		false,            // upsert = false
		"representation", // returning = full row back
		"",               // count = none
		"",               // schema = default/public
	).Execute()

	if err != nil {
		// Include the detailed Supabase error in the response
		sendJSON(w, http.StatusInternalServerError, schemas.Response{
			Error: "Failed to add member",
			Data:  err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusCreated, schemas.Response{
		Message: "Member added successfully",
		Data:    member,
	})
}

func UpdateMemberHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	memberID := r.URL.Query().Get("id")
	if memberID == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{
			Error: fmt.Sprintf("Missing member id: %s", memberID),
		})
		return
	}
	var member schemas.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Invalid JSON"})
		return
	}

	_, _, err := db.SupaClient.
		From("members").
		Update(map[string]interface{}{
			"name":   member.Name,
			"mobile": member.Mobile,
			"email":  member.Email,
			"role":   member.Role,
		}, "", ""). // <-- returning="" and count=""
		Eq("id", memberID).
		Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to update member"})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{Message: "Member updated successfully"})
}
func GetMembersByAdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	adminID := r.URL.Query().Get("admin_id")
	if adminID == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Missing admin_id"})
		return
	}

	data, _, err := db.SupaClient.From("members").Select("*", "", false).
		Eq("admin_id", adminID).Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to fetch members"})
		return
	}

	var members []schemas.Member
	if err := json.Unmarshal(data, &members); err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to parse members"})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{
		Message: "Members fetched successfully",
		Data:    members,
	})
}
func ToggleMemberActiveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	memberID := r.URL.Query().Get("id")
	active := r.URL.Query().Get("active") // pass "true" or "false"

	if memberID == "" || active == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Missing member id or active param"})
		return
	}

	isActive := active == "true"

	// ✅ Include returning + count arguments
	_, _, err := db.SupaClient.From("members").Update(
		map[string]interface{}{"active": isActive},
		"minimal", // or "representation" if you want the updated row back
		"",
	).Eq("id", memberID).Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to update active status"})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{
		Message: "Member status updated successfully",
	})
}

func DeleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	memberID := r.URL.Query().Get("id")
	if memberID == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Missing member id"})
		return
	}

	// Use Delete("minimal", "") for no return rows
	_, _, err := db.SupaClient.
		From("members").
		Delete("minimal", ""). // 👈 fixed here
		Eq("id", memberID).
		Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to delete member"})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{
		Message: "Member deleted successfully",
	})
}
