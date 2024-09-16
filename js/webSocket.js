
export function initializeWebSocket() {
    window.socket = new WebSocket("ws://localhost:8080/api/ws");

    window.socket.onopen = function(event) {
        console.log("Connected to WebSocket server");
        // socket.send("Hello Server!"); // Send a test message
    
        // send a message with users_list
        const requestMessage = JSON.stringify({
            Type: "USER_LIST",
            Payload: null
        });
        window.socket.send(requestMessage)
    
        // const requestMessageChat = JSON.stringify({
        //     Type: "CHAT_MESSAGE",
        //     Payload: null
        // });
        // socket.send(requestMessageChat)
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
                console.log(chatMessage)
            } else if (message.Type === "CHAT_HISTORY") {
                const messages = message.Payload;
    
                // Check if messages is a valid array before proceeding
                if (Array.isArray(messages) && messages.length > 0) {
                    for (let i = messages.length - 1; i >= 0; i--) {
                        displayMessage(messages[i].Sender, messages[i].Content);
                    }
                } else {
                    console.log("No chat history found or invalid message payload.");
                }
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

    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.innerHTML = `<strong>${sender}:</strong> ${content}`;
    messageContainer.appendChild(messageElement);

    // Scroll to the bottom (for new messages)
    messageContainer.scrollTop = messageContainer.scrollHeight;
}

//let offset = 10; // Initial offset for chat messages

// function loadMoreMessages() {
//     const chatMessagesDiv = document.getElementById('chat-messages');

//     if (chatMessagesDiv.scrollTop === 0) { // Scrolled to top
//         // Send request to load 10 more messages
//         const username = document.getElementById('user-name-chat').textContent;

//         const chatHistoryRequest = {
//             Type: 'CHAT_HISTORY',
//             Payload: {
//                 Sender: 'fatima', // Current user (replace dynamically)
//                 Recipient: username,
//                 Limit: 10,
//                 Offset: offset // Load the next batch of 10
//             }
//         };

//         socket.send(JSON.stringify(chatHistoryRequest));
//         offset += 10; // Increase offset for the next request
//     }
// }

//  function throttle(func, delay) {
//     let lastCall = 0;
//     return function (...args) {
//         const now = new Date().getTime();
//         if (now - lastCall < delay) {
//             return;
//         }
//         lastCall = now;
//         return func(...args);
//     };
// }