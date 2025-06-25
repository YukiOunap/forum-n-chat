package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func CheckLogin(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	_, session := CheckSession(database, w, r)

	// not logged in
	if session.Nickname == "" {
		log.Println("unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println(map[string]string{
		"user": session.Nickname,
	})

	// logged in
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"user": session.Nickname,
	})
}
