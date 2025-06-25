package server

import (
	"database/sql"
	"log"
)

func GetAllUsers(database *sql.DB) []User {

	rows, err := database.Query(`SELECT * FROM users`)
	if err != nil {
		log.Println("Error querying users:", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.Password)
		if err != nil {
			log.Println("Error scanning user:", err)
		}
		users = append(users, user)
	}

	return users
}
