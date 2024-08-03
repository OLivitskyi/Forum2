let socket;
let messageHandler = () => {};
let reactionHandler = () => {};

export const connectWebSocket = (sessionToken) => {
    socket = new WebSocket(`ws://localhost:8080/ws?session_token=${sessionToken}`);

    socket.onopen = () => {
        console.log("Connected to WebSocket server");
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        if (message.type === "message") {
            messageHandler(message);
        } else if (message.type === "reaction") {
            reactionHandler(message);
        }
    };

    socket.onclose = () => {
        console.log("Disconnected from WebSocket server");
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
};

export const sendMessage = (message) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.error("WebSocket is not connected");
    }
};

export const setMessageHandler = (handler) => {
    messageHandler = handler;
};

export const setReactionHandler = (handler) => {
    reactionHandler = handler;
};
