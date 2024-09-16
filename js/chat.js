import { loggedInUsername } from './script.js';
let offset = 0;

window.initializeChat = function initializeChat(event) {
    const clickedElement = event.currentTarget;
    const username = clickedElement.dataset.username;
    
    if (!username) {
        console.error("Username not found on clicked element.");
        return;
    }
    
    console.log('Open chat for', username);

    // Open chat UI and display the selected user's information
    openChatWindow(username);

    offset = 0;

    // Send a request to the server to get the last 10 messages
    const chatHistoryRequest = {
        Type: 'CHAT_HISTORY',
        Payload: {
            Sender: loggedInUsername, // Current user (replace dynamically)
            Recipient: username,
            Limit: 10, // Load the last 10 messages
            Offset: offset  // No offset initially
        }
    };

     // Make sure the socket is connected before sending
    if (window.socket) {
        window.socket.send(JSON.stringify(chatHistoryRequest));
    } else {
        console.error("WebSocket is not connected.");
    }
}

function openChatWindow(username) {
    // Get the chat container
    const container = document.getElementById('main-content');
    
    if (!container) {
        console.error("Chat container is missing.");
        return;
    }

    // Example: Clear existing content and show the new chat UI
    container.innerHTML = 
        `<div id="chat-div">
            <div id="user-info-chat">
                <img src="/static/images/user.png" id="user-pic" alt="User Picture">
                <p id="user-name-chat">${username}</p>
            </div>
            <div id="chat-messages">
                <!-- Messages will be displayed here -->
            </div>
            <textarea id="chat-input" placeholder="Type your message here..."></textarea>
            <button id="send-message-button">Send</button>
        </div>`
    ;

    const chatMessages = document.getElementById('chat-messages');
    chatMessages.addEventListener('scroll', throttle(loadMoreMessages, 500));

    // Set up event listeners for the chat input and send button
    document.getElementById('send-message-button').addEventListener('click', sendMessage);
    document.getElementById('chat-input').addEventListener('keypress', function(event) {
        if (event.key === 'Enter') {
            sendMessage();
        }
    });
}

function sendMessage() {
    const chatInput = document.getElementById('chat-input');
    const message = chatInput.value.trim();
    
    if (!message) {
        alert("Message is empty. Please type a message before sending.");
        return;
    }
    
    // Get the sender username from a global state, session, or another variable
    const sender = loggedInUsername; // Replace this with how you're managing the logged-in user

    // Send the message through WebSocket
    const chatMessage = {
        Type: 'CHAT_MESSAGE',
        Payload: {
            Sender: sender, // Set the actual sender username
            Recipient: document.getElementById('user-name-chat').textContent,
            Content: message
        }
    };
    
    // Assuming you have a global WebSocket variable
    if (socket) {
        socket.send(JSON.stringify(chatMessage));
        chatInput.value = ''; // Clear input after sending
    } else {
        console.error("WebSocket is not connected.");
    }
}

function loadMoreMessages() {
    const chatMessagesDiv = document.getElementById('chat-messages');
    // Check if the user has scrolled to the top
    if (chatMessagesDiv.scrollTop === 0) {
        // Increment the offset to load the next batch of messages
        offset += 10;
        const chatHistoryRequest = {
            Type: 'CHAT_HISTORY',
            Payload: {
                Sender: loggedInUsername,
                Recipient: document.getElementById('user-name-chat').textContent,
                Limit: 10,
                Offset: offset
            }
        };
        socket.send(JSON.stringify(chatHistoryRequest));
    }
}

// Throttle function to avoid too many requests on scroll
function throttle(func, limit) {
    let inThrottle;
    return function () {
        const args = arguments;
        const context = this;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    }
}
// Function to append new chat messages
function appendChatMessage(message) {
    const chatMessagesDiv = document.getElementById('chat-messages');
    const messageElement = document.createElement('div');
    messageElement.className = 'chat-message';
    messageElement.innerHTML = `<strong>${message.Sender}</strong>: ${message.Content}`;
    chatMessagesDiv.appendChild(messageElement);
    // Scroll to the bottom when new messages are added
    chatMessagesDiv.scrollTop = chatMessagesDiv.scrollHeight;
}
// Function to prepend chat history messages
function prependChatMessages(messages) {
    const chatMessagesDiv = document.getElementById('chat-messages');
    const initialScrollHeight = chatMessagesDiv.scrollHeight;
    messages.forEach(message => {
        const messageElement = document.createElement('div');
        messageElement.className = 'chat-message';
        messageElement.innerHTML = `<strong>${message.Sender}</strong>: ${message.Content}`;
        chatMessagesDiv.insertBefore(messageElement, chatMessagesDiv.firstChild);
    });
    // Maintain scroll position after prepending messages
    chatMessagesDiv.scrollTop = chatMessagesDiv.scrollHeight - initialScrollHeight;
}