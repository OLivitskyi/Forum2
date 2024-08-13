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
    try {
        const limit = 10;
        offset = loadMore ? offset + limit : 0;
        
        const response = await fetch(`/api/get-messages?user_id=${receiverID}&limit=${limit}&offset=${offset}`);
        const messages = await response.json();

        const messageList = document.querySelector(".message-list");

        if (!loadMore) {
            messageList.innerHTML = "";
        }

        if (Array.isArray(messages) && messages.length > 0) {
            // Сортуємо повідомлення від найновішого до найстарішого
            messages.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

            messages.forEach(msg => {
                const messageElement = document.createElement("div");
                messageElement.classList.add(msg.sender_id === receiverID ? "other-user-message" : "user-message");
                messageElement.innerHTML = `
                    <div class="message-content">
                        <strong>${msg.sender_name || 'Unknown User'}:</strong> ${msg.content}
                        <div class="message-time">${new Date(msg.created_at).toLocaleString()}</div>
                    </div>
                `;
                
                // Додаємо повідомлення на початок, якщо завантажуємо більше, інакше в кінець
                if (loadMore) {
                    messageList.prepend(messageElement);
                } else {
                    messageList.appendChild(messageElement);
                }
            });

            // Прокрутити до нових повідомлень тільки при першому завантаженні
            if (!loadMore) {
                messageList.scrollTop = messageList.scrollHeight;
            }
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
        // Перевіряємо, чи користувач прокрутив до верхньої частини списку повідомлень
        if (messageList.scrollTop === 0 && !loading) {
            const oldScrollHeight = messageList.scrollHeight;

            loading = true;
            loadMessages(currentReceiverID, true).finally(() => {
                loading = false;

                // Відновлення позиції скролу після завантаження нових повідомлень
                const newScrollHeight = messageList.scrollHeight;
                messageList.scrollTop = newScrollHeight - oldScrollHeight;
            });
        }
    }, 200));
};

export const setCurrentReceiver = (receiverID) => {
    currentReceiverID = receiverID;
    offset = 0;
};

const showPopupNotification = (message) => {
    const notification = document.getElementById("popup-notification");
    if (notification) {
        notification.textContent = message;
        notification.style.display = "block";

        setTimeout(() => {
            notification.style.display = "none";
        }, 3000);
    }
};

export const handlePrivateMessage = (message) => {
    const messageList = document.querySelector(".message-list");
    const messagesLink = document.getElementById("messages");

    if (messageList) {
        const messageElement = document.createElement("div");
        messageElement.classList.add(message.sender_id === currentReceiverID ? "other-user-message" : "user-message");
        messageElement.innerHTML = `
            <div class="message-content">
                <strong>${message.sender_name}:</strong> ${message.content}
                <div class="message-time">${new Date(message.timestamp).toLocaleString()}</div>
            </div>
        `;
        messageList.insertBefore(messageElement, messageList.firstChild); // Додаємо нове повідомлення зверху
        messageList.scrollTop = 0; // Прокручуємо до верху, щоб показати нове повідомлення
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
