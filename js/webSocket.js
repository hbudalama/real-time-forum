import { appendChatMessage,  prependChatMessages } from './chat.js';
import { loggedInUsername } from './script.js';
let typingTimeout;
const typingDelay = 2000; // 2 seconds delay to indicate "stopped typing"

export function initializeWebSocket() {
    const socket = new WebSocket("ws://localhost:8080/api/ws");
    window.socket = socket;

    socket.onopen = () => {
        console.log("Connected to WebSocket server");
        requestUserList();
    };
    
    window.socket.onmessage = (event) => {
        handleMessage(event.data);
    };
    
    window.socket.onclose = function(event) {
        console.log("Disconnected from server");
    };
    
    window.socket.onerror = function(error) {
        console.error("WebSocket error: ", error);
    };
}

function requestUserList() {
    const requestMessage = JSON.stringify({
        Type: "USER_LIST",
        Payload: null
    });
    window.socket.send(requestMessage);
}

function handleMessage(data) {
    try {
        const message = JSON.parse(data);
        switch (message.Type) {
            case "USER_LIST":
                updateUsersList(message.Payload);
                break;
            case "CHAT_MESSAGE":
                displayMessage(message.Payload.Sender, message.Payload.Content, message.Payload.CreatedDate);
                break;
            case "CHAT_HISTORY":
                handleChatHistory(message.Payload);
                break;
            case "NEW_MESSAGE_NOTIFICATION":
                handleNewMessageNotification(message.Payload);
                break;
            case "TYPING_STATUS":
                handleTypingStatus(message.Payload);
                break;
            default:
                console.log("Unhandled message type: ", message.Type);
        }
    } catch (err) {
        console.error("Error parsing websocket message: ", err);
    }
}

function updateUsersList(users) {
    const usersList = document.querySelector('.users-list');
    if (!usersList || !Array.isArray(users)) return;

    usersList.innerHTML = users.map(user => `
        <div data-username="${user.Username}" class="user-list-profile" onclick="window.initializeChat(event)">
            <img src="/static/images/user.png" class="user-icon">
            <div>
                <p>${user.Username}</p>
                <p class="user-status">${user.Status}</p> <!-- Display status -->
            </div>
        </div>
    `).join('');
}

function displayMessage(sender, content, createdDate) {
    const messageContainer = document.getElementById('chat-messages');
    if (!messageContainer) return;

     // Create a message object
     const message = {
        Sender: sender,
        Content: content,
        CreatedDate: createdDate
    };
    
    // Use appendChatMessage to add the message
    appendChatMessage(message);

    messageContainer.scrollTop = messageContainer.scrollHeight;
}

function handleChatHistory(messages) {
    console.log("Chat history received:", messages);
    if (Array.isArray(messages) && messages.length > 0) {
        prependChatMessages(messages);
    } else {
        console.log("No chat history found or invalid message payload.");
    }
}

function handleNewMessageNotification({ Sender, Content }) {
    alert(`New message from ${Sender}: ${Content}`);
}

// Handle typing status received from WebSocket
function handleTypingStatus(payload) {
    const typingStatusDiv = document.getElementById('typing-status');
    const currentChatUser = document.getElementById('user-name-chat').textContent;

    // Ensure typing status is only shown if the typing sender is the person you're chatting with
    if (payload.Sender === currentChatUser && payload.Recipient === loggedInUsername) {
        console.log("here")
        if (payload.IsTyping) {
            typingStatusDiv.style.display = 'block'; // Show typing indicator
        } else {
            typingStatusDiv.style.display = 'none'; // Hide typing indicator
        }
            // Clear the typing status after the specified delay
    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
        if (typingStatusDiv) {
            typingStatusDiv.style.display = 'none'; // Hide typing status
        }
        // sendTypingToRecipient(false);
    }, typingDelay);
    }
}
