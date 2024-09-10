function initializeChat(event) {
    const clickedElement = event.currentTarget;
    const username = clickedElement.dataset.username;
    
    if (!username) {
        console.error("Username not found on clicked element.");
        return;
    }
    
    console.log('Open chat for', username);

    // Open chat UI and display the selected user's information
    openChatWindow(username);

    // Send a request to the server to get the last 10 messages
    const chatHistoryRequest = {
        Type: 'CHAT_HISTORY',
        Payload: {
            Sender: 'fatima', // Current user (replace dynamically)
            Recipient: username,
            Limit: 10, // Load the last 10 messages
            Offset: 0  // No offset initially
        }
    };
     socket.send(JSON.stringify(chatHistoryRequest));
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

    document.getElementById('chat-messages').addEventListener('scroll', throttle(loadMoreMessages, 500));

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
        alert("Message can't empty.");
        return;
    }
    
    // Get the sender username from a global state, session, or another variable
    const sender = 'fatima'; // Replace this with how you're managing the logged-in user

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