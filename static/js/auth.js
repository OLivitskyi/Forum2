import { sendRequest } from './api.js';
import { navigateTo } from './router.js';
import { connectWebSocket, setupWebSocketHandlers } from './websocket.js'; // Додано імпорт setupWebSocketHandlers

export const isAuthenticated = async () => {
    try {
        const response = await fetch("/api/validate-session", {
            method: "GET",
            credentials: "same-origin",
        });

        if (response.ok) {
            return true;
        } else {
            return false;
        }
    } catch (error) {
        console.error("Failed to validate session:", error);
        return false;
    }
};

export const logout = async () => {
    const response = await sendRequest("/logout", "GET");
    if (response.ok) {
        clearCookies();
        navigateTo("/");
    } else {
        console.error("Logout failed");
    }
};

export const connectAfterLogin = (token) => {
    // Очищення старих кукі
    clearCookies();
    // Затримка для забезпечення, що кукі очищені
    setTimeout(() => {
        // Встановлення нового кукі
        document.cookie = `session_token=${token}; path=/;`;
        console.log("New cookie set: ", document.cookie);
        // Ініціалізація обробників WebSocket
        connectWebSocket(token);
        setupWebSocketHandlers();
    }, 100);
};




// Функція для видалення всіх кукі
function clearCookies() {
    const cookies = document.cookie.split(";");
    for (const cookie of cookies) {
        const eqPos = cookie.indexOf("=");
        const name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/;";
        console.log(`Cleared cookie: ${name}`);
    }
}
