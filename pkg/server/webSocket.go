package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rtf/pkg/db"
	"rtf/pkg/structs"
	"time"

	"github.com/gorilla/websocket"
)

const (
	messageTypeError          = "ERROR"
	messageTypeUserList       = "USER_LIST"
	messageTypeUnhandledEvent = "UNHANDLED_EVENT"
	messageTypeChatMessage    = "CHAT_MESSAGE"
	messageTypeChatHistory    = "CHAT_HISTORY"
	messageTypeNotification   = "NEW_MESSAGE_NOTIFICATION"
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
	CreatedDate string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn         *websocket.Conn
	SessionToken string
	Username     string
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

	// Get the session associated with the session token
	session, err := db.GetSession(cookie.Value)
	if err != nil || session == nil {
		log.Printf("Failed to get session by session token: %v", err)
		connection.Close()
		return
	}

	// Extract the username from the session
	loggedInUsername := session.User.Username

	// Store the connection along with the session token and username
	client := Client{
		Conn:         connection,
		SessionToken: cookie.Value,
		Username:     loggedInUsername,
	}
	clients = append(clients, client)

	// Send the updated user list to all clients (including the new connection)
	userListHandler()

	defer func() {
		// Remove the client from the list upon disconnect
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
			// Send the personalized user list to this client
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
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				chatHistoryHandler(connection, payload)
			} else {
				connection.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid payload for chat history request"})
			}

		case "TYPING_STATUS":
			// var typingStatus struct {
			// 	Sender    string `json:"Sender"`
			// 	Recipient string `json:"Recipient"`
			// 	IsTyping  bool   `json:"IsTyping"`
			// }

			// if err := json.Unmarshal(buffer, &typingStatus); err != nil {
			// 	log.Printf("Error parsing typing status: %v", err)
			// 	connection.WriteJSON(Message{Type: messageTypeError, Payload: "BAD REQUEST!"})
			// 	continue
			// }

			// // Broadcast typing status to the recipient (User B)
			// broadcastTypingStatus(typingStatus)
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

	// Create a message for each client with their personalized user list
	for _, client := range clients {
		// Get the session associated with the client
		session, err := db.GetSession(client.SessionToken)
		if err != nil || session == nil {
			log.Printf("Error getting session for client: %v", err)
			continue
		}

		loggedInUsername := session.User.Username

		// Build the user list, excluding the logged-in user
		users := make([]Users, 0)
		for _, username := range dbUsers {
			if username == loggedInUsername {
				continue // Skip the logged-in user
			}

			// Check if the user has a valid session
			session, err := db.GetSessionByUsername(username)
			status := "offline" // Default status

			if err == nil && session != nil && db.IsSessionValid(session.Token) {
				status = "online"
			}

			users = append(users, Users{
				Username: username,
				Status:   status,
			})
		}

		// Create a personalized message for this client
		message := Message{
			Type:    messageTypeUserList,
			Payload: users,
		}

		// Send the personalized user list to the client
		if err := client.Conn.WriteJSON(message); err != nil {
			log.Printf("Error writing user list to client %s: %v", loggedInUsername, err)
		}
	}
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

			// Send a notification message to the recipient (if needed)
			notificationMessage := Message{
				Type: "NEW_MESSAGE_NOTIFICATION",
				Payload: map[string]string{
					"Sender":  chatMsg.Sender,
					"Content": chatMsg.Content,
				},
			}
			if err := client.Conn.WriteJSON(notificationMessage); err != nil {
				log.Printf("Error sending notification to recipient: %v", err)
			}

			break }
	}

	createdDate := time.Now().Format("2006-01-02 15:04:05")

	chatMsg = ChatMessage {
		Sender: chatMsg.Sender,
		Recipient: chatMsg.Recipient,
		Content: chatMsg.Content,
		CreatedDate: createdDate,
	}

	// conn.WriteJSON(Message{Type: messageTypeChatMessage, Payload: chatMsg})

	// Optionally, send a confirmation back to the sender
	if err := conn.WriteJSON(Message{Type: messageTypeChatMessage, Payload: chatMsg}); err != nil {
		log.Printf("Error sending chat confirmation to sender: %v", err)
	}
}

func chatHistoryHandler(conn *websocket.Conn, chatRequest map[string]interface{}) {
	sender, ok := chatRequest["Sender"].(string)
	if !ok {
		log.Println("Invalid sender in chat history request")
		conn.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid sender in chat history request"})
		return
	}
	recipient, ok := chatRequest["Recipient"].(string)
	if !ok {
		log.Println("Invalid recipient in chat history request")
		conn.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid recipient in chat history request"})
		return
	}
	limit, ok := chatRequest["Limit"].(float64) // JSON numbers are parsed as float64
	if !ok {
		log.Println("Invalid limit in chat history request")
		conn.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid limit in chat history request"})
		return
	}
	offset, ok := chatRequest["Offset"].(float64)
	if !ok {
		log.Println("Invalid offset in chat history request")
		conn.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid offset in chat history request"})
		return
	}
	// Fetch chat history from the database
	messages, err := db.GetChatHistory(sender, recipient, int(limit), int(offset))
	if err != nil {
		log.Printf("Error fetching chat history: %v", err)
		conn.WriteJSON(Message{Type: messageTypeError, Payload: "Failed to fetch chat history"})
		return
	}
	if messages == nil {
		conn.WriteJSON(Message{Type: messageTypeChatHistory, Payload: []structs.ChatMessage{}})
		return
	}
	// Send the chat history back to the client
	conn.WriteJSON(Message{Type: messageTypeChatHistory, Payload: messages})
}
