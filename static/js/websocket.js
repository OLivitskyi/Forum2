let socket;
let messageHandler = () => {};
let reactionHandler = () => {};
let postHandler = () => {};

export const connectWebSocket = (sessionToken) => {
    if (!sessionToken) {
        console.error("No session token provided for WebSocket connection");
        return;
    }

    socket = new WebSocket(`ws://localhost:8080/ws?session_token=${sessionToken}`);

    socket.onopen = () => {
        console.log("Connected to WebSocket server");
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("Received message:", message);
        if (message.type === "message") {
            messageHandler(message);
        } else if (message.type === "reaction") {
            reactionHandler(message);
        } else if (message.subject) { // New post
            postHandler(message);
        }
    };

    socket.onclose = () => {
        console.log("Disconnected from WebSocket server");
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
};

export const sendPost = (post) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
        console.log("Sending post:", post);
        socket.send(JSON.stringify({ type: 'post', data: post }));
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

export const setPostHandler = (handler) => {
    postHandler = handler;
};

export const setupWebSocketHandlers = () => {
    setPostHandler((post) => {
        const postsContainer = document.getElementById("posts-container");
        if (postsContainer) {
            const categories = post.categories.map(category => `<span class="category">${category.name}</span>`).join(', ');
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.innerHTML = `
                <h3>${post.subject}</h3>
                <p>${post.content}</p>
                <div class="post-categories">Categories: ${categories}</div>
                <div>
                    <span>Likes: ${post.like_count}</span>
                    <span>Dislikes: ${post.dislike_count}</span>
                </div>
            `;
            postsContainer.prepend(postElement);
        }
    });
};
