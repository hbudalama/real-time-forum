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
	messageTypeChatHistory    = "CHAT_HISTORY"
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

	// Send the updated user list to all clients after a new connection
	userListHandler()

	defer func() {
		// Remove client from the list upon disconnect
		for i, c := range clients {
			if c.Conn == connection {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		connection.Close()

		// Send the updated user list to all clients after a disconnection
		userListHandler()
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
			// No need to pass the connection; send the list to all clients
			userListHandler()

		case messageTypeChatMessage:
			// Extract the chat message from the payload
			var chatMessage ChatMessage
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				chatMessage.Sender = payload["Sender"].(string)
				chatMessage.Recipient = payload["Recipient"].(string)
				chatMessage.Content = payload["Content"].(string)
			}
			chatMessageHandler(connection, chatMessage)

		case messageTypeChatHistory:
			//complete here
			var historyRequest struct {
                Sender       string 
                Recipient    string
                Limit        int
                Offset       int
            }
            if payload, ok := message.Payload.(map[string]interface{}); ok {
                historyRequest.Sender = payload["Sender"].(string)
                historyRequest.Recipient = payload["Recipient"].(string)
                historyRequest.Limit = int(payload["Limit"].(float64))
                historyRequest.Offset = int(payload["Offset"].(float64))
            }
			// chatHistoryHandler(connection, historyRequest)

		default:
			connection.WriteJSON(Message{Type: messageTypeUnhandledEvent, Payload: fmt.Sprintf("[%s] is not handled", message.Type)})
		}
	}
}

func userListHandler() {
	dbUsers, err := db.GetAllUsernames() // Assuming this returns a slice of usernames
	if err != nil {
		log.Printf("Error getting users list: %v", err)
		return
	}

	// Creating a list of User structs with their statuses
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

	// Create the message to send to all clients
	message := Message{
		Type:    "USER_LIST",
		Payload: users,
	}

	// Send the users list to all connected clients
	for _, client := range clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			log.Printf("Error writing users list to client: %v", err)
		}
	}

	fmt.Println(users)
}

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

// chatHistoryHandler(conn *websocket.Conn, chatMsg structs.chatMessage) {
// complete later
// }

func MessageHandler(message []byte) {
	fmt.Println("this is message handler: " + string(message))
}

