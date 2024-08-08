import { sendRequest } from './api.js';
import { navigateTo } from './routeUtils.js';
import { connectWebSocket, setupWebSocketHandlers } from './websocket.js';

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
        localStorage.removeItem('session_token');
        navigateTo("/");
    } else {
        console.error("Logout failed");
    }
};

export const connectAfterLogin = (token) => {
    localStorage.setItem('session_token', token);
    console.log("New token set in localStorage: ", localStorage.getItem('session_token'));
    connectWebSocket(token);
    setupWebSocketHandlers();
};