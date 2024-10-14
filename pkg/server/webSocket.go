package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rtf/pkg/db"
	"rtf/pkg/structs"
	"sort"
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
	messageTypeTypingStatus   = "TYPING_STATUS"
	messageTypeChatOpened     = "CHAT_OPENED"
	messageTypeChatClosed     = "CHAT_CLOSED"
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
	Sender      string
	Recipient   string
	Content     string
	CreatedDate string
}

// type LastMessage struct {
// 	Sender    string    // The sender of the last message
// 	Content   string    // Content of the last message
// 	Timestamp time.Time // When the last message was sent
// }

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn         *websocket.Conn
	SessionToken string
	Username     string
	IsOnline     bool // Track if the user is currently online
}

var clients []Client

// Global map to track active chat sessions
var activeChats = make(map[string]map[string]bool) // map[sender]map[recipient]bool

func Echo(w http.ResponseWriter, r *http.Request) {
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
		IsOnline:     true, // Set status to online when connected
	}
	clients = append(clients, client)

	// Send the updated user list to all clients (including the new connection)
	userListHandler()

	defer func() {
		// Remove the client from the list upon disconnect
		for i, c := range clients {
			if c.Conn == connection {
				// clients = append(clients[:i], clients[i+1:]...)
				clients[i].IsOnline = false
				break
			}
		}

		// Send the updated user list to all clients after a disconnection
		userListHandler()

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
			// Send the personalized user list to this client
			fmt.Println("here!!")
			userListHandler()

		case messageTypeChatMessage:
			// Extract the chat message from the payload with checks
			var chatMessage ChatMessage
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				// Check if Sender is a valid string
				if sender, ok := payload["Sender"].(string); ok && sender != "" {
					chatMessage.Sender = sender
				} else {
					log.Printf("Invalid or missing sender in payload")
					connection.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid or missing sender"})
					continue
				}

				// Check if Recipient is a valid string
				if recipient, ok := payload["Recipient"].(string); ok && recipient != "" {
					chatMessage.Recipient = recipient
				} else {
					log.Printf("Invalid or missing recipient in payload")
					connection.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid or missing recipient"})
					continue
				}

				// Check if Content is a valid string
				if content, ok := payload["Content"].(string); ok && content != "" {
					chatMessage.Content = content
				} else {
					log.Printf("Invalid or missing content in payload")
					connection.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid or missing content"})
					continue
				}

				// Call the handler for sending the chat message
				chatMessageHandler(connection, chatMessage)
			}
		case messageTypeChatHistory:
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				chatHistoryHandler(connection, payload)
			} else {
				connection.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid payload for chat history request"})
			}

		case messageTypeTypingStatus:
			// Extract the payload from the message
			var typingStatus structs.TypingStatus
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				// Parse the payload into the TypingStatus struct
				typingStatus.Sender = payload["Sender"].(string)
				typingStatus.Recipient = payload["Recipient"].(string)
				typingStatus.IsTyping = payload["IsTyping"].(bool)
			} else {
				log.Printf("Invalid payload for typing status")
				connection.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid payload for typing status"})
				continue
			}

			fmt.Printf("Typing status: %+v\n", typingStatus)

			// Broadcast typing status to the recipient (User B)
			broadcastTypingStatus(typingStatus)

		case messageTypeChatOpened:
			// Handle chat opened
			if payload, ok := message.Payload.(map[string]interface{}); ok {
				chatOpenedHandler(payload)
			} else {
				log.Printf("Invalid payload for chat opened")
				connection.WriteJSON(Message{Type: messageTypeError, Payload: "Invalid payload for chat opened"})
			}
		case messageTypeChatClosed:
			chatClosedHandler(connection)
		default:
			connection.WriteJSON(Message{Type: messageTypeUnhandledEvent, Payload: fmt.Sprintf("[%s] is not handled", message.Type)})
		}
	}
}

func userListHandler() {
	dbUsers, err := db.GetAllUsernames() // Fetch all usernames
	if err != nil {
		log.Printf("Error getting users list: %v", err)
		return
	}

	// Fetch the latest messages for each user from the Chat table
	lastMessages, err := db.GetLastMessages() // Custom function to get last messages
	if err != nil {
		log.Printf("Error getting last messages: %v", err)
		return
	}

	// Create a map of username to LastMessage struct
	if lastMessages == nil {
		log.Printf("there is no last message")

	}

	userLastMessage := make(map[string]*structs.LastMessage)
	for _, msg := range lastMessages {
		userLastMessage[msg.Sender] = &msg
	}

	// Create a map to track online status
	onlineStatus := make(map[string]bool)

	// Mark users as online based on active WebSocket connections
	for _, client := range clients {
		if client.IsOnline {
			onlineStatus[client.Username] = true
		}
	}

	// Send the user list to all clients
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

			// Determine the user's status based on the onlineStatus map
			status := "offline"
			if onlineStatus[username] {
				status = "online"
			}

			// Add the user to the list, even if there are no messages
			users = append(users, Users{
				Username: username,
				Status:   status,
			})
		}
		fmt.Println("before sorting: ", users)
		// Sort users: prioritize based on last message, then alphabetically
		sort.Slice(users, func(i, j int) bool {
			// Check if both users have messages
			msgI, okI := userLastMessage[users[i].Username]
			msgJ, okJ := userLastMessage[users[j].Username]

			if okI && okJ {
				// Compare timestamps (latest first)
				return msgI.Timestamp.After(msgJ.Timestamp)
			} else if okI {
				// If only user i has messages, they go first
				return true
			} else if okJ {
				// If only user j has messages, they go first
				return false
			} else {
				// If neither has messages, sort alphabetically
				return users[i].Username < users[j].Username
			}
		})

		// Send updated user list to the specific client
		message := Message{
			Type:    messageTypeUserList,
			Payload: users,
		}
		fmt.Println("after sorting: ", users)
		err = client.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending user list to client: %v", err)
		}
	}
}

func chatMessageHandler(conn *websocket.Conn, chatMsg ChatMessage) {
	chatMsg.CreatedDate = time.Now().Format(time.RFC3339)
	// Save the chat message to the database
	err := db.SaveChatMessage(chatMsg.Sender, chatMsg.Recipient, chatMsg.Content)
	if err != nil {
		log.Printf("Error saving chat message: %v", err)
		conn.WriteJSON(Message{Type: messageTypeError, Payload: "Failed to save message"})
		return
	}

	fmt.Printf("Sending from %s to %s\n", chatMsg.Sender, chatMsg.Recipient)

	// Check if the recipient is in an active chat with the sender
	if activeChats[chatMsg.Sender] == nil || !activeChats[chatMsg.Sender][chatMsg.Recipient] {
		log.Printf("Recipient %s is not in an active chat with sender %s. Message will not be sent.", chatMsg.Recipient, chatMsg.Sender)
		return // Don't send the message if the recipient isn't in the active chat
	}

	// Iterate through all connections and find both the recipient and the sender
	for _, client := range clients {
		session, err := db.GetSession(client.SessionToken)
		if err != nil || session == nil {
			continue
		}

		// Notify the recipient
		if session.User.Username == chatMsg.Recipient {
			log.Println("Sending message to recipient...")

			// Send the chat message to the recipient
			if err := client.Conn.WriteJSON(Message{Type: messageTypeChatMessage, Payload: chatMsg}); err != nil {
				log.Printf("Error sending chat message to recipient: %v", err)
			}

			// Send a notification to the recipient
			notificationMessage := Message{
				Type: messageTypeNotification,
				Payload: map[string]string{
					"Sender":  chatMsg.Sender,
					"Content": chatMsg.Content,
				},
			}
			if err := client.Conn.WriteJSON(notificationMessage); err != nil {
				log.Printf("Error sending notification to recipient: %v", err)
			}
		}

		// Optionally, notify the sender (if required)
		if session.User.Username == chatMsg.Sender {
			log.Println("Sending confirmation to sender...")

			// Send a confirmation back to the sender
			if err := conn.WriteJSON(Message{Type: messageTypeChatMessage, Payload: chatMsg}); err != nil {
				log.Printf("Error sending chat confirmation to sender: %v", err)
			}
		}
	}

	log.Printf("Message sent with CreatedDate: %s", chatMsg.CreatedDate)
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

func broadcastTypingStatus(typingStatus structs.TypingStatus) {
	for _, client := range clients {
		session, err := db.GetSession(client.SessionToken)
		if err != nil || session == nil {
			continue
		}
		fmt.Printf("Comparing sender: %s with recipient: %s, Equal: %t\n", session.User.Username, typingStatus.Recipient, session.User.Username == typingStatus.Recipient)

		// Check if the recipient is the logged-in user
		if session.User.Username == typingStatus.Recipient {
			// Create a message to notify the recipient about the typing status
			message := Message{
				Type: "TYPING_STATUS",
				Payload: map[string]interface{}{
					"Sender":    typingStatus.Sender,
					"Recipient": typingStatus.Recipient,
					"IsTyping":  typingStatus.IsTyping,
				},
			}
			log.Printf("Broadcasting typing status to %s: %+v", typingStatus.Recipient, message) // Debug log
			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("Error broadcasting typing status: %v", err)
			}
			break
		}
	}
}

func chatOpenedHandler(payload map[string]interface{}) {
	// Extract sender and recipient from the payload
	sender, ok1 := payload["Sender"].(string)
	recipient, ok2 := payload["Recipient"].(string)
	if !ok1 || !ok2 {
		log.Printf("Invalid payload for chat opened: %v", payload)
		return
	}

	// Check if a chat session already exists
	if _, exists := activeChats[sender]; !exists {
		activeChats[sender] = make(map[string]bool)
	}

	// Mark the recipient as having their chat opened
	activeChats[sender][recipient] = true

	// Optionally, send a notification or log that the chat has been opened
	fmt.Printf("Chat opened between %s and %s", sender, recipient)
}

func chatClosedHandler(connection *websocket.Conn) {
	var sender string

	// Identify the sender based on the WebSocket connection
	for _, client := range clients {
		if client.Conn == connection {
			sender = client.Username
			break
		}
	}

	if sender == "" {
		log.Printf("Sender could not be identified for closed chat")
		return
	}

	log.Printf("Chat closed by %s", sender)

	// Notify the other users or the recipient that this user closed the chat
	for _, client := range clients {
		if client.Username != sender {
			message := Message{
				Type: messageTypeChatClosed,
				Payload: map[string]string{
					"Sender": sender,
				},
			}
			client.Conn.WriteJSON(message)
		}
	}
}
