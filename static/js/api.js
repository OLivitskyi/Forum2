import { navigateTo } from './router.js';

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

export const getCategories = async () => {
    const response = await sendRequest("/api/get-categories", "GET");

    if (!response.ok) {
        throw new Error("Failed to fetch categories");
    }
    return await response.json();
};

export const createPost = async (body) => {
    const response = await sendRequest("/api/create-post", "POST", body);
    return response;
};

export const createCategory = async (body) => {
    const response = await sendRequest("/api/create-category", "POST", body);
    return response;
};

export const sendMessage = async (body) => {
    const response = await sendRequest("/api/send-message", "POST", body);
    return response;
};
