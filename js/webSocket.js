const socket = new WebSocket("ws://localhost:8080/api/ws");

socket.onopen = function(event) {
    console.log("Connected to WebSocket server");
    // socket.send("Hello Server!"); // Send a test message

    // send a message with users_list
    const requestMessage = JSON.stringify({
        Type: "USER_LIST",
        Payload: null
    });
    socket.send(requestMessage)
};

socket.onmessage = function(event) {
    console.log("Message from server: ", event.data);

    // handle in console
    try {
        const message = JSON.parse(event.data);
        if (message.Type === "USER_LIST") {
            console.log("Received user list: ", message.Payload);
            // you can update the UI with the list of users, e.g., display them in a sidebar
        } else {
            console.log("Unhandled message type: ", message.Type);
        }
    } catch (err) {
        console.error("Error parsing websocket message: ", err);
    }

};

socket.onclose = function(event) {
    console.log("Disconnected from server");
};

socket.onerror = function(error) {
    console.error("WebSocket error: ", error);
};

/*

type1
usersList
type: usersList
payload: {
    status {
             typing:
        }

    }

type2
chatMessages

type3
onlineUsers

type3
notifications

*/