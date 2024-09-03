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

    const requestMessageChat = JSON.stringify({
        Type: "CHAT_MESSAGE",
        Payload: null
    });
    socket.send(requestMessageChat)
};

socket.onmessage = function(event) {
    console.log("Message from server: ", event.data);

    // handle in console
    try {
        const message = JSON.parse(event.data);
        if (message.Type === "USER_LIST") {
            console.log("Received user list: ", message.Payload);
            // you can update the UI with the list of users, e.g., display them in a sidebar
        } else if (message.Type === "CHAT_MESSAGE") {
            console.log("let's pretend this is the Chat")
        }
        else {
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
    { username: haneen
      status: online }

    { username: fatema
      status: offline }
      
    { username: spiderman
      status: typing
      }


type2
chatMessages
Type: CHAT_MESSAGE
Payload:
    { sender: haneen
     receipent: fatema
     content:hey pookie! are you coming to reboot today? }

    { sender: fatima
     receipent: haneen
     content: hi pooks, Yes! }
    
type3
notifications

*/