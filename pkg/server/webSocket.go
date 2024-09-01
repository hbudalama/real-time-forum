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
)

type Message struct {
	Type    string
	Payload any
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
		// fmt.Println("Received message: ", message.Type) // Log the received message

		// // Echo the message back to the client
		// err = connection.WriteJSON(message)
		// if err != nil {
		// 	log.Printf("Error writing message: %v", err)
		// 	break
		// }

		switch message.Type {
		case messageTypeUserList:
			userListHandler(connection)

		default:
			connection.WriteJSON(Message{Type: messageTypeUnhandledEvent, Payload: fmt.Sprintf("[%s] is not handled", message.Type)})
		}
	}
}

func userListHandler(conn *websocket.Conn) {
	users, err := db.GetAllUsernames()
	if err != nil {
		log.Printf("error geting users list: %v", err)
	}

	// write json message {type:USERS_LIST payload: array of users}
	conn.WriteJSON(Message{Type: messageTypeUserList, Payload: users})
	fmt.Println(users)
	
}

func MessageHandler(message []byte) {
	fmt.Println("this is message handler: " + string(message))
}

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
