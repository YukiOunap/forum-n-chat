// generate session and insert to database
package server

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

func CreateSession(database *sql.DB, w http.ResponseWriter, nickname string) string {

	// generate random session ID
	sessionID := ""
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println("unable to read the random bytes", err.Error())
	}
	sessionID = base64.URLEncoding.EncodeToString(randomBytes)

	// Set session expiration time (1 hour)
	expiresAt := time.Now().Add(1 * time.Hour).UTC()

	//insert session into DB
	_, err = database.Exec("INSERT INTO sessions (session_id, nickname, expires_at) VALUES (?, ?, ?)", sessionID, nickname, expiresAt)
	if err != nil {
		log.Println("unable to insert the session into database", err.Error())
	}

	// set sessionID to cookie
	cookie := &http.Cookie{
		Name:    "sessionID",
		Value:   sessionID,
		Expires: expiresAt,
	}
	http.SetCookie(w, cookie)

	return sessionID
}
