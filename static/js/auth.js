import { sendRequest } from './api.js';
import { navigateTo } from './router.js';

export const isAuthenticated = async () => {
    const response = await sendRequest("/api/validate-session", "GET");
    return response.status === 200;
};

export const logout = async () => {
    const response = await sendRequest("/logout", "GET");
    if (response.ok) {
        document.cookie = 'session_token=; Max-Age=0; path=/;';
        navigateTo("/");
    } else {
        console.error("Logout failed");
    }
};

export const connectAfterLogin = (token) => {
    document.cookie = `session_token=${token}; path=/;`;
    connectWebSocket(token);
};
