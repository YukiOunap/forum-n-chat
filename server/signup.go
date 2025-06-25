// sign up handler: check if the information are unregistered and create new user in database
package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func SignUp(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	// store form values in User struct
	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	age, _ := strconv.Atoi(r.Form.Get("age"))
	newUser := User{
		Nickname:  r.Form.Get("nickname"),
		Age:       age,
		Gender:    r.Form.Get("gender"),
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Password:  r.Form.Get("password"),
	}

	stmt, err := database.Prepare("INSERT INTO users(nickname, age, gender, first_name, last_name, email, password) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("fail to prepare query", err.Error())
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(newUser.Nickname, newUser.Age, newUser.Gender, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Password)
	if err != nil {
		message := err.Error()
		if err.Error() == "UNIQUE constraint failed: users.email" {
			message = "The email is already registered."
		} else if err.Error() == "UNIQUE constraint failed: users.nickname" {
			message = "The nickname is already registered."
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": message,
		})
		return
	}

	// send new user via WS
	broadcasts <- Message{
		Type: "newUser",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("The user %v is successfully registered!", newUser.Nickname),
	})
}
