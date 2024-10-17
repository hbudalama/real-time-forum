import { loggedInUsername } from './script.js';
let offset = 0;
let typingTimeout; // Added variable to manage typing timeout
const typingDelay = 1000; // Time in milliseconds to wait before sending typing status
export let activeChatRecipient = null;

// @ts-expect-error fuck u
window.initializeChat = function initializeChat(event) {
    const clickedElement = event.currentTarget;
    const username = clickedElement.dataset.username;


    if (!username) {
        console.error("Username not found on clicked element.");
        return;
    }

    activeChatRecipient = username;
    chatOpened(username);
    
    console.log('Open chat for', username);



    // Open chat UI and display the selected user's information
    openChatWindow(username);

    offset = 0;
    // Send a request to the server to get the last 10 messages on chat open
    requestChatHistory(username, offset);
}

/**
 * 
 * @param {string} recipient 
 */
export function chatOpened(recipient) {
    const chatOpened = {
        Type: 'CHAT_OPENED',
        Payload: {
            Recipient: recipient
        }
    };
    console.log("CHAT OPENED")
    window.socket.send(JSON.stringify(chatOpened));
}

window.addEventListener('beforeunload', function () {
    if (activeChatRecipient) {
        console.log("CHAT CLOSING")
        chatClosed();
    }
});

// window.addEventListener('visibilitychange', function () {
//     if (document.visibilityState === 'hidden' && activeChatRecipient) {
//         console.log("CHAT CLOSING due to tab switch or visibility change");
//         chatClosed();
//     }
// });

export function chatClosed() {
    activeChatRecipient = null;
    const chatClosed = {
        Type: 'CHAT_CLOSED',
        Payload: {
        }
    };
    if (window.socket && window.socket.readyState === WebSocket.OPEN) {
        window.socket.send(JSON.stringify(chatClosed));
        console.log("CHAT CLOSED sent to WebSocket");
    } else {
        console.error("WebSocket is not connected.");
    }
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
    const recipientUsername = document.getElementById('user-name-chat').textContent;
    document.getElementById('chat-input').addEventListener('keydown', (e) => {
        handleTyping(e, recipientUsername)
    });
}

function sendMessage() {
    const chatInput = document.getElementById('chat-input');
    const message = chatInput.value.trim();
    
    if (!message) {
        alert("Message is empty. Please type a message before sending.....");
        return;
    }

    // Send the message through WebSocket
    const chatMessage = {
        Type: 'CHAT_MESSAGE',
        Payload: {
            Recipient: activeChatRecipient,
            Content: message,
        }
    };

    if (window.socket) {
        window.socket.send(JSON.stringify(chatMessage));
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
    const chatMessagesDiv = document.getElementById('chat-messages'); // Ensure chat-messages div is defined
    if (!chatMessagesDiv) {
        console.error("Chat messages div is missing.");
        return;
    }

        const messageElement = document.createElement('div');
        messageElement.className = 'chat-message';
    
        // Format the timestamp before displaying it
        const timestamp = new Date(message.CreatedDate).toLocaleString(); 
        messageElement.innerHTML = `<strong>${message.Sender} ${timestamp}</strong>: ${message.Content}`;
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
        
        // Format the timestamp before displaying it
        const timestamp = new Date(message.CreatedDate).toLocaleString();
        messageElement.innerHTML = `<strong>${message.Sender} ${timestamp}</strong>: ${message.Content}`;
        chatMessagesDiv.insertBefore(messageElement, chatMessagesDiv.firstChild);
    });
    // Adjust scroll position to maintain it relative to the loaded messages
    chatMessagesDiv.scrollTop = initialScrollTop + (chatMessagesDiv.scrollHeight - initialScrollHeight);
}

function handleTyping(event, recipientUsername) {
    if (event.key === 'Enter') {
        event.preventDefault(); 
        sendMessage();
    }
    // Clear timeout for typing status if typing stops
    clearTimeout(typingTimeout);
    
    // Send typing status to the server
    sendTypingToRecipient(true, recipientUsername);

    // Set a timeout to indicate typing has stopped after a delay
    typingTimeout = setTimeout(() => {
        sendTypingToRecipient(false, recipientUsername);
    }, typingDelay);
}

// Function to send typing status to the recipient
export function sendTypingToRecipient(isTyping, recipientUsername) {
    // Check if the recipient is valid
    if (!recipientUsername) {
        console.error("Recipient username not found.");
        return;
    }

    const message = {
        Type: "TYPING_STATUS",
        Payload: {
            Sender: loggedInUsername,
            Recipient: recipientUsername,
            IsTyping: isTyping
        }
    };
    if (window.socket) {
        window.socket.send(JSON.stringify(message));
    }
}

