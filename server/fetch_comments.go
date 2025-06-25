package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func FetchComments(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// get filter from query
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	postID := r.URL.Query().Get("post")
	offset := r.URL.Query().Get("offset")

	rows, err := database.Query(`
		SELECT * FROM comments
		WHERE post_id = ?
		ORDER BY id DESC
		LIMIT 10 OFFSET ?`, postID, offset)
	if err != nil {
		log.Println("Error querying comments:", err)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.Author, &comment.Content, &comment.Time)
		if err != nil {
			log.Println("Error scanning comment:", err)
			return
		}
		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		log.Printf("Error encoding messages: %v", err)
	}
}
