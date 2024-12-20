package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"rtf/pkg/db"
	"rtf/pkg/structs"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || !db.IsSessionValid(cookie.Value) {
		json.NewEncoder(w).Encode(map[string]bool{"isAuthenticated": false})
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"isAuthenticated": true})
}

func GetAgesHandler(w http.ResponseWriter, r *http.Request) {
	var ages []int
	for i := 16; i <= 70; i++ {
		ages = append(ages, i)
	}
	json.NewEncoder(w).Encode(map[string][]int{"ages": ages})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !MethodsGuard(w, r, "GET", "POST") {
		return
	}

	if r.Method == "POST" {

		if err := r.ParseMultipartForm(1024); err != nil {
			http.Error(w, `{"reason": "Can't parse form data"}`, http.StatusBadRequest)
			Error400Handler(w, r, "Can't parse form data")
			return
		}

		identifier := strings.TrimSpace(r.Form.Get("username"))
		password := strings.TrimSpace(r.Form.Get("password"))

		fmt.Println(identifier)
		if identifier == "" {
			http.Error(w, `{reason: empty username}`, http.StatusBadRequest)
			Error400Handler(w, r, "Empty username")
			return
		}
		if password == "" {
			http.Error(w, `{reason: empty password}`, http.StatusBadRequest)
			Error400Handler(w, r, "Empty password")
			return
		}

		var username string
		var exists bool
		var err error

		if strings.Contains(identifier, "@") {
			// Check if it's an email
			exists, err = db.CheckEmailExists(identifier)
			if err != nil {
				http.Error(w, `{"reason": "Server error"}`, http.StatusInternalServerError)
				Error500Handler(w, r)
				return
			}
			if !exists {
				http.Error(w, `{"reason": "Email not found"}`, http.StatusNotFound)
				return
			}
			username, err = db.GetUsernameByEmail(identifier)
			if err != nil {
				http.Error(w, `{"reason": "Server error"}`, http.StatusInternalServerError)
				Error500Handler(w, r)
				return
			}
		} else {
			// Check if it's a username
			exists, err = db.CheckUsernameExists(identifier)
			if err != nil {
				http.Error(w, `{"reason": "Server error"}`, http.StatusInternalServerError)
				Error500Handler(w, r)
				return
			}
			if !exists {
				http.Error(w, `{"reason": "Username not found"}`, http.StatusNotFound)
				return
			}
			username = identifier
		}

		passwordMatches, err := db.CheckPassword(username, password)
		if err != nil {
			http.Error(w, `{"reason": "Server error"}`, http.StatusInternalServerError)
			Error500Handler(w, r)
			return
		}
		if !passwordMatches {
			http.Error(w, `{"reason": "Invalid password"}`, http.StatusUnauthorized)
			return
		}

		token, err := db.CreateSession(username)
		if err != nil {
			http.Error(w, `{"reason": "Server error"}`, http.StatusInternalServerError)
			Error500Handler(w, r)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Path:     "/",
			Secure:   false,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"isAuthenticated": true}`))
		return
	}

	var Ages []int
	for i := 16; i <= 70; i++ {
		Ages = append(Ages, i)
	}

	fmt.Println("i logged in")
	data := structs.PageData{Ages: Ages}
	tmpl := template.Must(template.ParseFiles(filepath.Join("pages", "index.html")))
	tmpl.Execute(w, data)
	fmt.Println("im in server")
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/api/posts" {
		// Get all posts from the database
		posts := db.GetAllPosts()

		// Set the response content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Encode the posts slice into JSON and write it to the response
		err := json.NewEncoder(w).Encode(posts)
		if err != nil {
			http.Error(w, `{reason: error encoding posts}`, http.StatusInternalServerError)
			return
		}
		return
	}

	// Existing code for handling individual posts by ID...
	postIDStr := r.URL.Path[len("/api/posts/"):]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, `{reason: invalid post id}`, http.StatusBadRequest)
		Error400Handler(w, r, "invalid post id")
		return
	}

	post, err := db.GetPost(postID)
	if err != nil {
		http.Error(w, `{reason: post not found}`, http.StatusNotFound)
		Error404Handler(w, r)
		return
	}
	comments, err := db.GetComments(postID)
	if err != nil {
		http.Error(w, `{reason: invalid comment}`, http.StatusInternalServerError)
		Error500Handler(w, r)
		return
	}
	var user *structs.User
	cookie, err := r.Cookie("session_token")
	if err == nil {
		token := cookie.Value
		session, err := db.GetSession(token)
		if err == nil {
			u := session.User
			user = &u
		}
	}
	data := struct {
		Post         structs.Post
		Comments     []structs.Comment
		LoggedInUser *structs.User
	}{
		Post:         post,
		Comments:     comments,
		LoggedInUser: user,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		Error500Handler(w, r)
		return
	}
	tmpl.Execute(w, data)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		postID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		post, err := db.GetPost(postID)
		if err != nil {
			Error404Handler(w, r)
			return
		}

		comments, err := db.GetComments(postID)
		if err != nil {
			Error500Handler(w, r)
			return
		}

		var user *structs.User
		if LoginGuard(w, r) {
			cookie, err := r.Cookie("session_token")
			if err == nil {
				token := cookie.Value
				session, err := db.GetSession(token)
				if err == nil {
					user = &session.User
				}
			}
		}

		ctx := struct {
			Post         structs.Post
			Comments     []structs.Comment
			LoggedInUser *structs.User
		}{
			Post:         post,
			Comments:     comments,
			LoggedInUser: user,
		}

		tmpl, err := template.ParseFiles(filepath.Join("pages", "index.html"))
		if err != nil {
			log.Printf("can't parse the template: %s\n", err.Error())
			Error500Handler(w, r)
			return
		}

		err = tmpl.Execute(w, ctx)
		if err != nil {
			log.Printf("can't execute the template: %s\n", err.Error())
			Error500Handler(w, r)
			return
		}
	} else if r.Method == http.MethodPost {
		// This part handles creating a new post
		if !LoginGuard(w, r) {
			http.Error(w, `{"reason": Unauthorized}`, http.StatusUnauthorized)
			return
		}

		// Parse form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["post-category"]

		if title == "" || content == "" {
			http.Error(w, `{"reason:": Missing post title or content}`, http.StatusBadRequest)
			Error400Handler(w, r, "Missing post title or content")
			return
		}

		// Retrieve the current user
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, `{"reason": Failed to retrieve session token}`, http.StatusInternalServerError)
			Error500Handler(w, r)
			return
		}
		token := cookie.Value
		session, err := db.GetSession(token)
		if err != nil {
			http.Error(w, `{"reason": Invalid session}`, http.StatusUnauthorized)
			return
		}

		// Use the existing CreatePost function
		postID, err := db.CreatePost(title, content, session.User.Username, categories[0])
		if err != nil {
			http.Error(w, `{"reason": Failed to create post"}`, http.StatusInternalServerError)
			Error500Handler(w, r)
			return
		}

		// // Save post categories
		// err = db.AddPostCategories(postID, categories)
		// if err != nil {
		// 	http.Error(w, `{"reason": Failed to save post categories}`, http.StatusInternalServerError)
		// 	Error500Handler(w, r)
		// 	return
		// }

		// Redirect to the newly created post page or return a success response
		http.Redirect(w, r, fmt.Sprintf("/api/posts/%d", postID), http.StatusSeeOther)
	} else {
		http.Error(w, `{"reason": Method not allowed}`, http.StatusMethodNotAllowed)
		Error404Handler(w, r)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if !MethodsGuard(w, r, "POST") {
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, `{"reason": "cant parse signup form"}`, http.StatusBadRequest)
		Error400Handler(w, r, "Can't parse form data")
		return
	}

	log.Println(r.Form)

	username := strings.TrimSpace(r.Form.Get("username"))
	firstName := strings.TrimSpace(r.Form.Get("firstName"))
	lastName := strings.TrimSpace(r.Form.Get("lastName"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))
	confirmPassword := strings.TrimSpace(r.Form.Get("password"))

	gender := r.Form.Get("gender") == "1"
	// true :=  1 == 1
	// false:=0 == 1

	age, err := strconv.Atoi(r.Form.Get("age"))
	if err != nil {
		log.Printf("Can't parse date for signup: %s\n", err)
		http.Error(w, `{"reason": "Invalid age format"}`, http.StatusBadRequest)
		return
	}

	if username == "" {
		http.Error(w, `{"reason": "Username is required"}`, http.StatusBadRequest)
		return
	}
	if email == "" {
		http.Error(w, `{"reason": "Email is required"}`, http.StatusBadRequest)
		return
	}
	if !validEmail(email) {
		http.Error(w, `{"reason": "Invalid email format"}`, http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, `{"reason": "Password is required"}`, http.StatusBadRequest)
		return
	}
	if !validatePassword(password) {
		http.Error(w, `{"reason": "Password must be at least 8 characters long"}`, http.StatusBadRequest)
		return
	}
	if password != confirmPassword {
		http.Error(w, `{"reason": "Passwords do not match"}`, http.StatusBadRequest)
		return
	}

	exists, err := db.CheckUsernameExists(username)
	if err != nil {
		Error500Handler(w, r)
		return
	}
	if exists {
		http.Error(w, `{"reason": "Username already taken"}`, http.StatusBadRequest)
		return
	}

	emailExists, err := db.CheckEmailExists(email)
	if err != nil {
		Error500Handler(w, r)
		return
	}
	if emailExists {
		http.Error(w, `{"reason": "Email already taken"}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		Error500Handler(w, r)
		return
	}
	_, err = db.AddUser(username, firstName, lastName, gender, age, email, string(hashedPassword))
	if err != nil {
		errMsg := fmt.Sprintf(`{"reason": "%s"}`, err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error404Handler(w, r)
		return
	}

	if !MethodsGuard(w, r, "GET") {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := structs.HomeContext{}
	if LoginGuard(w, r) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Printf("can't get the cookie: %s\n", err.Error())
			return
		}

		token := cookie.Value

		session, err := db.GetSession(token)
		if err != nil {
			log.Printf("can't get the session: %s\n", err.Error())
			return
		}
		ctx.LoggedInUser = &session.User
	}

	ctx.Posts = db.GetAllPosts()

	users, err := db.GetAllUsers()
	if err != nil {
		log.Printf("error getting users %s", err.Error())
		Error500Handler(w, r)
		return
	}
	ctx.Users = users

	tmpl, err := template.ParseFiles(filepath.Join("pages", "index.html"))
	if err != nil {
		log.Printf("can't parse the template here: %s\n", err.Error())
		Error500Handler(w, r)
		return
	}

	fmt.Println(ctx.Posts)

	err = tmpl.Execute(w, ctx)
	if err != nil {
		log.Printf("can't execute the template: %s\n", err.Error())
		Error500Handler(w, r)
		return
	}
}

func AddPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !LoginGuard(w, r) {
		http.Redirect(w, r, "/api/login", http.StatusTemporaryRedirect)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.FormValue("post-category") // Update to match the form field name

	log.Printf("Received categories: %v\n", categories) // Debug print

	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" || len(categories) == 0 {
		http.Error(w, `{"reason": "The post must have a title, content, and category"}`, http.StatusBadRequest)
		return
	}

	var user structs.User
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, `{"reason": "can't get the cookie"}`, http.StatusBadRequest)
		return
	}

	token := cookie.Value

	session, err := db.GetSession(token)
	if err != nil {
		http.Error(w, `{"reason": "can't get the session"}`, http.StatusBadRequest)
		return
	}
	user = session.User

	postID, err := db.CreatePost(title, content, user.Username, categories)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `{"reason": "failed to create post"}`, http.StatusInternalServerError)
		Error500Handler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "postID": postID})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !MethodsGuard(w, r, "GET") {
		http.Error(w, "Only GET requests allowed", http.StatusMethodNotAllowed)
		return
	}

	// Attempt to get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Attempt to delete the session in the database
		err = db.DeleteSession(cookie.Value)
		if err != nil {
			log.Printf("LogoutHandler: %s", err.Error())
			Error500Handler(w, r)
			return
		}
	} else if err != http.ErrNoCookie {
		// If there's an error other than no cookie, log it
		log.Printf("can't get the cookie: %s\n", err.Error())
		Error500Handler(w, r)
		return
	}

	// Clear the session cookie regardless of whether it was found
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})

	fmt.Println("i logged out")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func Error400Handler(w http.ResponseWriter, r *http.Request, reason string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"reason": "` + reason + `"}`))
}

func Error500Handler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	// w.Write([]byte(`{"reason": "Interval Server Error"}`))
}

func Error404Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"reason": "Not found"}`))
}

func validEmail(email string) bool {
	return strings.Contains(email, "@") && strings.HasSuffix(email, ".com")
}

func validatePassword(password string) bool {
	return len(password) >= 8
}

func RenderAddPostForm(w http.ResponseWriter, r *http.Request, errorMessage string) {
	var user structs.User
	if LoginGuard(w, r) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Printf("can't get the cookie: %s\n", err.Error())
			return
		}

		token := cookie.Value

		session, err := db.GetSession(token)
		if err != nil {
			log.Printf("can't get the session: %s\n", err.Error())
			return
		}
		user = session.User
	}

	ctx := structs.HomeContext{
		LoggedInUser: &user,
		ErrorMessage: errorMessage,
		Posts:        db.GetAllPosts(),
	}

	tmpl, err := template.ParseFiles(filepath.Join("pages", "index.html"))
	if err != nil {
		log.Printf("can't parse the template: %s\n", err.Error())
		Error500Handler(w, r)
		return
	}

	err = tmpl.Execute(w, ctx)
	if err != nil {
		log.Printf("can't execute the template: %s\n", err.Error())
		Error500Handler(w, r)
		return
	}
}

func PostAPIHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Path[len("/api/posts/"):] // Extract the post ID from the URL
	postID, err := strconv.Atoi(postIDStr)       // Convert the post ID to an integer
	if err != nil {                              // Handle error if the post ID is not a valid integer
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := db.GetPost(postID) // Retrieve the post details from the database using the post ID
	if err != nil {                 // Handle error if the post is not found
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the response content type to JSON
	json.NewEncoder(w).Encode(post)                    // Encode the post data as JSON and write it to the response
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve session information from the database
	session, err := db.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	// Fetch the user info from the session
	user := session.User

	// Respond with the user information in JSON format
	json.NewEncoder(w).Encode(map[string]interface{}{
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"age":       user.Age,
		// Add any other fields you want to return
	})
}

func GetAllUsernamesHandler(w http.ResponseWriter, r *http.Request) {
    users, err := db.GetAllUsernames()
    if err != nil {
        http.Error(w, `{"reason": "Error fetching usernames"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}


