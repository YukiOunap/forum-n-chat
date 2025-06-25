// log in handler: display and manage log in functions
package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func LogIn(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	// store form values in User struct
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	user := User{
		Nickname: r.FormValue("nickname"), // possibly email
		Password: r.FormValue("password"),
	}

	// check if there is the user
	registeredUser, err := ExtractRegisteredUser(database, user)
	if err != nil {
		response := map[string]string{
			"loginStatus": "fail",
			"message":     "The nickname/e-mail is not registered.",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// check if the password matches
	if user.Password != registeredUser.Password {
		response := map[string]string{
			"loginStatus": "fail",
			"message":     "Incorrect password.",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// login success

	// create new session & cookie
	sessionID := CreateSession(database, w, registeredUser.Nickname)
	cookie := http.Cookie{
		Name:     "session",
		Value:    sessionID,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"loginStatus": "success",
		"user":        registeredUser.Nickname,
	})
}
