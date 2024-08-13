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
    if (!loadMore) {
        offset = 0; // Reset offset if not loading more messages
    }

    try {
        const response = await fetch(`/api/get-messages?user_id=${receiverID}&limit=10&offset=${offset}`);
        const messages = await response.json();

        const messageList = document.querySelector(".message-list");

        if (!loadMore) {
            messageList.innerHTML = "";
        }

        if (Array.isArray(messages) && messages.length > 0) {
            messages.forEach(msg => {
                const messageElement = document.createElement("div");
                messageElement.classList.add(msg.sender_id === receiverID ? "other-user-message" : "user-message");
                messageElement.innerHTML = `
                    <div class="message-content">
                        <strong>${msg.sender_name || 'Unknown User'}:</strong> ${msg.content}
                        <div class="message-time">${new Date(msg.created_at).toLocaleString()}</div>
                    </div>
                `;
                if (loadMore) {
                    messageList.insertBefore(messageElement, messageList.firstChild);
                } else {
                    messageList.appendChild(messageElement);
                }
            });

            offset += 10; // Increase offset for the next load
        } else if (!loadMore) {
            console.log("No messages found for this user.");
            messageList.innerHTML = "<p>No previous messages.</p>";
        }
    } catch (error) {
        console.error("Failed to load messages:", error);
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

const showPopupNotification = (message) => {
    const notification = document.getElementById("popup-notification");
    if (notification) {
        notification.textContent = message;
        notification.style.display = "block";

        setTimeout(() => {
            notification.style.display = "none";
        }, 3000); // Hide after 3 seconds
    }
};

export const handlePrivateMessage = (message) => {
    const messageList = document.querySelector(".message-list");
    const messagesLink = document.getElementById("messages");

    if (messageList && message.sender_id === currentReceiverID) {
        const messageElement = document.createElement("div");
        messageElement.classList.add("other-user-message");
        messageElement.innerHTML = `
            <div class="message-content">
                <strong>${message.sender_name}:</strong> ${message.content}
                <div class="message-time">${new Date(message.timestamp).toLocaleString()}</div>
            </div>
        `;
        messageList.appendChild(messageElement);
        messageList.scrollTop = messageList.scrollHeight;
    } else {
        showPopupNotification('You have a new message!');
        
        const messageCountElement = messagesLink.querySelector(".message-count");
        const currentCount = parseInt(messageCountElement.textContent, 10) || 0;
        messageCountElement.textContent = currentCount + 1;
    }
};

export const markMessagesAsRead = (receiverID) => {
    const messagesLink = document.getElementById("messages");
    const messageCountElement = messagesLink.querySelector(".message-count");
    const currentCount = parseInt(messageCountElement.textContent, 10) || 0;

    if (receiverID === currentReceiverID) {
        const unreadMessages = document.querySelectorAll(".other-user-message.unread");
        const unreadCount = unreadMessages.length;
        unreadMessages.forEach(msg => msg.classList.remove("unread"));

        messageCountElement.textContent = Math.max(currentCount - unreadCount, 0);
    }
};
