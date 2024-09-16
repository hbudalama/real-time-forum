package db
import (
	"log"
	"rtf/pkg/structs"
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
    log.Printf("Fetched chat history: %v", messages)
    if len(messages) == 0 {
        return nil, nil // Return nil if no messages found
    }
    return messages, nil
}
