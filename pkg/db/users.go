package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"github.com/mattn/go-sqlite3"
	"rtf/pkg/structs"
)

func AddUser(username string, email string, hashedPassword string) (*structs.User, error) {
	// Insert user into the database
	dbMutex.Lock()
	defer dbMutex.Unlock()
	_, err := db.Exec("INSERT INTO user (username, email, password) VALUES ($1, $2, $3)", username, email, hashedPassword)
	if err != nil {
		log.Printf("AddUser: %s\n", err.Error())
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			if strings.Contains(sqliteErr.Error(), "username") {
				return nil, errors.New("username already exists")
			} else if strings.Contains(sqliteErr.Error(), "email") {
				return nil, errors.New("email already exists")
			}
		}

		return &structs.User{}, errors.New("internal server error ??")
	}

	return nil, nil
}

func GetUser(username string) (*structs.User, error) {
	user := structs.User{}
	err := db.QueryRow("SELECT username, email FROM User WHERE username = ?", username).Scan(&user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("GetUser: %s\n", err.Error())
		return nil, err
	}

	return &user, nil
}
