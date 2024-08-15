// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"rtf/pkg/db"
// 	"rtf/pkg/server"
// )

// func main() {
// 	dbErr := db.Connect()
// 	if dbErr != nil {
// 		log.Fatal("Error connecting to the database: ", dbErr)
// 	}
// 	defer db.Close()
// 	mux := http.NewServeMux()
// 	// Serve static files
// 	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
// 	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

// 	authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			if r.URL.Path == "/login" || r.URL.Path == "/signup" || r.URL.Path == "/api/login" || r.URL.Path == "/api/signup" {
// 				next.ServeHTTP(w, r)
// 				return
// 			}
// 			// Check if the user is authenticated
// 			cookie, err := r.Cookie("session_token")
// 			if err != nil || !db.IsSessionValid(cookie.Value) {
// 				http.Redirect(w, r, "/login", http.StatusSeeOther)
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 		}
// 	}

// 	// Handle login and registration separately
// 	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, "pages/index.html")
// 	})
// 	// mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
// 	// 	http.ServeFile(w, r, "pages/index.html")
// 	// })

// 	// Handle the single-page application for logged-in users
// 	mux.HandleFunc("/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, "pages/index.html")
// 	}))

// 	// API handlers
// 	mux.HandleFunc("/api/login", server.LoginHandler)
// 	mux.HandleFunc("/api/posts", server.PostHandler)
// 	mux.HandleFunc("/api/posts/{id}/comments", server.CommentsHandler)
// 	mux.HandleFunc("/api/posts/{id}", server.GetPostHandler)
// 	mux.HandleFunc("/api/posts/{id}/dislike", server.AddDislikesHandler)
// 	mux.HandleFunc("/api/posts/{id}/like", server.AddLikesHandler)
// 	mux.HandleFunc("/api/signup", server.SignupHandler)
// 	mux.HandleFunc("/api/logout", server.LogoutHandler)
// 	mux.HandleFunc("/api/comments/{id}/like", server.LikeCommentHandler)
// 	mux.HandleFunc("/api/comments/{id}/dislike", server.DislikeCommentHandler)
// 	mux.HandleFunc("/api/check_session", server.CheckSessionHandler)
// 	mux.HandleFunc("/api/get_ages", server.GetAgesHandler)

// 	log.Println("Serving on http://localhost:8080")
// 	err := http.ListenAndServe(":8080", mux)
// 	if err != nil {
// 		log.Println("Error starting server:", err)
// 		os.Exit(1)
// 	}
// }

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

	// authMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		if r.URL.Path == "/api/login" || r.URL.Path == "/api/signup" {
	// 			next.ServeHTTP(w, r)
	// 			return
	// 		}
	// 		// Check if the user is authenticated
	// 		cookie, err := r.Cookie("session_token")
	// 		if err != nil || !db.IsSessionValid(cookie.Value) {
	// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 			return
	// 		}
	// 		next.ServeHTTP(w, r)
	// 	}
	// }

	// Handle login and registration separately
	// mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "pages/index.html")
	// })

	// Handle the single-page application for logged-in users
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "pages/index.html")
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
	mux.HandleFunc("/api/check_session", server.CheckSessionHandler)
	mux.HandleFunc("/api/get_ages", server.GetAgesHandler)

	log.Println("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("Error starting server:", err)
		os.Exit(1)
	}
}
