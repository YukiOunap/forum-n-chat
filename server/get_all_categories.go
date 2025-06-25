package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func FetchCategories(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT * FROM categories")
	if err != nil {
		log.Println("unable to extract categories", err.Error())
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.Name); err != nil {
			log.Println("unable to scan categories", err.Error())
		}
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categories,
	}); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
