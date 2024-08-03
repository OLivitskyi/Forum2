import { navigateTo } from './router.js';

// Function to send an HTTP request with session token
export const sendRequest = async (url, method, body = null) => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    // Get session token from cookies
    const sessionToken = document.cookie.split('; ').find(row => row.startsWith('session_token='));
    if (sessionToken) {
        headers.append('Authorization', `Bearer ${sessionToken.split('=')[1]}`);
    }

    const response = await fetch(url, {
        method: method,
        headers: headers,
        body: body ? JSON.stringify(body) : null
    });

    // If unauthorized, redirect to login page
    if (response.status === 401) {
        navigateTo("/");
    }

    return response;
};

// Function to check if the user is authenticated
export const isAuthenticated = async () => {
    const response = await sendRequest("/validate-session", "GET");
    return response.status === 200;
};

// Function to logout the user and redirect to login page
export const logout = async () => {
    const response = await sendRequest("/logout", "GET");
    if (response.ok) {
        navigateTo("/");
    } else {
        console.error("Logout failed");
    }
};
