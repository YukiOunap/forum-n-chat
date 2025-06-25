package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func AddCategory(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// fetch form data to instance
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newCategory := r.FormValue("newCategory")

	// insert category data into database
	response := map[string]string{"status": "error"}
	_, err = database.Exec("INSERT INTO categories(name) VALUES(?)", newCategory)
	if err != nil {
		log.Println("error inserting new category to database:", err.Error())
		if err.Error() == "UNIQUE constraint failed: categories.name" {
			response["status"] = "duplicates"
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// make update for WS broadcast
	broadcasts <- Message{Type: "categoryUpdate", Content: newCategory}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response["status"] = "success"
	json.NewEncoder(w).Encode(response)
}
