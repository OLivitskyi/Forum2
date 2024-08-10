let socket;
let messageHandler = () => {};
let reactionHandler = () => {};
let postHandler = () => {};
let commentHandler = () => {};

export const connectWebSocket = (sessionToken) => {
    if (!sessionToken) {
        console.error("No session token provided for WebSocket connection");
        return;
    }

    if (socket && socket.readyState === WebSocket.OPEN) {
        console.log("WebSocket is already connected");
        return;
    }

    socket = new WebSocket(`ws://localhost:8080/ws?session_token=${sessionToken}`);

    socket.onopen = () => {
        console.log("Connected to WebSocket server");
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("Received message:", message);

        switch (message.type) {
            case "message":
                messageHandler(message);
                break;
            case "reaction":
                reactionHandler(message);
                break;
            case "post":
                postHandler(message.data); // Передаємо тільки дані посту
                break;
            case "comment":
                commentHandler(message.data); // Передаємо тільки дані коментаря
                break;
            default:
                console.warn("Unknown message type:", message);
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
        const message = { type: 'post', data: post };
        console.log("Sending post:", message);
        socket.send(JSON.stringify(message));
    } else {
        console.error("WebSocket is not connected");
    }
};

export const sendComment = (comment) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
        const message = { type: 'comment', data: comment };
        console.log("Sending comment:", message);
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

export const setPostHandler = (handler) => {
    postHandler = handler;
};

export const setCommentHandler = (handler) => {
    commentHandler = handler;
};

export const setupWebSocketHandlers = () => {
    setPostHandler((post) => {
        console.log("Handling new post:", post);
        const postsContainer = document.getElementById("posts-container");
        if (postsContainer) {
            // Перевіряємо, чи існує вже елемент з таким ID
            if (document.getElementById(`post-${post.id}`)) {
                console.warn(`Post with ID ${post.id} already exists, skipping`);
                return;
            }

            const categories = post.categories.map(category => `<span class="category">${category.name}</span>`).join(', ');
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.id = `post-${post.id}`; // Додаємо ID до посту
            postElement.innerHTML = `
                <h3>${post.subject}</h3>
                <p>${post.content}</p>
                <div class="post-categories">Categories: ${categories}</div>
                <div>
                    <span>Likes: ${post.like_count}</span>
                    <span>Dislikes: ${post.dislike_count}</span>
                </div>
            `;
            postElement.onclick = () => navigateToPostDetails(post.id); // Прив'язуємо подію кліку

            postsContainer.prepend(postElement);
        }
    });

    setCommentHandler((comment) => {
        console.log("Handling new comment:", comment);
        const commentsContainer = document.getElementById("comments-container");
        if (commentsContainer) {
            const commentElement = document.createElement("div");
            commentElement.classList.add("comment");
            commentElement.innerHTML = `
                <h4>${comment.user.username}</h4>
                <p>${comment.content}</p>
                <div>
                    <span>Likes: ${comment.like_count}</span>
                    <span>Dislikes: ${comment.dislike_count}</span>
                </div>
            `;
            commentsContainer.appendChild(commentElement);
        }
    });
};
