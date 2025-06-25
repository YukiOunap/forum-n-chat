package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func FetchPostsByCategory(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// get filter from query
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	filter := r.URL.Query().Get("category")
	offset := r.URL.Query().Get("offset")

	var rows *sql.Rows
	var err error

	if filter == "all" {
		query := "SELECT * FROM posts ORDER BY posts.id DESC LIMIT 10 OFFSET ?;"
		rows, err = database.Query(query, offset)
	} else {
		query := `
			SELECT posts.*
			FROM posts
			LEFT JOIN post_categories ON posts.id = post_categories.post_id
			WHERE post_categories.category_name = ?
			ORDER BY posts.id DESC
			LIMIT 10 OFFSET ?;`
		rows, err = database.Query(query, filter, offset)
	}
	if err != nil {
		log.Println("Error querying posts:", err)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Time, &post.Author, &post.NumberOfComments); err != nil {
			log.Println("Error scanning posts:", err)
			return
		}
		posts = append(posts, post)
	}

	// send info to client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
	}); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
