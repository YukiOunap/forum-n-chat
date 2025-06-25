package server

import (
	"database/sql"
	"log"
)

func InsertMessage(database *sql.DB, message Message, nickname string) PrivateMessage {

	// insert post data into database
	row := database.QueryRow(`
	INSERT INTO private_messages(sender, receiver, content)
	VALUES (?, ?, ?)
	RETURNING sender, receiver, content, created_at`,
		message.Sender, message.Receiver, message.Content)

	var insertedMessage PrivateMessage
	err := row.Scan(&insertedMessage.Sender, &insertedMessage.Receiver, &insertedMessage.Content, &insertedMessage.Time)
	if err != nil {
		log.Println("Unable to retrieve inserted message:", err.Error())
		return PrivateMessage{}
	}

	// update latest time for messaging between the users
	_, err = database.Exec(`
    INSERT INTO last_message_time (user_1, user_2, last_message_time)
    VALUES (
    	CASE WHEN ? < ? THEN ? ELSE ? END,
    	CASE WHEN ? < ? THEN ? ELSE ? END,
    	CURRENT_TIMESTAMP)
    ON CONFLICT(user_1, user_2) DO UPDATE
    SET last_message_time = CURRENT_TIMESTAMP`,
		message.Sender, message.Receiver, // ? < ? の条件1
		message.Sender, message.Receiver, // THEN ? ELSE ? の条件
		message.Sender, message.Receiver, // ? < ? の条件2
		message.Receiver, message.Sender) // THEN ? ELSE ? の条件
	if err != nil {
		log.Println("error inserting last_message_time", err)
	}

	log.Println("TEST", nickname, message)
	_, err = database.Exec(`
			INSERT INTO message_is_read (me, other, is_read)
			VALUES (?, ?, false)
			ON CONFLICT(me, other) DO UPDATE SET is_read = excluded.is_read`,
		message.Receiver, message.Sender)
	if err != nil {
		log.Println("error inserting is_read", err)
	}

	return insertedMessage
}
