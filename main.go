// hosting Go server which offers a web forum service
package main

import (
	"database/sql"
	"git/ykaneko/real-time-forum/server"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	PORT               = ":8080"
	DATABASE_FILE_PATH = "database.db"
)

func main() {

	// open and set the connection to database
	database, err := sql.Open("sqlite3", DATABASE_FILE_PATH)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}
	defer database.Close()

	// enable foreign key setting on SQLite
	if _, err := database.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatalf("unable to activate foreign key mode: %v", err)
	}

	// Serve static files
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "forum.html"))
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.Handle("/profile_pictures/", http.StripPrefix("/profile_pictures/", http.FileServer(http.Dir("web/profile_pictures"))))

	// function handlers

	// login page
	http.HandleFunc("/check-login", func(w http.ResponseWriter, r *http.Request) {
		server.CheckLogin(database, w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		server.LogIn(database, w, r)
	})
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		server.SignUp(database, w, r)
	})

	// content control

	// home page
	http.HandleFunc("/fetch-categories", func(w http.ResponseWriter, r *http.Request) {
		server.FetchCategories(database, w, r)
	})
	http.HandleFunc("/fetch-posts", func(w http.ResponseWriter, r *http.Request) {
		server.FetchPostsByCategory(database, w, r)
	})
	http.HandleFunc("/fetch-post-categories", func(w http.ResponseWriter, r *http.Request) {
		server.FetchPostCategory(database, w, r)
	})
	http.HandleFunc("/fetch-users", func(w http.ResponseWriter, r *http.Request) {
		server.FetchUsers(database, w, r)
	})
	http.HandleFunc("/fetch-is-read", func(w http.ResponseWriter, r *http.Request) {
		log.Println("fetch-is-read")
		server.FetchIsRead(database, w, r)
	})
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		server.LogOut(database, w, r)
	})
	http.HandleFunc("/add-category", func(w http.ResponseWriter, r *http.Request) {
		server.AddCategory(database, w, r)
	})
	http.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		server.CreatePost(database, w, r)
	})

	// post page
	http.HandleFunc("/fetch-post", func(w http.ResponseWriter, r *http.Request) {
		server.FetchPost(database, w, r)
	})
	http.HandleFunc("/fetch-comments", func(w http.ResponseWriter, r *http.Request) {
		server.FetchComments(database, w, r)
	})
	http.HandleFunc("/create-comment", func(w http.ResponseWriter, r *http.Request) {
		server.CreateComment(database, w, r)
	})

	// chat page
	http.HandleFunc("/render-chat-history", func(w http.ResponseWriter, r *http.Request) {
		server.FetchChatHistory(database, w, r)
	})
	http.HandleFunc("/fetch-latest-messages", func(w http.ResponseWriter, r *http.Request) {
		server.FetchLastMessages(database, w, r)
	})
	http.HandleFunc("/update-read-status", func(w http.ResponseWriter, r *http.Request) {
		log.Println("update-read-status")
	})

	// WebSocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.HandleConnections(database, w, r)
	})
	go server.HandleBroadcasts()

	log.Printf("server opened on http://localhost%v\n", PORT)

	// setup the server
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
