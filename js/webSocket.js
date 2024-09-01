const socket = new WebSocket("ws://localhost:8080/api/ws");

socket.onopen = function(event) {
    console.log("Connected to WebSocket server");
    socket.send("Hello Server!"); // Send a test message
};

socket.onmessage = function(event) {
    console.log("yay! ig");
    console.log("Message from server: ", event.data);
};

socket.onclose = function(event) {
    console.log("Disconnected from server");
};

socket.onerror = function(error) {
    console.error("WebSocket error: ", error);
};
