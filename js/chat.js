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
    // Send a request to the server to get the last 10 messages on chat open
    requestChatHistory(username, offset);
}

function requestChatHistory(recipient, offset) {
    const chatHistoryRequest = {
        Type: 'CHAT_HISTORY',
        Payload: {
            Sender: loggedInUsername, // Current user
            Recipient: recipient,
            Limit: 10,
            Offset: offset
        }
    };

    // Send request for chat history
    if (window.socket && window.socket.readyState === WebSocket.OPEN) {
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
                <p id="typing-status" style="display:none;">is typing...</p>
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
    document.getElementById('chat-input').addEventListener('keypress', handleTyping);
     
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
    console.log("Load more messages function triggered");

    const chatMessagesDiv = document.getElementById('chat-messages');

    console.log("Scroll Height:", chatMessagesDiv.scrollHeight);
    console.log("Client Height:", chatMessagesDiv.clientHeight);
    console.log("Scroll Top:", chatMessagesDiv.scrollTop);
    
    // Check if the user has scrolled to the top
    if (chatMessagesDiv.scrollTop <= 10) {
        console.log("Near top, loading more messages...")
        // Increment the offset to load the next batch of messages
        offset += 10;
        const recipient = document.getElementById('user-name-chat').textContent;
        requestChatHistory(recipient, offset);  // Load more messages using offset
    }
}

// Throttle function to avoid too many requests on scroll
function throttle(func, limit) {
    let inThrottle;
    return function () {
        console.log("Throttle function triggered");
        const args = arguments;
        const context = this;
        if (!inThrottle) {
            console.log("Throttled function called");
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => {
                console.log("Throttling ended");
                inThrottle = false;
            }, limit);
        }
    };
}

// Function to append new chat messages
export function appendChatMessage(message) {
    const chatMessagesDiv = document.getElementById('chat-messages');

    const messageElement = document.createElement('div');
    messageElement.className = 'chat-message';

    messageElement.innerHTML = `<strong>${message.Sender} ${message.CreatedDate}</strong>: ${message.Content}`;
    chatMessagesDiv.appendChild(messageElement);

}


// Function to prepend chat history messages
export function prependChatMessages(messages) {
    const chatMessagesDiv = document.getElementById('chat-messages');
    const initialScrollHeight = chatMessagesDiv.scrollHeight;
    const initialScrollTop = chatMessagesDiv.scrollTop;
    messages.forEach(message => {
        const messageElement = document.createElement('div');
        messageElement.className = 'chat-message';
        messageElement.innerHTML = `<strong>${message.Sender} ${message.CreatedDate}</strong>: ${message.Content}`;
        chatMessagesDiv.insertBefore(messageElement, chatMessagesDiv.firstChild);
    });
    // Adjust scroll position to maintain it relative to the loaded messages
    chatMessagesDiv.scrollTop = initialScrollTop + (chatMessagesDiv.scrollHeight - initialScrollHeight);
}

function handleTyping(event) {
    if (event.key === 'Enter') {
        sendMessage();
    }
    // Send typing status to the server
    sendTypingToRecipient(true);
}

export function sendTypingToRecipient(isTyping) {
    const recipient = document.getElementById('user-name-chat').textContent;

     // Check if the recipient is correctly retrieved
     if (!recipient) {
        console.error("Recipient is not set! Check the user-name-chat element.");
        return; // Exit if the recipient is not found
    }
    const message = {
        Type: "TYPING_STATUS",
        Payload: {
            Sender: loggedInUsername,
            Recipient: recipient,
            IsTyping: isTyping
        }
    };
    window.socket.send(JSON.stringify(message));
    console.log("Recipient username:", message.Payload.Recipient);
    console.log("sender:", message.Payload.Sender ? message.Payload.Sender : "not there")
}