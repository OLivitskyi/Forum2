import { navigateTo } from './routeUtils.js';

export const sendRequest = async (url, method, body = null) => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    
    const sessionToken = document.cookie.split('; ').find(row => row.startsWith('session_token='));
    if (sessionToken) {
        headers.append('Authorization', `Bearer ${sessionToken.split('=')[1]}`);
    }

    const response = await fetch(url, {
        method: method,
        headers: headers,
        body: body ? JSON.stringify(body) : null
    });

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

export const createPost = async (post) => {
    return fetch('/api/create-post', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(post)
    });
};

export const createCategory = async (body) => {
    const response = await sendRequest("/api/create-category", "POST", body);
    return response;
};

export const sendMessage = async (body) => {
    const response = await sendRequest("/api/send-message", "POST", body);
    return response;
};

export const getPosts = async () => {
    const response = await sendRequest("/api/get-posts", "GET");
    if (!response.ok) {
        throw new Error("Failed to fetch posts");
    }
    return await response.json();
};

export const getPost = async (body) => {
    const response = await fetch("/api/get-post", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(body)
    });
    if (!response.ok) {
        throw new Error("Failed to fetch post");
    }
    return await response.json();
};
