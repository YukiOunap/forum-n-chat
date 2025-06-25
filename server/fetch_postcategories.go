package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func FetchPostCategory(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("SELECT * FROM post_categories")
	if err != nil {
		log.Println("Error querying post_categories:", err)
		return
	}
	defer rows.Close()

	var postCategories []PostCategory
	for rows.Next() {
		var postCategory PostCategory
		err := rows.Scan(&postCategory.PostID, &postCategory.CategoryName)
		if err != nil {
			log.Println("Error scanning post_category:", err)
			return
		}
		postCategories = append(postCategories, postCategory)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"postCategories": postCategories,
	}); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
