import { navigateToPostDetails } from "./routeUtils.js";
import { handlePrivateMessage } from "./handlers/messageHandlers.js";
import { handleUserStatus, requestUserStatus } from "./handlers/userStatusHandlers.js";
import { getUserInfo } from "./api.js"; 
import { renderUserList } from "./components/userList.js"; 

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

    socket = new WebSocket(
        `ws://localhost:8080/ws?session_token=${sessionToken}`
    );

    socket.onopen = () => {
        console.log("Connected to WebSocket server");
        isConnected = true;
        reconnectAttempts = 0;

        while (messageQueue.length > 0) {
            const message = messageQueue.shift();
            sendMessage(message);
        }

        requestUserStatus();
    };

    socket.onmessage = async (event) => {
        const message = JSON.parse(event.data);
        console.log("Received message:", message);

        if (message.type === "error" && message.data === "Unauthorized") {
            console.error("Unauthorized access, please log in again.");
            socket.close();
            return;
        }

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
                handlePrivateMessage(message.data);
                break;
            default:
                console.warn("Unknown message type:", message.type);
        }
        const currentUserInfo = await getUserInfo();
        const currentUserId = currentUserInfo.user_id;
        renderUserList("user-status-list", JSON.parse(localStorage.getItem("users")) || [], currentUserId, (userId) => {
            console.log(`User ${userId} clicked in global user list.`);
        });
    };

    socket.onclose = (event) => {
        isConnected = false;
        console.log(
            "Disconnected from WebSocket server, code:",
            event.code,
            "reason:",
            event.reason
        );

        if (event.code === 1000) {
            console.log("Connection closed normally.");
            return;
        }

        if (reconnectAttempts < maxReconnectAttempts) {
            reconnectAttempts++;
            console.log(
                `Attempting to reconnect... (${reconnectAttempts}/${maxReconnectAttempts})`
            );
            setTimeout(connectWebSocket, 5000);
        } else {
            console.error(
                "Max reconnect attempts reached. Could not reconnect to WebSocket server."
            );
        }
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
};
export const sendPost = (post) => {
    sendMessage({ type: "post", data: post });
};

export const sendComment = (comment) => {
    sendMessage({ type: "comment", data: comment });
};
export const sendMessage = (message) => {
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

        const categories = post.categories
            .map((category) => `<span class="category">${category.name}</span>`)
            .join(", ");
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

        postElement.addEventListener("click", () => {
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

export const initializeWebSocket = (token = null) => {
    sessionToken = token || localStorage.getItem("session_token");

    if (sessionToken) {
        connectWebSocket();
    } else {
        console.error(
            "No session token available to initialize WebSocket connection"
        );
    }
};
