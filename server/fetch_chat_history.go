package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func FetchChatHistory(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	// get filter from query
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	sender := r.URL.Query().Get("sender")
	receiver := r.URL.Query().Get("receiver")
	offset := r.URL.Query().Get("offset")

	query := `
		SELECT * FROM private_messages
		WHERE (sender = ? AND receiver = ?)
		   OR (sender = ? AND receiver = ?)
		ORDER BY created_at DESC
		LIMIT 10 OFFSET ?`

	rows, err := database.Query(query, receiver, sender, sender, receiver, offset)
	if err != nil {
		log.Println("Error querying private_messages:", err)
		return
	}
	defer rows.Close()

	var messages []PrivateMessage
	for rows.Next() {
		var message PrivateMessage
		err := rows.Scan(&message.Sender, &message.Receiver, &message.Content, &message.Time)
		if err != nil {
			log.Println("Error scanning private_messages:", err)
			return
		}
		messages = append(messages, message)
	}

	log.Println(messages)

	// set is_read status as true
	_, err = database.Exec(`
        UPDATE message_is_read
		SET is_read = true
		WHERE me = ? AND other = ?`,
		sender, receiver)
	if err != nil {
		log.Println("error setup read status", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("Error encoding messages: %v", err)
	}
}
