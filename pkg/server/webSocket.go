package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"rtf/pkg/db"
	"github.com/gorilla/websocket"
)

const (
	messageTypeError          = "ERROR"
	messageTypeUserList       = "USER_LIST"
	messageTypeUnhandledEvent = "UNHANDLED_EVENT"
	messageTypeChatMessage	  = "CHAT_MESSAGE"
)

type Message struct {
	Type    string
	Payload any
}

type Users struct {
	Username string 
	Status string
}

type ChatMessage struct {
	Sender string
	Recipient string
	Content string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients []*websocket.Conn

func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I'm in echo")
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer connection.Close()

	clients = append(clients, connection)

	for {
		_, buffer, err := connection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			index := slices.Index(clients, connection)
			clients = slices.Delete(clients, index, index+1)
			break
		}
		var message Message
		errr := json.Unmarshal(buffer, &message)
		if errr != nil {
			log.Printf("WRONG PAYLOAD!!!!: %v", err)
			connection.WriteJSON(Message{Type: messageTypeError, Payload: "BAD REQUEST!"})
			continue
		}
		
		switch message.Type {
		case messageTypeUserList:
			userListHandler(connection)

		case messageTypeChatMessage:
			chatMessageHandler(connection)

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
        users[i] = Users{
            Username: username,
            // Assign status based on your application's logic
            Status: "online", // Default status, or customize based on your logic
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


func chatMessageHandler(conn *websocket.Conn) {
	chatMessages := []ChatMessage{
        {Sender: "haneen", Recipient: "fatema", Content: "hey pookie! are you coming to reboot today?"},
        {Sender: "fatema", Recipient: "haneen", Content: "hi pooks, Yes!"},
    }

    // Write JSON message with type CHAT_MESSAGE and payload as an array of chat messages
    message := Message{
        Type:    "CHAT_MESSAGE",
        Payload: chatMessages,
    }

    if err := conn.WriteJSON(message); err != nil {
        log.Printf("error writing chat message: %v", err)
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
