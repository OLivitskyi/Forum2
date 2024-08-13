import { navigateToPostDetails } from './routeUtils.js';
import { handlePrivateMessage } from './handlers/messageHandlers.js';

let socket;
let isConnected = false;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
const messageQueue = [];
let sessionToken = null;

const connectWebSocket = () => {
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
        isConnected = true;
        reconnectAttempts = 0;

        while (messageQueue.length > 0) {
            const message = messageQueue.shift();
            sendMessage(message);
        }

        // Send login message
        sendMessage({ type: 'login' });
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("Received message:", message);

        switch (message.type) {
            case "post":
                handlePost(message.data);
                break;
            case "comment":
                handleComment(message.data);
                break;
            case "user_status":
                handleUserStatus(message.data);
                break;
            case "private_message":
                handlePrivateMessage(message.data); // Нова обробка приватних повідомлень
                break;
            default:
                console.warn("Unknown message type:", message.type);
        }
    };

    socket.onclose = (event) => {
        isConnected = false;
        console.log("Disconnected from WebSocket server, code:", event.code, "reason:", event.reason);

        if (reconnectAttempts < maxReconnectAttempts) {
            reconnectAttempts++;
            console.log(`Attempting to reconnect... (${reconnectAttempts}/${maxReconnectAttempts})`);
            setTimeout(connectWebSocket, 5000);
        } else {
            console.error("Max reconnect attempts reached. Could not reconnect to WebSocket server.");
        }

        // Send logout message
        sendMessage({ type: 'logout' });
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
};

const sendMessage = (message) => {
    if (isConnected && socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.warn("WebSocket is not connected. Queuing message...");
        messageQueue.push(message);
    }
};

const handlePost = (post) => {
    const postsContainer = document.getElementById("posts-container");
    if (postsContainer) {
        if (document.getElementById(`post-${post.id}`)) {
            console.warn(`Post with ID ${post.id} already exists, skipping`);
            return;
        }

        const categories = post.categories.map(category => `<span class="category">${category.name}</span>`).join(', ');
        const postElement = document.createElement("div");
        postElement.classList.add("post");
        postElement.id = `post-${post.id}`;
        postElement.innerHTML = `
            <h3>${post.subject}</h3>
            <p>${post.content}</p>
            <div class="post-categories">Categories: ${categories}</div>
            <div>
                <span>Likes: ${post.like_count}</span>
                <span>Dislikes: ${post.dislike_count}</span>
            </div>
        `;

        postElement.addEventListener('click', () => {
            navigateToPostDetails(post.id);
        });

        postsContainer.prepend(postElement);
    }
};

const handleComment = (comment) => {
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
};

const handleUserStatus = (users) => {
    const userContainer = document.getElementById("box1");
    if (!userContainer) return;

    const currentUserId = localStorage.getItem('user_id'); 

    userContainer.innerHTML = ''; 

    users
        .filter(user => user.user_id !== currentUserId)
        .forEach(user => {
            const userElement = document.createElement("div");
            userElement.classList.add("user-box");
            userElement.dataset.userId = user.user_id;
            const statusClass = user.is_online ? "logged-in" : "logged-out";
            userElement.innerHTML = `<span class="${statusClass}">●</span>${user.username}`;
            userContainer.appendChild(userElement);
        });
};


// Додаємо функцію для запиту статусів користувачів
export const requestUserStatus = () => {
    sendMessage({ type: 'request_user_status' });
};

export const sendPost = (post) => {
    sendMessage({ type: 'post', data: post });
};

export const sendComment = (comment) => {
    sendMessage({ type: 'comment', data: comment });
};

export const sendPrivateMessage = async (receiverID, content) => {
    const username = localStorage.getItem('user_name');
    
    if (!username) {
        console.error('User name not found in localStorage');
        return;
    }

    const message = {
        type: 'private_message',
        data: { 
            receiver_id: receiverID, 
            content: content, 
            sender_name: username, // Використовуємо ім'я користувача з localStorage
            timestamp: new Date().toISOString() 
        }
    };

    // Відправляємо повідомлення на сервер
    sendMessage(message);

    // Додаємо повідомлення до списку локально, щоб відправник побачив його одразу
    handlePrivateMessage(message.data);
};

export const initializeWebSocket = (token) => {
    sessionToken = token;
    connectWebSocket();
};
