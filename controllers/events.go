package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-backend-technical-test/database"
	"go-backend-technical-test/models"
	"github.com/gorilla/mux"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 
	}

	rows, err := database.DB.Query("SELECT id, name, description, date FROM events LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Date)
		if err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		events = append(events, event)
	}

	var totalRows int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM events").Scan(&totalRows)
	if err != nil {
		http.Error(w, "Error fetching total rows", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"events":    events,
		"totalRows": totalRows,
		"limit":     limit,
		"offset":    offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var event models.Event
	query := "SELECT id, name, description, date FROM events WHERE id = ?"
	err = database.DB.QueryRow(query, id).Scan(&event.ID, &event.Name, &event.Description, &event.Date)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO events (name, description, date) VALUES (?, ?, ?)"
	_, err := database.DB.Exec(query, event.Name, event.Description, event.Date)
	if err != nil {
		http.Error(w, "Error saving event to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Event Created")
}

func EditEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	queryExisting := "SELECT id, name, description, date FROM events WHERE id = ?"
	err = database.DB.QueryRow(queryExisting, id).Scan(&event.ID, &event.Name, &event.Description, &event.Date)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := "UPDATE events SET name = ?, description = ?, date = ? WHERE id = ?"
	_, err = database.DB.Exec(query, event.Name, event.Description, event.Date, id)
	if err != nil {
		http.Error(w, "Error updating event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Event Updated")
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	query := "DELETE FROM events WHERE id = ?"
	_, err = database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Event Deleted")
}
