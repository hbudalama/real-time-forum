package db

import (
	"log"
	"rtf/pkg/structs"
)

// import (
// 	"fmt"
// 	"net/http"
// 	"rtf/pkg/structs"
// 	"github.com/gorilla/websocket"

// )

// // Map to store connected clients
// var clients = make(map[string]*websocket.Conn)

// // Channel to broadcast messages to all clients
// var broadcast = make(chan structs.Message)

// // Upgrader to upgrade HTTP connections to WebSocket connections
// var upgrader = websocket.Upgrader{
//     CheckOrigin: func(r *http.Request) bool {
//         return true
//     },
// }

// // handleConnections handles incoming WebSocket connections
// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	// Extract the session token from the request header
// 	token := r.Header.Get("Authorization")
// 	if token == "" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Retrieve the session from the database
// 	session, err := GetSession(token)
// 	if err != nil || session == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Extract the username from the session
// 	username := session.User.Username

// 	// Upgrade the HTTP connection to a WebSocket connection
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Upgrade:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Register the connection
// 	clients[username] = conn

// 	// Listen for incoming messages
// 	for {
// 		var msg structs.Message
// 		err := conn.ReadJSON(&msg)
// 		if err != nil {
// 			fmt.Println("ReadJSON:", err)
// 			delete(clients, username)
// 			break
// 		}
// 		msg.SenderUsername = username
// 		broadcast <- msg
// 		saveMessageToDB(msg) // Save the message to the database
// 	}
// }

// // handleMessages handles broadcasting messages to the appropriate recipient
// func handleMessages() {
// 	for {
// 		msg := <-broadcast // Receive a message from the broadcast channel

// 		// Find the recipient's connection
// 		recipientConn, ok := clients[msg.RecipientUsername]
// 		if ok {
// 			// Send the message to the recipient
// 			err := recipientConn.WriteJSON(msg)
// 			if err != nil {
// 				fmt.Println("WriteJSON:", err)
// 				recipientConn.Close()
// 				delete(clients, msg.RecipientUsername)
// 			}
// 		}
// 	}
// }

// // saveMessageToDB saves a message to the database
// func saveMessageToDB(msg structs.Message) {
// 	_, err := db.Exec("INSERT INTO PrivateMessages (SenderUsername, RecipientUsername, Content) VALUES (?, ?, ?)",
// 		msg.SenderUsername, msg.RecipientUsername, msg.Content)
// 	if err != nil {
// 		fmt.Println("Error saving message to DB:", err)
// 	}
// }

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

	rows, err := db.Query(`
        SELECT SenderUsername, RecipientUsername, Content, CreatedDate
        FROM Chat
        WHERE (SenderUsername = ? AND RecipientUsername = ?) OR (SenderUsername = ? AND RecipientUsername = ?)
        ORDER BY CreatedDate DESC
        LIMIT ? OFFSET ?, sender, recipient, recipient, sender, limit, offset`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	log.Println("Query executed successfully")

	defer rows.Close()

	var messages []structs.ChatMessage
	for rows.Next() {
		var message structs.ChatMessage
		var createdDate string
		if err := rows.Scan(&message.Sender, &message.Recipient, &message.Content, &createdDate); err != nil {
			log.Println("Query executed successfully")
			return nil, err
		}
		message.Content += " (" + createdDate + ")" // Append date to the content
		messages = append(messages, message)
	}

	log.Printf("Fetched chat history: %v", messages)

	if len(messages) == 0 {
		return nil, nil // Return nil if no messages found
	}

	return messages, nil
}
