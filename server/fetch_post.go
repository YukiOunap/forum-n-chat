package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func FetchPost(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// get filter from query
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")

	post := Post{}

	row := database.QueryRow("SELECT * FROM posts WHERE id = ?", id)
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Time, &post.Author, &post.NumberOfComments); err != nil {
		log.Println("Error scanning post:", err)
		return
	}

	log.Println(post)

	// send info to client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"post": post,
	}); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
