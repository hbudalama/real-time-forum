package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rtf/pkg/db"
	"rtf/pkg/structs"
	"strconv"
	"strings"
)

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
    if !LoginGuard(w, r) {
        http.Redirect(w, r, "/api/login", http.StatusTemporaryRedirect)
        return
    }

    postIDStr := r.URL.Path[len("/api/posts/") : len(r.URL.Path)-len("/comments")]
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    if r.Method == http.MethodGet {
        post, err := db.GetPost(postID)
        if err != nil {
            http.Error(w, `{"reason": "Error fetching post"}`, http.StatusInternalServerError)
            return
        }

        comments, err := db.GetComments(postID)
        if err != nil {
            http.Error(w, `{"reason": "Error fetching comments"}`, http.StatusInternalServerError)
            return
        }

        data := struct {
            Post     structs.Post      `json:"post"`
            Comments []structs.Comment `json:"comments"`
        }{
            Post:     post,
            Comments: comments,
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(data); err != nil {
            http.Error(w, `{"reason": "Error encoding comments"}`, http.StatusInternalServerError)
        }
        return
    }

	if r.Method == http.MethodPost {
		comment := r.FormValue("comment")
		if strings.TrimSpace(comment) == "" {
			http.Error(w, `{"success": false, "message": "Comment cannot be empty"}`, http.StatusBadRequest)
			return
		}
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Printf("can't get the cookie: %s\n", err.Error())
			http.Error(w, `{"success": false, "message": "Session expired"}`, http.StatusUnauthorized)
			return
		}
		token := cookie.Value
		session, err := db.GetSession(token)
		if err != nil {
			log.Printf("can't get the session: %s\n", err.Error())
			http.Error(w, `{"success": false, "message": "Session expired"}`, http.StatusUnauthorized)
			return
		}
		username := session.User.Username
		err = db.AddComment(postID, username, comment)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Error adding comment"}`, http.StatusInternalServerError)
			return
		}
	
		// Return updated comments
		post, err := db.GetPost(postID)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Error fetching post"}`, http.StatusInternalServerError)
			return
		}
		comments, err := db.GetComments(postID)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Error fetching comments"}`, http.StatusInternalServerError)
			return
		}
	
		data := struct {
			Success   bool             `json:"success"`
			Post      structs.Post     `json:"post"`
			Comments  []structs.Comment `json:"comments"`
		}{
			Success:  true,
			Post:     post,
			Comments: comments,
		}
	
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"success": false, "message": "Error encoding comments"}`, http.StatusInternalServerError)
		}
	}
}

func AddLikesHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusFound)
        return
    }
    postIDStr := r.PathValue("id")

    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var user structs.User
    cookie, err = r.Cookie("session_token")
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

    err = db.InsertOrUpdateInteraction(postID, user.Username, 1)
    if err != nil {
        http.Error(w, "Unable to like post", http.StatusInternalServerError)
        return
    }

    likes, dislikes, err := db.GetPostInteractions(postID) // here
    if err != nil {
        http.Error(w, "Unable to retrieve post interactions", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "success":  true,
        "likes":    likes,
        "dislikes": dislikes,
    }
    json.NewEncoder(w).Encode(response)
}

func AddDislikesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("i'm in dislike")

	if !LoginGuard(w, r) {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	postIDStr := r.PathValue("id")
	fmt.Printf("postIDStr: %s\n", postIDStr)
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		Error400Handler(w, r, "invalid post ID")
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("can't get the cookie: %s\n", err.Error())
		Error500Handler(w, r)
		return
	}

	token := cookie.Value
	session, err := db.GetSession(token)
	if err != nil {
		log.Printf("can't get the session: %s\n", err.Error())
		Error500Handler(w, r)
		return
	}
	username := session.User.Username

	err = db.InsertOrUpdateInteraction(postID, username, 0)
	if err != nil {
		Error500Handler(w, r)
		return
	}

	// Get the updated dislike count after the operation
	likes, dislikes, err := db.GetPostInteractions(postID)
	if err != nil {
		Error500Handler(w, r)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success":  true,
		"likes":    likes,
		"dislikes": dislikes,
	}
	json.NewEncoder(w).Encode(response)
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	commentIDStr := r.PathValue("id")

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var user structs.User
	cookie, err = r.Cookie("session_token")
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

	err = db.AddCommentInteraction(commentID, user.Username, 1)
	if err != nil {
		http.Error(w, "Unable to like comment", http.StatusInternalServerError)
		return
	}

	// http.Redirect(w, r, r.Referer(), http.StatusFound)
	likes, dislikes, err := db.GetCommentInteractions(commentID)
    if err != nil {
        http.Error(w, "Unable to retrieve comment interactions", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "success":  true,
        "likes":    likes,
        "dislikes": dislikes,
    }
    json.NewEncoder(w).Encode(response)
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	commentIDStr := r.PathValue("id")

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var user structs.User
	cookie, err = r.Cookie("session_token")
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

	err = db.AddCommentInteraction(commentID, user.Username, 0)
	if err != nil {
		http.Error(w, "Unable to dislike comment", http.StatusInternalServerError)
		return
	}

	// http.Redirect(w, r, r.Referer(), http.StatusFound)
	likes, dislikes, err := db.GetCommentInteractions(commentID)
    if err != nil {
        http.Error(w, "Unable to retrieve comment interactions", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "success":  true,
        "likes":    likes,
        "dislikes": dislikes,
    }
    json.NewEncoder(w).Encode(response)
}
