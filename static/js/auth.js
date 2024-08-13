import { sendRequest, getUserInfo } from './api.js';
import { navigateTo } from './routeUtils.js';
import { initializeWebSocket } from './websocket.js';

export const isAuthenticated = async () => {
    try {
        const response = await fetch("/api/validate-session", {
            method: "GET",
            credentials: "same-origin",
        });

        return response.ok;
    } catch (error) {
        console.error("Failed to validate session:", error);
        return false;
    }
};

export const logout = async () => {
    const response = await sendRequest("/logout", "GET");
    if (response.ok) {
        localStorage.removeItem('session_token');
        localStorage.removeItem('user_id');
        localStorage.removeItem('user_name');
        navigateTo("/");
    } else {
        console.error("Logout failed");
    }
};

export const connectAfterLogin = async (token) => {
    localStorage.setItem('session_token', token);
    console.log("New token set in localStorage: ", localStorage.getItem('session_token'));

    const userInfo = await getUserInfo();
    if (userInfo) {
        localStorage.setItem('user_id', userInfo.user_id);
        localStorage.setItem('user_name', userInfo.username);
    } else {
        console.error('Failed to retrieve user info after login');
    }

    initializeWebSocket(token);
};
