package controller

import (
	"clubapis/db"
	schemas "clubapis/schema"
	"encoding/json"
	"fmt"
	"net/http"
)

// ----------------- Add Collection -----------------
func AddCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var col schemas.Collection
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Invalid JSON"})
		return
	}

	data, _, err := db.SupaClient.From("collections").
		Insert(map[string]interface{}{
			"admin_id":  col.AdminID,
			"member_id": col.MemberID,
			"amount":    col.Amount,
			"reason":    col.Reason,
			"for_month": col.ForMonth,
			"notes":     col.Notes,
		}, false, "representation", "", "").
		Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{
			Error: "Failed to add collection",
			Data:  err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusCreated, schemas.Response{
		Message: "Collection added successfully",
		Data:    data,
	})
}

// ----------------- Update Collection -----------------
func UpdateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Missing collection id"})
		return
	}

	var col schemas.Collection
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Invalid JSON"})
		return
	}

	data, _, err := db.SupaClient.From("collections").
		Update(map[string]interface{}{
			"member_id": col.MemberID,
			"amount":    col.Amount,
			"reason":    col.Reason,
			"for_month": col.ForMonth,
			"notes":     col.Notes,
		}, "representation", "").
		Eq("id", id).
		Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to update collection"})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{
		Message: "Collection updated successfully",
		Data:    data,
	})
}

// ----------------- Delete Collection -----------------
func DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Missing collection id"})
		return
	}

	_, _, err := db.SupaClient.From("collections").
		Delete("minimal", ""). // <-- two arguments: returning, count
		Eq("id", id).
		Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to delete collection"})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{
		Message: fmt.Sprintf("Collection with ID %s deleted successfully", id),
	})
}

// ----------------- Get All Collections by Admin -----------------
func GetCollectionsByAdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	adminID := r.URL.Query().Get("admin_id")
	if adminID == "" {
		sendJSON(w, http.StatusBadRequest, schemas.Response{Error: "Missing admin_id"})
		return
	}

	data, _, err := db.SupaClient.From("collections").
		Select("*", "0", false).
		Eq("admin_id", adminID).
		Execute()

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to fetch collections"})
		return
	}

	var collections []schemas.Collection
	if err := json.Unmarshal(data, &collections); err != nil {
		sendJSON(w, http.StatusInternalServerError, schemas.Response{Error: "Failed to parse collections"})
		return
	}

	sendJSON(w, http.StatusOK, schemas.Response{
		Message: "Collections fetched successfully",
		Data:    collections,
	})
}
