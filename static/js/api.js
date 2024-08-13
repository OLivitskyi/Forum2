import { navigateTo } from './routeUtils.js';

export const sendRequest = async (url, method, body = null) => {
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');

    const sessionToken = localStorage.getItem('session_token');
    if (sessionToken) {
        headers.append('Authorization', `Bearer ${sessionToken}`);
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
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('session_token')}`
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
    const response = await sendRequest("/api/posts", "GET");
    if (!response || !response.ok) {
        const errorText = await response.text();
        throw new Error("Failed to fetch posts: " + errorText);
    }
    return await response.json();
};

export const getPost = async (postId) => {
    const response = await sendRequest(`/api/get-post/${postId}`, "GET");
    if (!response || !response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to fetch post: ${errorText}`);
    }
    return await response.json();
};

export const createComment = async (comment) => {
    const response = await sendRequest(`/api/post/comments/new`, "POST", comment);

    if (response.status === 201) {
        console.log("Comment created successfully.");
        return {};
    }

    if (!response || !response.ok) {
        const errorText = await response.text();
        throw new Error("Failed to create comment: " + errorText);
    }

    const responseText = await response.text();
    if (responseText.trim() === "") {
        throw new Error("Empty response from the server");
    }

    try {
        return JSON.parse(responseText);
    } catch (error) {
        throw new Error("Failed to parse JSON response: " + error.message);
    }
};

export const getComments = async (postId) => {
    try {
        const response = await sendRequest(`/api/post-comments/${postId}`, "GET");
        if (!response || !response.ok) {
            const errorText = await response.text();
            throw new Error("Failed to fetch comments: " + errorText);
        }

        const comments = await response.json();

        return Array.isArray(comments) ? comments : [];
    } catch (error) {
        console.error("Failed to get comments:", error);
        return [];
    }
};

function validateUUID(uuid) {
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    return uuidRegex.test(uuid);
}

export const getUserInfo = async () => {
    try {
        const response = await fetch('/api/get-user-info', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('session_token')}`,
            }
        });

        if (!response.ok) {
            throw new Error('Failed to fetch user info');
        }

        const userInfo = await response.json();
        return userInfo;
    } catch (error) {
        console.error('Error fetching user info:', error);
        return null;
    }
};
