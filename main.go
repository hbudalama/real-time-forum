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

	// Handle the single-page application for logged-in users
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "pages/index.html")
	})

	// API handlers
	mux.HandleFunc("/api/ws", server.Echo)

	mux.HandleFunc("/api/login", server.LoginHandler) //fetched
	mux.HandleFunc("/api/posts", server.PostHandler) // fetched
	mux.HandleFunc("/api/posts/{id}/comments", server.CommentsHandler) //fetched
	mux.HandleFunc("/api/posts/{id}", server.GetPostHandler) // unused path but everything works til now
	mux.HandleFunc("/api/add-post", server.AddPostsHandler) //fetched
	mux.HandleFunc("/api/posts/{id}/dislike", server.AddDislikesHandler) //fetched
	mux.HandleFunc("/api/posts/{id}/like", server.AddLikesHandler) //fetched
	mux.HandleFunc("/api/signup", server.SignupHandler) //fetched
	mux.HandleFunc("/api/logout", server.LogoutHandler) //fetched
	mux.HandleFunc("/api/comments/{id}/like", server.LikeCommentHandler) //fetched
	mux.HandleFunc("/api/comments/{id}/dislike", server.DislikeCommentHandler) //fetched
	mux.HandleFunc("/api/check_session", server.CheckSessionHandler) // fetched
	mux.HandleFunc("/api/get_ages", server.GetAgesHandler) //fetched
	mux.HandleFunc("/api/get_user_info", server.GetUserInfoHandler) //fetched
	mux.HandleFunc("/api/usernames", server.GetAllUsernamesHandler)
	

	log.Println("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("Error starting server:", err)
		os.Exit(1)
	}
}
