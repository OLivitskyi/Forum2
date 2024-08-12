import { sendPrivateMessage } from '../websocket.js';
import { debounce } from '../routeUtils.js';

let currentReceiverID = null;
let offset = 0;
let loading = false;

export const setupMessageForm = () => {
    const sendButton = document.getElementById("send-button");
    const messageInput = document.getElementById("message-input");

    sendButton.addEventListener("click", () => {
        if (currentReceiverID && messageInput.value) {
            sendPrivateMessage(currentReceiverID, messageInput.value);
            messageInput.value = "";
        }
    });
};

export const loadMessages = async (receiverID, loadMore = false) => {
    if (loadMore) {
        offset += 10;
    } else {
        offset = 0; // reset offset if it's a fresh load
    }

    try {
        const response = await fetch(`/api/get-messages?user_id=${receiverID}&limit=10&offset=${offset}`);
        const messages = await response.json();

        const messageList = document.querySelector(".message-list");

        if (!loadMore) {
            messageList.innerHTML = ""; // Очищення попередніх повідомлень
        }

        if (Array.isArray(messages) && messages.length > 0) {
            messages.forEach(msg => {
                const messageElement = document.createElement("div");
                messageElement.classList.add(msg.sender_id === receiverID ? "other-user-message" : "user-message");
                messageElement.innerHTML = `
                    <div class="message-content">
                        <strong>${msg.sender_name}:</strong> ${msg.content}
                        <div class="message-time">${new Date(msg.timestamp).toLocaleString()}</div>
                    </div>
                `;
                if (loadMore) {
                    messageList.insertBefore(messageElement, messageList.firstChild); // Додає нові повідомлення на початок
                } else {
                    messageList.appendChild(messageElement);
                }
            });
        } else if (!loadMore) {
            console.log("No messages found for this user.");
            messageList.innerHTML = "<p>No previous messages.</p>";
        }
    } catch (error) {
        console.error("Failed to load messages:", error);
    }
};

export const handlePrivateMessage = (message) => {
    const messageList = document.querySelector(".message-list");
    if (messageList) {
        const messageElement = document.createElement("div");
        messageElement.classList.add(message.sender_id === currentReceiverID ? "other-user-message" : "user-message");
        messageElement.innerHTML = `
            <div class="message-content">
                <strong>${message.sender_name}:</strong> ${message.content}
                <div class="message-time">${new Date(message.timestamp).toLocaleString()}</div>
            </div>
        `;
        messageList.appendChild(messageElement);
        messageList.scrollTop = messageList.scrollHeight;
    }
};

export const setupMessageListScroll = () => {
    const messageList = document.querySelector(".message-list");

    messageList.addEventListener("scroll", debounce(() => {
        if (messageList.scrollTop === 0 && !loading) {
            loading = true;
            loadMessages(currentReceiverID, true)
                .finally(() => loading = false);
        }
    }, 200));
};

export const setCurrentReceiver = (receiverID) => {
    currentReceiverID = receiverID;
    offset = 0; // Reset offset when switching to a new receiver
};
