package structs

import "time"

type User struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	Age       int
	Gender    int
}

type Post struct {
	ID           int
	Title        string
	Content      string
	CreatedDate  time.Time
	Username     string
	Category     string
	Interactions []Interaction
	Likes        int
	Dislikes     int
	Comments     int
}

type Comment struct {
	ID          int
	Content     string
	CreatedDate time.Time
	PostID      int
	Username    string
	Likes       int
	Dislikes    int
}

type Interaction struct {
	UserId int
	Kind   int
}

type CommentInteraction struct {
	UserId int
	Kind   int
}

type Session struct {
	Token  string
	Expiry time.Time
	User   User
}

type HomeContext struct {
	LoggedInUser *User
	Posts        []Post
	ErrorMessage string
	Users        []User
}

type PostContext struct {
	LoggedInUser *User
	Categories   []string
	Post         Post
	Comments     []Comment
}

type ChatMessage struct {
	Sender      string
	Recipient   string
	Content     string
	CreatedDate string
}

type PageData struct {
	Ages []int
}

type TypingStatus struct {
    Sender    string 
    Recipient string 
    IsTyping  bool
}

type LastMessage struct {
	Sender    string    // The sender of the last message
	Content   string    // Content of the last message
	Timestamp time.Time// When the last message was sent
}