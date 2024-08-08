import { sendMessage } from '../websocket.js';
import { showError, clearError } from '../errorHandler.js';
import { navigateTo } from '../routeUtils.js'; // Оновлено шлях до routeUtils.js

export const setupMessageForm = () => {
    const messageForm = document.getElementById("message-form");
    if (messageForm) {
        messageForm.removeEventListener("submit", handleSendMessage);
        messageForm.addEventListener("submit", handleSendMessage, { once: true });
    }
};

export async function handleSendMessage(e) {
    e.prevent.preventDefault();
    const receiverID = document.getElementById("receiver-id").value;
    const content = document.getElementById("message-content").value;
    if (!receiverID || !content) {
        showError("Receiver ID and content are required");
        return;
    }
    try {
        await sendMessage({ type: "message", receiver_id: receiverID, content });
        clearError();
        navigateTo("/messages");
    } catch (error) {
        showError("Failed to send message");
    }
};
