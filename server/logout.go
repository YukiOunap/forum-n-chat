// log out handler: delete valid session
package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func LogOut(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// check session and extract userID
	cookie, session := CheckSession(database, w, r)

	// delete session from database
	_, err := database.Exec("DELETE FROM sessions WHERE session_id = ?", session.SessionID)
	if err != nil {
		log.Println("unable to delete the session data in database", err.Error())
	}

	// unable cookie
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
