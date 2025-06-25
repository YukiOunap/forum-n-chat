// check session function: extract sessions from browser and database. Return cookie and session struct but return the struct with UserID 0 when there's no valid session
package server

import (
	"database/sql"
	"log"
	"net/http"
)

func CheckSession(database *sql.DB, w http.ResponseWriter, r *http.Request) (*http.Cookie, Session) {

	// Initialized session structure (invalid state)
	session := Session{Nickname: ""}

	// extract session from cookie
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		log.Println("No valid session in browser:", err)
		return nil, session
	}

	// extract session data from database
	err = database.QueryRow("SELECT * FROM Sessions WHERE session_id = ?", cookie.Value).Scan(&session.SessionID, &session.Nickname, &session.Time)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("Failed to query session from database:", err)
		}
		return nil, session
	}

	return cookie, session
}
