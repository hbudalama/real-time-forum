package db

import (
	"database/sql"
	"fmt"
	"log"
	"rtf/pkg/structs"
	"time"

	"github.com/google/uuid"
)

func CreateSession(username string) (string, error) {
	fmt.Println("im in session")
	token := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour)
	dbMutex.Lock()
	defer dbMutex.Unlock()
	_, err := db.Exec("UPDATE User SET sessionToken = ?, sessionExpiration = ? WHERE username = ?", token, expiry, username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetSession(token string) (*structs.Session, error) {
	session := structs.Session{}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	err := db.QueryRow("SELECT sessionToken, sessionExpiration, username FROM User WHERE sessionToken = ?", token).Scan(&session.Token, &session.Expiry, &session.User.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("GetSession: %s\n", err.Error())
		return nil, err
	}
	return &session, nil
}

func DeleteSession(token string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	_, err := db.Exec("UPDATE User SET sessionToken = NULL, sessionExpiration = NULL WHERE sessionToken = ?", token)
	return err
}

// func IsSessionValid(token string) bool {
// 	// Query your session table to check if the token exists and is valid
// 	return true // Return true if valid, false otherwise
// }

func IsSessionValid(token string) bool {
	session, err := GetSession(token)
	if err != nil || session == nil {
		return false
	}

	// Check if the session has expired
	if session.Expiry.Before(time.Now()) {
		return false
	}

	return true
}

func GetSessionByUsername(username string) (*structs.Session, error) {
	session := structs.Session{}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	err := db.QueryRow("SELECT sessionToken, sessionExpiration FROM User WHERE username = ?", username).Scan(&session.Token, &session.Expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("GetSessionByUsername: %s\n", err.Error())
		return nil, err
	}
	return &session, nil
}