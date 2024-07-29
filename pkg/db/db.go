package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var dbMutex = sync.Mutex{}

func Connect() error {
	database, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return err
	}
	sqlFile, err := os.ReadFile("./sql/createTables.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}
	_, err = database.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("failed to execute SQL commands: %v", err)
	}
	db = database
	return nil
}

func Close() error {
	return db.Close()
}

func CheckUsernameExists(username string) (bool, error) {
	var foundUsername string
	dbMutex.Lock()
	defer dbMutex.Unlock()
	err := db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&foundUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CheckPassword(username, password string) (bool, error) {
	hashedPassword, err := GetHashedPasswordByUsername(username)
	if err != nil {
		return false, err
	}
	if hashedPassword == "" {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetHashedPasswordByUsername(username string) (string, error) {
	var hashedPassword string
	dbMutex.Lock()
	defer dbMutex.Unlock()
	err := db.QueryRow("SELECT password FROM user WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return hashedPassword, nil
}

func CheckEmailExists(email string) (bool, error) {
    var foundEmail string
    dbMutex.Lock()
    defer dbMutex.Unlock()
    
    email = strings.TrimSpace(email) // Trim whitespace

    log.Printf("Checking email existence for: %s", email) // Debugging line

    err := db.QueryRow("SELECT email FROM user WHERE email = ?", email).Scan(&foundEmail)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        log.Printf("Error querying email: %v", err) // Debugging line
        return false, err
    }

    log.Printf("Email found: %s", foundEmail) // Debugging line
    return true, nil
}



func GetUsernameByEmail(email string) (string, error) {
	var username string
	dbMutex.Lock()
	defer dbMutex.Unlock()
	err := db.QueryRow("SELECT username FROM user WHERE email = ?", email).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return username, nil
}