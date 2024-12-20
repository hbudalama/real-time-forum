package db

import (
	"fmt"
	"log"
	"rtf/pkg/structs"
	"time"
)

func SaveChatMessage(sender, recipient, content string) error {
	log.Printf("Saving message from sender: %s to recipient: %s with content: %s", sender, recipient, content)
	stmt, err := db.Prepare("INSERT INTO Chat (SenderUsername, RecipientUsername, Content) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sender, recipient, content)
	return err
}

func GetChatHistory(sender, recipient string, limit, offset int) ([]structs.ChatMessage, error) {
    log.Printf("SQL query: sender=%s, recipient=%s", sender, recipient)
    query := `
        SELECT SenderUsername, RecipientUsername, Content, CreatedDate
        FROM Chat
        WHERE (SenderUsername = ? AND RecipientUsername = ?) OR (SenderUsername = ? AND RecipientUsername = ?)
        ORDER BY CreatedDate DESC
        LIMIT ? OFFSET ?`
    rows, err := db.Query(query, sender, recipient, recipient, sender, limit, offset)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        return nil, err
    }
    defer rows.Close()
    var messages []structs.ChatMessage
    for rows.Next() {
        var message structs.ChatMessage
        if err := rows.Scan(&message.Sender, &message.Recipient, &message.Content, &message.CreatedDate); err != nil {
            return nil, err
        }
        messages = append(messages, message)
    }

    if len(messages) == 0 {
        return nil, nil // Return nil if no messages found
    }
    return messages, nil
}


func GetLastMessages(loggedInUsername string) ([]structs.LastMessage, error) {
    query := `
        SELECT CASE 
                  WHEN SenderUsername = ? THEN RecipientUsername 
                  ELSE SenderUsername 
              END AS Username, 
              Content, 
              MAX(CreatedDate) AS Timestamp
        FROM Chat
        WHERE SenderUsername = ? OR RecipientUsername = ?
        GROUP BY Username;
    `

    rows, err := db.Query(query, loggedInUsername, loggedInUsername, loggedInUsername) // Provide the logged-in user as a parameter
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var lastMessages []structs.LastMessage
    var lm structs.LastMessage
    for rows.Next() {
        // var lm structs.LastMessage
        var timestampStr string // Temporary string for scanning
        err := rows.Scan(&lm.Sender, &lm.Content, &timestampStr)
        if err != nil {
            return nil, err
        }

        // Parse the timestamp string to time.Time
        lm.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestampStr) // Adjust the format if necessary
        if err != nil {
            return nil, err
        }

        lastMessages = append(lastMessages, lm)
        fmt.Println("lm.Sender", lm.Sender)
        fmt.Println("lm.Content", lm.Content)
        fmt.Println(len(lastMessages))
    }
    return lastMessages, nil
}



