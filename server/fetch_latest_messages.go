package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func FetchLastMessages(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// get filter from query
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	user := r.URL.Query().Get("user")

	rows, err := database.Query("SELECT * FROM last_message_time WHERE user_1 = ? OR user_2 = ?", user, user)
	if err != nil {
		log.Println("Error querying user_last_message_time:", err)
		return
	}
	defer rows.Close()

	var latestMassages []lastMessageTime
	for rows.Next() {
		var latestMassage lastMessageTime
		err := rows.Scan(&latestMassage.User1, &latestMassage.User2, &latestMassage.LastMessageTime)
		if err != nil {
			log.Println("Error scanning user_last_message_time:", err)
			return
		}
		latestMassages = append(latestMassages, latestMassage)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"latestMessages": latestMassages,
	}); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
