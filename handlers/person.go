package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"people-api/database"
	"people-api/models"

	"github.com/gorilla/mux"
)

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO people (name, email, phone, stack) VALUES (?, ?, ?, ?)",
		req.Name, req.Email, req.Phone, req.Stack,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	response := map[string]interface{}{
		"id":      id,
		"message": "Person created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, email, phone, stack, created_at FROM people")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Phone, &p.Stack, &p.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		people = append(people, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p models.Person
	err = database.DB.QueryRow(
		"SELECT id, name, email, phone, stack, created_at FROM people WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Name, &p.Email, &p.Phone, &p.Stack, &p.CreatedAt)

	if err != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(
		"UPDATE people SET name = ?, email = ?, phone = ?, stack = ? WHERE id = ?",
		req.Name, req.Email, req.Phone, req.Stack, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Person updated successfully"})
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("DELETE FROM people WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Person deleted successfully"})
} 