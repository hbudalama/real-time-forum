package main

import (
	"log"
	"net/http"
	"os"
	"rtf/pkg/db"
	"rtf/pkg/server"
)

func main() {
	dbErr := db.Connect()
	if dbErr != nil {
		log.Fatal("Error connecting to the database: ", dbErr)
	}
	defer db.Close()

	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// Serve the single HTML file for all routes except API routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// API handlers
	mux.HandleFunc("/api/login", server.LoginHandler)
	mux.HandleFunc("/api/posts", server.PostHandler)
	mux.HandleFunc("/api/posts/{id}/comments", server.CommentsHandler)
	mux.HandleFunc("/api/posts/{id}", server.GetPostHandler)
	mux.HandleFunc("/api/posts/{id}/dislike", server.AddDislikesHandler)
	mux.HandleFunc("/api/posts/{id}/like", server.AddLikesHandler)
	mux.HandleFunc("/api/signup", server.SignupHandler)
	mux.HandleFunc("/api/logout", server.LogoutHandler)
	mux.HandleFunc("/api/comments/{id}/like", server.LikeCommentHandler)
	mux.HandleFunc("/api/comments/{id}/dislike", server.DislikeCommentHandler)
	mux.HandleFunc("/api/myPosts", server.MyPostsHandler)
	mux.HandleFunc("/api/newest", server.NewestPostsHandler)

	log.Println("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("Error starting server:", err)
		os.Exit(1)
	}
}
