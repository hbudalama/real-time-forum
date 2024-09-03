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
}

function openChatWindow(username) {
    // Get the chat container
    const container = document.getElementById('main-content');
    
    if (!container) {
        console.error("Chat container is missing.");
        return;
    }

    // Example: Clear existing content and show the new chat UI
    container.innerHTML = `
        <div id="chat-div">
            <div id="user-info-chat">
                <img src="/static/images/user.png" id="user-pic" alt="User Picture">
                <p id="user-name-chat">${username}</p>
            </div>
            <div id="chat-messages">
                <!-- Messages will be displayed here -->
            </div>
            <textarea id="chat-input" placeholder="Type your message here..."></textarea>
            <button id="send-message-button">Send</button>
        </div>
    `;

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
        console.error("Message is empty.");
        return;
    }
    
    // Send the message through WebSocket
    const chatMessage = {
        Type: 'chatMessage',
        Payload: JSON.stringify({
            Sender: 'currentUsername', // Replace with actual current username
            Recipient: document.getElementById('user-name-chat').textContent,
            Content: message
        })
    };
    
    // Assuming you have a global WebSocket variable
    if (WebSocket) {
        WebSocket.Send(JSON.stringify(chatMessage));
        chatInput.value = ''; // Clear input after sending
    } else {
        console.error("WebSocket is not connected.");
    }
}