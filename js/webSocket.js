import { appendChatMessage } from './chat.js';
import { prependChatMessages } from './chat.js';

export function initializeWebSocket() {
    window.socket = new WebSocket("ws://localhost:8080/api/ws");

    window.socket.onopen = function(event) {
        console.log("Connected to WebSocket server");

        // Send a message to request the users list
        const requestMessage = JSON.stringify({
            Type: "USER_LIST",
            Payload: null
        });
        window.socket.send(requestMessage);
    };
    
    window.socket.onmessage = function(event) {
        try {
            const message = JSON.parse(event.data);
    
            if (message.Type === "USER_LIST") {
                console.log("Received user list: ", message.Payload);
    
                // Update the UI with the list of users and their statuses
                const usersList = document.querySelector('.users-list');
                usersList.innerHTML = message.Payload.map(user => `
                    <div data-username="${user.Username}" class="user-list-profile" onclick="window.initializeChat(event)">
                        <img src="/static/images/user.png" class="user-icon">
                        <div>
                            <p>${user.Username}</p>
                            <p class="user-status">${user.Status}</p> <!-- Display status -->
                        </div>
                    </div>
                `).join('');
            } else if (message.Type === "CHAT_MESSAGE") {
                const chatMessage = message.Payload;
                displayMessage(chatMessage.Sender, chatMessage.Content);
                console.log(chatMessage);
            } else if (message.Type === "CHAT_HISTORY") {
                console.log("Chat history received:", message.Payload);
                const messages = message.Payload;
    
                // Check if messages is a valid array before proceeding
                if (Array.isArray(messages) && messages.length > 0) {
                    prependChatMessages(messages);
                } else {
                    console.log("No chat history found or invalid message payload.");
                }
            } else if (message.Type === "NEW_MESSAGE_NOTIFICATION") {
                // New case for handling message notifications
                const sender = message.Payload.Sender;
                const content = message.Payload.Content;
                alert(`New message from ${sender}: ${content}`);
            } else {
                console.log("Unhandled message type: ", message.Type);
            }
        } catch (err) {
            console.error("Error parsing websocket message: ", err);
        }
    };
    
    window.socket.onclose = function(event) {
        console.log("Disconnected from server");
    };
    
    window.socket.onerror = function(error) {
        console.error("WebSocket error: ", error);
    };
}

function displayMessage(sender, content) {
    const messageContainer = document.getElementById('chat-messages');
    if (!messageContainer) return;

     // Create a message object
     const message = {
        Sender: sender,
        Content: content
    };
    
    // Use appendChatMessage to add the message
    appendChatMessage(message);

    messageContainer.scrollTop = messageContainer.scrollHeight;
}