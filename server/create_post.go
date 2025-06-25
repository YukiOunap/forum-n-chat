package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type PostFormData struct {
	Title      string
	Content    string
	Author     string
	Categories []string
}

func CreatePost(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// fetch form data to instance
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newPost := PostFormData{
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		Author:     r.FormValue("author"),
		Categories: r.Form["category"],
	}

	// insert post data into database
	row := database.QueryRow(`
		INSERT INTO Posts(title, content, author) 
		VALUES(?, ?, ?)
		RETURNING id, title, content, created_at, author, number_of_comments`,
		newPost.Title, newPost.Content, newPost.Author)

	var insertedPost Post
	err = row.Scan(&insertedPost.ID, &insertedPost.Title, &insertedPost.Content, &insertedPost.Time, &insertedPost.Author, &insertedPost.NumberOfComments)
	if err != nil {
		log.Println("Unable to retrieve inserted post:", err.Error())
	}

	// insert category data for the post into database
	for _, category := range newPost.Categories {
		_, err := database.Exec("INSERT INTO post_categories(post_id, category_name) VALUES(?, ?)", insertedPost.ID, category)
		if err != nil {
			log.Println("unable to insert category to new post in database:", err.Error())
			return
		}
	}

	// make update for WS broadcast
	insertedPost.Categories = newPost.Categories
	broadcasts <- Message{Type: "postUpdate", Content: insertedPost}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
