package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{}

var (
	onlineUsers = make(map[string]map[string]*websocket.Conn) // a map has usernames as keys and the value is the map of the users' connections by sessionID as keys
	mu          sync.Mutex
	broadcasts  = make(chan Message)
)

type Message struct {
	Type     string      `json:"type"`
	Sender   string      `json:"sender"`
	Receiver string      `json:"receiver"`
	Content  interface{} `json:"content"`
}

func HandleConnections(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	_, session := CheckSession(database, w, r)
	log.Println("WS connection established for ", session)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}

	// close ws when all the other process stopped
	defer func() {
		log.Println("Closing WS:", session)
		CloseConnection(session)
	}()

	// add online user
	mu.Lock()
	if onlineUsers[session.Nickname] == nil {
		onlineUsers[session.Nickname] = make(map[string]*websocket.Conn)
	}
	onlineUsers[session.Nickname][session.SessionID] = ws
	mu.Unlock()

	// broadcast online status update
	msg := Message{
		Type:    "statusUpdate",
		Sender:  session.Nickname,
		Content: "online",
	}
	broadcasts <- msg

	var wg sync.WaitGroup
	wg.Add(1)

	// wait for private messages
	go func() {
		defer wg.Done()
		for {

			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			HandleMessages(database, msg, session.Nickname)
		}
	}()

	wg.Wait()
}

func CloseConnection(session Session) {
	mu.Lock()
	if user, exists := onlineUsers[session.Nickname]; exists {
		if _, exists := user[session.SessionID]; exists {
			delete(user, session.SessionID)
			if len(user) == 0 {
				delete(onlineUsers, session.Nickname)
				msg := Message{
					Type:    "statusUpdate",
					Sender:  session.Nickname,
					Content: "offline",
				}
				broadcasts <- msg
			}
		}
	}
	mu.Unlock()
}

func HandleBroadcasts() {
	for broadcast := range broadcasts {
		messageJSON, err := json.Marshal(broadcast)
		if err != nil {
			log.Println("Error marshaling status update:", err)
			continue
		}

		mu.Lock()
		for _, user := range onlineUsers {
			broadcastToAllConnections(messageJSON, user)
		}
		mu.Unlock()
	}
}

func HandleMessages(database *sql.DB, msg Message, nickname string) {

	insertedMessage := InsertMessage(database, msg, nickname)

	msg.Content = insertedMessage
	messageJSON, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshaling status update:", err)
		return
	}

	// update the msg to the sender
	mu.Lock()
	sender, online := onlineUsers[msg.Sender]
	if online {
		broadcastToAllConnections(messageJSON, sender)
	}
	mu.Unlock()

	// update the msg to the receiver
	mu.Lock()
	receiver, online := onlineUsers[msg.Receiver]
	if online {
		broadcastToAllConnections(messageJSON, receiver)
	}
	mu.Unlock()

}

func broadcastToAllConnections(messageJSON []byte, user map[string]*websocket.Conn) {
	for _, ws := range user {
		err := ws.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			log.Println("Error sending status update:", err)
			continue
		}
	}
}
