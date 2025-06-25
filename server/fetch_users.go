package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func FetchUsers(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// filter unnecessary info
	otherUsers := []string{}
	for _, userData := range GetAllUsers(database) {
		otherUsers = append(otherUsers, userData.Nickname)
	}
	otherOnlineUsers := []string{}
	for userName := range onlineUsers {
		otherOnlineUsers = append(otherOnlineUsers, userName)
	}

	// send info to client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"users":       otherUsers,
		"onlineUsers": otherOnlineUsers,
	}); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
