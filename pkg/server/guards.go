package server

import (
	"log"
	"net/http"
	"strings"
	"time"
	"rtf/pkg/db"
)

func MethodsGuard(w http.ResponseWriter, r *http.Request, methods ...string) bool {
	method := strings.ToUpper(r.Method)
	log.Printf("Request method: %s\n", method)
	for _, v := range methods {
		if method == strings.ToUpper(v) {
			return true
		}
	}
	log.Printf("Method not allowed: %s\n", method)
	return false
}

func PostExistsGuard(w http.ResponseWriter, r *http.Request) bool {
	return true
}

func LoginGuard(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	token := cookie.Value
	session, err := db.GetSession(token)
	if err != nil || session == nil {
		// delete old cookie
		return false
	}
	// exist, err := db.CheckUsernameExists(session.User.Username)
	// if err != nil {
	// 	log.Printf("LoginGuard: Error checking username: %s\n", err.Error())
	// 	return
	// }
	// if exist && session.Expiry.Before(time.Now()){
	// 	log.Printf("LoginGuard: Error checking username: %s\n", err.Error())
	// }
	if session.Expiry.Before(time.Now()) {
		db.DeleteSession(token)
		return false
	}
	return true
}
