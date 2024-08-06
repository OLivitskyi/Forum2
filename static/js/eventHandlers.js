import { navigateTo } from './router.js';
import { createPost, createCategory, getCategories, getPosts, sendMessage } from './api.js';
import { logout } from './auth.js';
import { setFormMessage } from './formHandler.js';
import { showError, clearError } from './errorHandler.js';

export const handleLoginFormSubmit = () => {
    const loginForm = document.getElementById("login");
    if (loginForm) {
        loginForm.removeEventListener("submit", handleLogin);
        loginForm.addEventListener("submit", handleLogin, { once: true });
    }

    async function handleLogin(e) {
        e.preventDefault();
        clearError();

        const formData = new FormData(e.target);
        const response = await fetch("/api/login", {
            method: "POST",
            body: formData,
        });

        if (response.ok) {
            const token = await response.text();
            document.cookie = `session_token=${token}; path=/;`;

            clearError();
            navigateTo("/homepage");
        } else {
            const errorText = await response.text();
            showError(errorText || "Login failed. Please try again.");
        }
    }
};

export const handleLogout = () => {
    const logoutButton = document.getElementById("logout");
    if (logoutButton) {
        logoutButton.removeEventListener("click", handleLogoutClick);
        logoutButton.addEventListener("click", handleLogoutClick, { once: true });
    }

    async function handleLogoutClick(e) {
        e.preventDefault();
        await logout();
    }
};

export const handleCreatePostFormSubmit = (clearError, showError) => {
    const form = document.getElementById("create-post-form");
    if (form) {
        form.removeEventListener("submit", handleSubmit);
        form.addEventListener("submit", handleSubmit, { once: true });
    }

    async function handleSubmit(e) {
        e.preventDefault();
        clearError();

        const title = document.getElementById("title").value;
        const content = document.getElementById("content").value;
        const categorySelect = document.getElementById("category-select");
        const category = categorySelect.value;

        if (!title || !content || !category) {
            showError("Title, content, and category are required");
            return;
        }

        const body = {
            title,
            content,
            category_ids: [category]
        };

        try {
            const response = await createPost(body);

            if (response.ok) {
                clearError();
                navigateTo("/homepage");
            } else {
                const errorText = await response.text();
                showError(errorText || "Failed to create post");
            }
        } catch (error) {
            showError("An error occurred. Please try again.");
        }
    }
};

export const handleCreateCategoryFormSubmit = (clearError, showError) => {
    const form = document.getElementById("create-category-form");
    if (form) {
        form.removeEventListener("submit", handleSubmit);
        form.addEventListener("submit", handleSubmit, { once: true });
    }

    async function handleSubmit(e) {
        e.preventDefault();
        clearError();

        const categoryName = document.getElementById("category-name").value;
        const messageElement = document.getElementById("category-message");

        if (!categoryName) {
            showError("Category name is required");
            return;
        }

        const body = {
            name: categoryName
        };

        try {
            const response = await createCategory(body);
            if (response.ok) {
                setFormMessage(messageElement, "success", "Category created successfully");
                document.getElementById("category-name").value = "";
                await loadCategories();
            } else {
                const errorText = await response.text();
                setFormMessage(messageElement, "error", `Error: ${errorText}`);
            }
        } catch (error) {
            showError("An error occurred. Please try again.");
        }
    }
};

export const loadCategories = async () => {
    try {
        const categorySelect = document.getElementById("category-select");
        if (!categorySelect) return;

        categorySelect.innerHTML = '<option value="">Select a category</option>';

        const categories = await getCategories();

        categories.forEach(category => {
            const option = document.createElement("option");
            option.value = category.id;
            option.textContent = category.name;
            categorySelect.appendChild(option);
        });
    } catch (error) {
        console.error("Failed to load categories:", error);
        const messageElement = document.getElementById("category-message");
        if (messageElement) {
            setFormMessage(messageElement, "error", error.message || "Failed to load categories");
        }
    }
};

export const setupMessageForm = () => {
    const messageForm = document.getElementById("message-form");
    if (messageForm) {
        messageForm.removeEventListener("submit", handleSendMessage);
        messageForm.addEventListener("submit", handleSendMessage, { once: true });
    }

    async function handleSendMessage(e) {
        e.preventDefault();

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
    }
};

export const loadAndRenderPosts = async () => {
    try {
        const postsContainer = document.getElementById("posts-container");
        if (!postsContainer) return;

        const posts = await getPosts();
        postsContainer.innerHTML = "";

        posts.forEach(post => {
            const categories = post.categories.map(category => `<span class="category">${category.name}</span>`).join(', ');
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.innerHTML = `
                <h3>${post.subject}</h3>
                <p>${post.content}</p>
                <div class="post-categories">Categories: ${categories}</div>
                <div>
                    <span>Likes: ${post.like_count}</span>
                    <span>Dislikes: ${post.dislike_count}</span>
                </div>
            `;
            postsContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error("Failed to load posts:", error);
    }
};
