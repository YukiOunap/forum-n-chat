package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func FetchIsRead(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// get filter from query
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	me := r.URL.Query().Get("me")
	other := r.URL.Query().Get("other")
	log.Println(me, other)

	isRead := messageIsRead{}
	row := database.QueryRow("SELECT * FROM message_is_read WHERE me = ? AND other = ?", me, other)
	err := row.Scan(&isRead.Me, &isRead.Other, &isRead.IsRead)
	if err != nil {
		log.Println("Error scanning message_is_read:", err)
		return
	}

	log.Println(isRead)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"isRead": isRead.IsRead,
	}); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
