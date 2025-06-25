package server

import (
	"database/sql"
	"log"
	"net/http"
)

type CommentFormData struct {
	PostID     string
	Content    string
	Author     string
	Categories []string
}

func CreateComment(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// fetch form data to instance
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	newComment := CommentFormData{
		PostID:  r.FormValue("postId"),
		Content: r.FormValue("content"),
		Author:  r.FormValue("author"),
	}

	// starts transaction
	tx, err := database.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back: %v\n", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("Failed to commit transaction: %v\n", err)
			}
		}
	}()

	// insert post data into database
	row := tx.QueryRow(`
		INSERT INTO Comments(post_id, content, author)
		VALUES(?, ?, ?)
		RETURNING id, post_id, author, content, created_at`,
		newComment.PostID, newComment.Content, newComment.Author)

	var createdComment Comment
	err = row.Scan(&createdComment.ID, &createdComment.PostID, &createdComment.Author, &createdComment.Content, &createdComment.Time)
	if err != nil {
		log.Println("Unable to retrieve inserted comment:", err.Error())
		return
	}

	// insert comment for database and increment comment number for the post
	_, err = tx.Exec(`
		UPDATE posts
		SET number_of_comments = number_of_comments + 1
		WHERE id = ?`, newComment.PostID)
	if err != nil {
		log.Printf("Failed to update number_of_comments: %v\n", err)
		return
	}

	// make update for WS broadcast
	broadcasts <- Message{
		Type:    "commentUpdate",
		Content: createdComment}
}
