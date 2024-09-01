package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func Echo(w http.ResponseWriter, r *http.Request) {
    fmt.Println("I'm in echo")
    connection, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade to WebSocket: %v", err)
        return
    }
    defer connection.Close()

    clients[connection] = true

    for {
        _, message, err := connection.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            delete(clients, connection)  // Remove client from the map if there's an error
            break
        }

        fmt.Println("Received message: ", string(message)) // Log the received message
        
        // Echo the message back to the client
        err = connection.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Printf("Error writing message: %v", err)
            break
        }
    }
}



func MessageHandler(message []byte) {
	fmt.Println("this is message handler: " + string(message))
}

func WriteMessage(message []byte) {
    for conn := range clients {
        fmt.Println("this is write message:")
        err := conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Printf("Error writing message: %v", err)
            conn.Close()
            delete(clients, conn)  // Remove the client from the map
        }
    }
}
