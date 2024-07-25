package structs

import "time"

type User struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
	age       int
	gender    int
}

type Post struct {
	ID           int
	Title        string
	Content      string
	CreatedDate  time.Time
	Username     string
	Categories   []string
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
}

type PostContext struct {
	LoggedInUser *User
	Categories   []string
	Post         Post
	Comments     []Comment
}

type Message struct {
	SenderUsername    string
	RecipientUsername string
	Content           string
	CreatedDate       string
}
