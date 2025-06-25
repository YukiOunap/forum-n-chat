package server

import "time"

type Session struct {
	SessionID string    `json:"sessionId"`
	Nickname  string    `json:"nickname"`
	Time      time.Time `json:"time"`
}

type User struct {
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Post struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Time             time.Time `json:"time"`
	Author           string    `json:"author"`
	Categories       []string  `json:"categories"`
	NumberOfComments int       `json:"numberOfComments"`
}

type Category struct {
	Name string `json:"name"`
}

type PostCategory struct {
	PostID       int    `json:"postId"`
	CategoryName string `json:"categoryName"`
}

type Comment struct {
	ID      int       `json:"id"`
	PostID  int       `json:"postId"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
	Author  string    `json:"author"`
}

type PrivateMessage struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

type lastMessageTime struct {
	User1           string    `json:"user1"`
	User2           string    `json:"user2"`
	LastMessageTime time.Time `json:"lastMessageTime"`
}

type messageIsRead struct {
	Me     string `json:"sender"`
	Other  string `json:"receiver"`
	IsRead bool   `json:"isRead"`
}
