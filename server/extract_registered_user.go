package server

import (
	"database/sql"
)

func ExtractRegisteredUser(database *sql.DB, user User) (User, error) {

	registeredUser := User{}

	row := database.QueryRow(`SELECT * FROM users WHERE nickname = ?`, user.Nickname)
	err := row.Scan(&registeredUser.Nickname, &registeredUser.Age, &registeredUser.Gender, &registeredUser.FirstName, &registeredUser.LastName, &registeredUser.Email, &registeredUser.Password)
	if err == nil {
		return registeredUser, nil
	}

	row = database.QueryRow(`SELECT * FROM users WHERE email = ?`, user.Nickname)
	err = row.Scan(&registeredUser.Nickname, &registeredUser.Age, &registeredUser.Gender, &registeredUser.FirstName, &registeredUser.LastName, &registeredUser.Email, &registeredUser.Password)
	if err == nil {
		return registeredUser, nil
	}

	return User{}, err
}
