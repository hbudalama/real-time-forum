package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rtf/pkg/db"

	"github.com/gorilla/websocket"
)

const (
	messageTypeError          = "ERROR"
	messageTypeUserList       = "USER_LIST"
	messageTypeUnhandledEvent = "UNHANDLED_EVENT"
	messageTypeChatMessage    = "CHAT_MESSAGE"
)

type Message struct {
	Type    string
	Payload any
}

type Users struct {
	Username string
	Status   string
}

type ChatMessage struct {
	Sender    string
	Recipient string
	Content   string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn         *websocket.Conn
	SessionToken string
}

var clients []Client

func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I'm in echo")
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer connection.Close()

	// Retrieve session token from the HTTP request (e.g., from cookies)
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("No session cookie found")
		connection.Close()
		return
	}

	// Store the connection along with the session token
	client := Client{
		Conn:         connection,
		SessionToken: cookie.Value,
	}
	clients = append(clients, client)

	defer func() {
		// Remove client from the list upon disconnect
		for i, c := range clients {
			if c.Conn == connection {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		connection.Close()
	}()

	for {
		_, buffer, err := connection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			// Remove the client from the slice on error (disconnection)
			for i, c := range clients {
				if c.Conn == connection {
					clients = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			break
		}

		var message Message
		if err := json.Unmarshal(buffer, &message); err != nil {
			log.Printf("WRONG PAYLOAD!!!!: %v", err)
			connection.WriteJSON(Message{Type: messageTypeError, Payload: "BAD REQUEST!"})
			continue
		}

		switch message.Type {
		case messageTypeUserList:
			userListHandler(connection)

		case messageTypeChatMessage:
			// Extract the chat message from the payload
			var chatMessage ChatMessage
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				chatMessage.Sender = payload["Sender"].(string)
				chatMessage.Recipient = payload["Recipient"].(string)
				chatMessage.Content = payload["Content"].(string)
			}
			chatMessageHandler(connection, chatMessage)

		default:
			connection.WriteJSON(Message{Type: messageTypeUnhandledEvent, Payload: fmt.Sprintf("[%s] is not handled", message.Type)})
		}
	}
}

func userListHandler(conn *websocket.Conn) {
	dbUsers, err := db.GetAllUsernames() // Assuming this returns a slice of usernames
	if err != nil {
		log.Printf("error getting users list: %v", err)
		return
	}

	// Creating a list of User structs with some statuses
	users := make([]Users, len(dbUsers))
	for i, username := range dbUsers {
		// Check if the user has a valid session
		session, err := db.GetSessionByUsername(username)
		status := "offline" // Default status

		if err == nil && session != nil && db.IsSessionValid(session.Token) {
			status = "online"
		}

		users[i] = Users{
			Username: username,
			Status:   status,
		}
	}

	// Send the users list as JSON
	message := Message{
		Type:    "USER_LIST",
		Payload: users,
	}

	if err := conn.WriteJSON(message); err != nil {
		log.Printf("error writing users list: %v", err)
	}

	fmt.Println(users)
}

// func chatMessageHandler(conn *websocket.Conn) {
// 	chatMessages := []ChatMessage{
//         {Sender: "haneen", Recipient: "fatema", Content: "hey pookie! are you coming to reboot today?"},
//         {Sender: "fatema", Recipient: "haneen", Content: "hi pooks, Yes!"},
//     }

//     // Write JSON message with type CHAT_MESSAGE and payload as an array of chat messages
//     message := Message{
//         Type:    "CHAT_MESSAGE",
//         Payload: chatMessages,
//     }

//     if err := conn.WriteJSON(message); err != nil {
//         log.Printf("error writing chat message: %v", err)
//     }
// }

func chatMessageHandler(conn *websocket.Conn, chatMsg ChatMessage) {
	// Save the chat message to the database
	err := db.SaveChatMessage(chatMsg.Sender, chatMsg.Recipient, chatMsg.Content)
	if err != nil {
		log.Printf("Error saving chat message: %v", err)
		conn.WriteJSON(Message{Type: messageTypeError, Payload: "Failed to save message"})
		return
	}

	// Iterate through all connections and find the recipient
	for _, client := range clients {
		// Get session from the token
		session, err := db.GetSession(client.SessionToken)
		if err != nil || session == nil {
			continue
		}

		// Check if the session username matches the recipient
		if session.User.Username == chatMsg.Recipient {
			// Send the chat message to the recipient
			if err := client.Conn.WriteJSON(Message{Type: messageTypeChatMessage, Payload: chatMsg}); err != nil {
				log.Printf("Error sending chat message to recipient: %v", err)
			}
		}
	}

	// Optionally, you can also send a confirmation message back to the sender
	if err := conn.WriteJSON(Message{Type: messageTypeChatMessage, Payload: chatMsg}); err != nil {
		log.Printf("Error sending chat confirmation to sender: %v", err)
	}
}

func MessageHandler(message []byte) {
	fmt.Println("this is message handler: " + string(message))
}

//--broadcasting--
// func WriteMessage(message []byte) {
// 	for conn := range clients {
// 		fmt.Println("this is write message:")
// 		err := conn.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			log.Printf("Error writing message: %v", err)
// 			conn.Close()
// 			delete(clients, conn) // Remove the client from the map
// 		}
// 	}
// }
